package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/wgosse64/det/internal/fmtsize"
)

func (m Model) headerView() string {
	title := titleStyle.Render("DET")
	curPath := m.startPath
	var curSize int64
	var curFiles int64
	if m.current != nil {
		curPath = m.current.Path
		curSize = m.current.Size
		curFiles = int64(m.current.CountFiles())
	}
	pathPart := pathStyle.Render(truncatePath(curPath, m.width-30))

	totals := totalsStyle.Render(fmt.Sprintf(
		"%s · %s files",
		fmtsize.Bytes(curSize),
		fmtsize.Comma(curFiles),
	))

	left := lipgloss.JoinHorizontal(lipgloss.Center, title, " ", pathPart)

	pad := m.width - lipgloss.Width(left) - lipgloss.Width(totals)
	if pad < 1 {
		pad = 1
	}

	rows := []string{left + strings.Repeat(" ", pad) + totals}
	if m.devMode {
		rows = append(rows, devTagStyle.Render("DEV — DELETIONS DISABLED"))
	}
	return strings.Join(rows, "\n")
}

// scanBandLine renders the live scanner status as a single fixed-height row.
// During scan: spinner + truncated current path. After: a quiet completion
// note. The path is truncated, never wrapped, so this row never pushes the
// rest of the UI off-screen.
func (m Model) scanBandLine() string {
	if m.scanning {
		spin := m.spinner.View()
		// " <spin> scanning " ≈ 12 chars of fixed prefix; leave a bit of
		// breathing room on the right so trailing path doesn't kiss the edge.
		const fixed = 14
		avail := m.width - fixed
		if avail < 10 {
			avail = 10
		}
		path := truncatePath(m.progress.CurrentPath, avail)
		return dimStyle.Render(fmt.Sprintf(" %s scanning %s", spin, path))
	}
	return statusStyle.Render(" ✓ scan complete")
}

func (m Model) footerView() string {
	if m.confirm != nil {
		return ""
	}
	if m.help.ShowAll {
		return m.help.View(m.keys)
	}
	return wrapHints(shortHints(), m.width)
}

type hint struct{ key, desc string }

func shortHints() []hint {
	return []hint{
		{"↑↓", "move"},
		{"→/⏎", "enter"},
		{"←", "back"},
		{"d", "trash"},
		{"o", "open"},
		{"y", "yank"},
		{"v", "blocks"},
		{"m", "treemap"},
		{"t", "theme"},
		{".", "hidden"},
		{"s", "sort"},
		{"r", "rescan"},
		{"?", "help"},
		{"q", "quit"},
	}
}

// wrapHints lays out hints in a uniform grid. The preferred shape is 6
// columns × 2 rows; the column count drops to whatever fits when the window
// is narrower than that.
func wrapHints(hints []hint, width int) string {
	if len(hints) == 0 {
		return ""
	}
	if width <= 0 {
		width = 80
	}
	const gap = 2
	const preferredCols = 7
	keyStyle := dimStyle.Bold(true)

	// Uniform cell width so the grid columns line up vertically.
	cellW := 0
	for _, h := range hints {
		if w := lipgloss.Width(h.key + " " + h.desc); w > cellW {
			cellW = w
		}
	}

	maxCols := preferredCols
	if len(hints) < maxCols {
		maxCols = len(hints)
	}
	cols := 1
	for c := maxCols; c >= 1; c-- {
		if c*cellW+(c-1)*gap <= width {
			cols = c
			break
		}
	}

	rows := (len(hints) + cols - 1) / cols
	var lines []string
	for r := 0; r < rows; r++ {
		var parts []string
		for c := 0; c < cols; c++ {
			i := r*cols + c
			if i >= len(hints) {
				break
			}
			h := hints[i]
			cell := keyStyle.Render(h.key) + " " + dimStyle.Render(h.desc)
			plainW := lipgloss.Width(h.key + " " + h.desc)
			if plainW < cellW {
				cell += strings.Repeat(" ", cellW-plainW)
			}
			parts = append(parts, cell)
		}
		lines = append(lines, strings.Join(parts, strings.Repeat(" ", gap)))
	}
	return strings.Join(lines, "\n")
}

func (m Model) statusLine() string {
	if m.status == "" {
		return ""
	}
	return statusStyle.Render(m.status)
}

func (m Model) confirmView() string {
	c := m.confirm
	if c == nil {
		return ""
	}
	prompt := fmt.Sprintf("Move to trash?\n\n  %s\n  (%s)\n\n[y] yes   [n] no",
		c.target.Name, fmtsize.Bytes(c.target.Size))
	if m.devMode {
		prompt += "\n\n" + devTagStyle.Render("DEV mode: nothing will actually be deleted")
	}
	box := confirmBoxStyle.Render(prompt)
	return lipgloss.Place(m.width, m.height-2, lipgloss.Center, lipgloss.Center, box)
}

func truncatePath(p string, max int) string {
	if max <= 0 || len(p) <= max {
		return p
	}
	if max < 4 {
		return p[len(p)-max:]
	}
	return "…" + p[len(p)-max+1:]
}
