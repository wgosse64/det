package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/wgosse64/det/internal/fmtsize"
	"github.com/wgosse64/det/internal/scan"
)

const (
	barWidth     = 20
	sizeColWidth = 9
	pctColWidth  = 5
)

func (m Model) treeView(rows int) string {
	if len(m.visible) == 0 {
		if m.scanning {
			return dimStyle.Render("  scanning…")
		}
		return dimStyle.Render("  (empty)")
	}

	if rows <= 0 {
		rows = 1
	}

	start := m.scrollOffset
	end := start + rows
	if end > len(m.visible) {
		end = len(m.visible)
	}

	var b strings.Builder
	for i := start; i < end; i++ {
		b.WriteString(m.renderRow(m.visible[i], i == m.cursor))
		if i < end-1 {
			b.WriteByte('\n')
		}
	}

	used := end - start
	for k := used; k < rows; k++ {
		b.WriteByte('\n')
	}
	return b.String()
}

func (m Model) renderRow(n *scan.Node, selected bool) string {
	pct := n.PercentOfParent()
	color := HeatColor(pct)

	cursor := "  "
	if selected {
		cursor = "▸ "
	}

	sizeStr := fmtsize.Bytes(n.Size)
	if len(sizeStr) < sizeColWidth {
		sizeStr = strings.Repeat(" ", sizeColWidth-len(sizeStr)) + sizeStr
	}
	pctStr := fmt.Sprintf("%4.0f%%", pct*100)
	name := n.Name
	if n.IsDir {
		name = name + "/"
	}
	if n.Err != nil {
		name = name + "  ⚠"
	}

	if selected {
		// Selected row gets a uniform background; foreground stays the heat color.
		plain := fmt.Sprintf("%s%s  %s  %s  %s", cursor, fmtsize.Bar(pct, barWidth), sizeStr, pctStr, name)
		return cursorStyle.Foreground(color).Render(padRight(plain, m.width))
	}

	bar := lipgloss.NewStyle().Foreground(color).Render(fmtsize.Bar(pct, barWidth))
	sizeRendered := lipgloss.NewStyle().Foreground(color).Render(sizeStr)
	pctRendered := dimStyle.Render(pctStr)

	var nameRendered string
	switch {
	case n.Err != nil:
		nameRendered = errorStyle.Render(name)
	case n.IsDir:
		nameRendered = dirNameStyle.Render(name)
	default:
		nameRendered = fileNameStyle.Render(name)
	}

	return fmt.Sprintf("%s%s  %s  %s  %s", cursor, bar, sizeRendered, pctRendered, nameRendered)
}

func padRight(s string, w int) string {
	if lipgloss.Width(s) >= w {
		return s
	}
	return s + strings.Repeat(" ", w-lipgloss.Width(s))
}
