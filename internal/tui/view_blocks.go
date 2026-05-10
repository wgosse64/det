package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/wgosse64/det/internal/fmtsize"
)

// blocksView renders the children of m.current as a row-major grid of colored
// unicode blocks. Each child gets a number of cells proportional to its size;
// every cell is colored by the child's percent-of-parent.
func (m Model) blocksView(rows int) string {
	if rows <= 1 || m.width <= 4 {
		return dimStyle.Render("(too small)")
	}
	if m.current == nil || len(m.visible) == 0 {
		return dimStyle.Render("  (empty)")
	}

	// Reserve a few rows for the legend below.
	gridRows := rows - 4
	if gridRows < 3 {
		gridRows = rows - 1
	}
	if gridRows < 1 {
		gridRows = 1
	}
	gridCols := m.width
	totalCells := gridRows * gridCols

	parentSize := m.current.Size
	if parentSize <= 0 {
		return dimStyle.Render("  (no size)")
	}

	// Allocate cells per child proportional to size, ensure at least 1 cell
	// per child until cells run out.
	type alloc struct {
		idx   int
		cells int
	}
	allocs := make([]alloc, 0, len(m.visible))
	used := 0
	for i, n := range m.visible {
		c := int(float64(n.Size) / float64(parentSize) * float64(totalCells))
		if c < 1 && totalCells-used > 0 {
			c = 1
		}
		if c > totalCells-used {
			c = totalCells - used
		}
		allocs = append(allocs, alloc{idx: i, cells: c})
		used += c
		if used >= totalCells {
			break
		}
	}

	// Build a flat color array of `totalCells` entries.
	cells := make([]lipgloss.Color, 0, totalCells)
	owners := make([]int, 0, totalCells)
	for _, a := range allocs {
		n := m.visible[a.idx]
		col := HeatColor(n.PercentOfParent())
		for k := 0; k < a.cells; k++ {
			cells = append(cells, col)
			owners = append(owners, a.idx)
		}
	}
	for len(cells) < totalCells {
		cells = append(cells, lipgloss.Color("#1e1e1e"))
		owners = append(owners, -1)
	}

	var b strings.Builder
	for r := 0; r < gridRows; r++ {
		for c := 0; c < gridCols; c++ {
			i := r*gridCols + c
			style := lipgloss.NewStyle().Foreground(cells[i])
			if owners[i] == m.cursor {
				style = style.Reverse(true)
			}
			b.WriteString(style.Render("█"))
		}
		if r < gridRows-1 {
			b.WriteByte('\n')
		}
	}

	// Legend: highlighted child name + size + pct.
	if m.cursor < len(m.visible) {
		sel := m.visible[m.cursor]
		legend := dimStyle.Render("► ") +
			lipgloss.NewStyle().Foreground(HeatColor(sel.PercentOfParent())).Bold(true).Render(
				sel.Name,
			) +
			dimStyle.Render("  ") +
			fileNameStyle.Render(fmtsize.Bytes(sel.Size))
		b.WriteString("\n\n" + legend)
	}

	return b.String()
}
