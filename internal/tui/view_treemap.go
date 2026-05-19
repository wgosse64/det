package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/wgosse64/det/internal/fmtsize"
)

// treemapResult bundles the rendered string with the per-cell ownership
// grid, so the mouse handler can map (x, y) inside the panel back to a
// child index without recomputing the layout.
type treemapResult struct {
	view   string
	rects  []tmRect
	rows   int
	cols   int
	owners []int // flat row-major; -1 = empty
}

// styled cell drawn into the treemap grid before flattening to a string.
type tmCell struct {
	ch   rune
	fg   lipgloss.Color
	bg   lipgloss.Color
	bold bool
}

func (m Model) renderTreemap(rows, width int) treemapResult {
	if m.current == nil || len(m.visible) == 0 || rows < 1 || width < 4 {
		return treemapResult{view: dimStyle.Render("(empty)"), rows: rows, cols: width}
	}

	weights := make([]float64, len(m.visible))
	for i, n := range m.visible {
		weights[i] = float64(n.Size)
	}
	rects := squarifyRects(weights, tmRect{x: 0, y: 0, w: width, h: rows})

	grid := make([][]tmCell, rows)
	for i := range grid {
		grid[i] = make([]tmCell, width)
		for j := range grid[i] {
			grid[i][j] = tmCell{ch: ' '}
		}
	}
	owners := make([]int, rows*width)
	for i := range owners {
		owners[i] = -1
	}

	sepColor := currentTheme.DimFg

	for idx, rc := range rects {
		if rc.w <= 0 || rc.h <= 0 || idx >= len(m.visible) {
			continue
		}
		n := m.visible[idx]
		heat := HeatColor(n.PercentOfParent())
		isSel := idx == m.cursor

		atRight := rc.x+rc.w >= width
		atBottom := rc.y+rc.h >= rows

		// Reserve one cell on the right and bottom for inter-rect separators
		// (only when the rect has room and isn't already at the panel edge).
		innerW := rc.w
		innerH := rc.h
		if !atRight && rc.w >= 2 {
			innerW--
		}
		if !atBottom && rc.h >= 2 {
			innerH--
		}

		// Selected box uses a contrasting accent foreground; everything else
		// uses its heat color. Fill character stays `█` for both so the
		// transition reads as a color shift, not a texture change.
		fillColor := heat
		if isSel {
			fillColor = currentTheme.TreemapSelectionFg
		}

		// 1) Fill inner area with the chosen color.
		for y := rc.y; y < rc.y+innerH && y < rows; y++ {
			for x := rc.x; x < rc.x+innerW && x < width; x++ {
				grid[y][x] = tmCell{ch: '█', fg: fillColor}
				owners[y*width+x] = idx
			}
		}
		// Owners covers the separator cells too so clicks on a separator
		// still select that rect (more forgiving hit-test).
		for y := rc.y; y < rc.y+rc.h && y < rows; y++ {
			for x := rc.x; x < rc.x+rc.w && x < width; x++ {
				if owners[y*width+x] == -1 {
					owners[y*width+x] = idx
				}
			}
		}

		// 2) Separators on right + bottom, drawn AFTER fill so they overwrite.
		if !atRight && rc.w >= 2 {
			sx := rc.x + rc.w - 1
			for y := rc.y; y < rc.y+rc.h && y < rows; y++ {
				grid[y][sx] = tmCell{ch: '│', fg: sepColor}
			}
		}
		if !atBottom && rc.h >= 2 {
			sy := rc.y + rc.h - 1
			for x := rc.x; x < rc.x+rc.w && x < width; x++ {
				grid[sy][x] = tmCell{ch: '─', fg: sepColor}
			}
			// Corner where right + bottom meet.
			if !atRight && rc.w >= 2 {
				grid[sy][rc.x+rc.w-1] = tmCell{ch: '┘', fg: sepColor}
			}
		}

		// 3) Labels — only on boxes large enough to show meaningful text
		// without aggressive clipping. drawTreemapLabel applies its own
		// minimum-size check and skips silently if the rect is too tight.
		drawTreemapLabel(grid, rc, innerW, innerH, n.Name+suffixForDir(n.IsDir),
			fmtsize.Bytes(n.Size), fillColor, isSel, width)
	}

	var b strings.Builder
	for r := 0; r < rows; r++ {
		for c := 0; c < width; c++ {
			cl := grid[r][c]
			style := lipgloss.NewStyle().Foreground(cl.fg)
			if cl.bg != "" {
				style = style.Background(cl.bg)
			}
			if cl.bold {
				style = style.Bold(true)
			}
			b.WriteString(style.Render(string(cl.ch)))
		}
		if r < rows-1 {
			b.WriteByte('\n')
		}
	}

	return treemapResult{
		view:   b.String(),
		rects:  rects,
		rows:   rows,
		cols:   width,
		owners: owners,
	}
}

func suffixForDir(isDir bool) string {
	if isDir {
		return "/"
	}
	return ""
}

// minLabelInnerW is the floor for showing any text. Boxes narrower than this
// would only fit a heavily-truncated name (~3 visible chars + ellipsis) which
// reads as visual noise rather than a label, so we skip them entirely.
const minLabelInnerW = 8

// drawTreemapLabel writes name (and optionally size) into a rectangle's inner
// area, but only if the box is large enough that the label won't be clipped
// to uselessness. Label background is the rect's actual fill color so the
// selection accent stays consistent across the whole rectangle.
func drawTreemapLabel(grid [][]tmCell, rc tmRect, innerW, innerH int,
	name, size string, fillColor lipgloss.Color, selected bool, gridW int) {

	if innerW < minLabelInnerW || innerH < 1 {
		return
	}

	labelFg := contrastingFg(fillColor)

	primary := name
	if len(primary) > innerW {
		primary = truncateMiddle(primary, innerW)
	}
	primaryX := rc.x + (innerW-len(primary))/2
	primaryY := rc.y
	writeRow(grid, primaryY, primaryX, primary, labelFg, fillColor, selected, gridW)

	// Add the size on a second row if there's vertical room AND the size
	// string fits without truncation. Truncated sizes are useless.
	if innerH >= 2 && len(size) <= innerW {
		secondaryX := rc.x + (innerW-len(size))/2
		secondaryY := rc.y + 1
		writeRow(grid, secondaryY, secondaryX, size, labelFg, fillColor, selected, gridW)
	}
}

func writeRow(grid [][]tmCell, y, x int, s string, fg, bg lipgloss.Color, bold bool, gridW int) {
	if y < 0 || y >= len(grid) {
		return
	}
	for i, ch := range s {
		col := x + i
		if col < 0 || col >= gridW {
			continue
		}
		grid[y][col] = tmCell{ch: ch, fg: fg, bg: bg, bold: bold}
	}
}

// truncateMiddle shortens s to fit width, replacing the middle with "…".
// Returns s unchanged when it already fits.
func truncateMiddle(s string, width int) string {
	if len(s) <= width {
		return s
	}
	if width <= 1 {
		return s[:width]
	}
	if width <= 3 {
		return s[:width-1] + "…"
	}
	keep := width - 1
	left := keep / 2
	right := keep - left
	return s[:left] + "…" + s[len(s)-right:]
}

// contrastingFg returns black or white depending on the perceived luminance
// of the given hex color. Uses the standard 299/587/114 weights.
func contrastingFg(c lipgloss.Color) lipgloss.Color {
	s := string(c)
	if len(s) != 7 || s[0] != '#' {
		return lipgloss.Color("#000000")
	}
	r := hexByte(s[1:3])
	g := hexByte(s[3:5])
	b := hexByte(s[5:7])
	lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	if lum > 140 {
		return lipgloss.Color("#000000")
	}
	return lipgloss.Color("#ffffff")
}

func hexByte(s string) int {
	n := 0
	for _, ch := range s {
		n *= 16
		switch {
		case ch >= '0' && ch <= '9':
			n += int(ch - '0')
		case ch >= 'a' && ch <= 'f':
			n += int(ch - 'a' + 10)
		case ch >= 'A' && ch <= 'F':
			n += int(ch - 'A' + 10)
		}
	}
	return n
}
