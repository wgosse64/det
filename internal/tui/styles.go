package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Heat gradient stops used for both bar fill and size text.
var heatStops = []struct {
	pct float64
	r   int
	g   int
	b   int
}{
	{0.0, 0x4e, 0xc9, 0xb0}, // teal
	{0.10, 0xdc, 0xdc, 0xaa}, // yellow
	{0.30, 0xce, 0x91, 0x78}, // orange
	{0.60, 0xf4, 0x47, 0x47}, // red
	{1.0, 0xf4, 0x47, 0x47},
}

// HeatColor returns a lipgloss color for the given fraction in [0, 1].
func HeatColor(pct float64) lipgloss.Color {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	for i := 0; i < len(heatStops)-1; i++ {
		a, b := heatStops[i], heatStops[i+1]
		if pct >= a.pct && pct <= b.pct {
			t := 0.0
			if b.pct > a.pct {
				t = (pct - a.pct) / (b.pct - a.pct)
			}
			r := int(float64(a.r) + t*float64(b.r-a.r))
			g := int(float64(a.g) + t*float64(b.g-a.g))
			bl := int(float64(a.b) + t*float64(b.b-a.b))
			return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, bl))
		}
	}
	last := heatStops[len(heatStops)-1]
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", last.r, last.g, last.b))
}

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#5f5fd7")).
			Bold(true).
			Padding(0, 1)

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9cdcfe")).
			Bold(true)

	totalsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#dcdcaa"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7f848e"))

	devTagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#d70000")).
			Bold(true).
			Padding(0, 1)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#3a3d41")).
			Bold(true)

	dirNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9cdcfe")).
			Bold(true)

	fileNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#d4d4d4"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f44747")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#b5cea8"))

	confirmBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#f44747")).
			Padding(1, 2)
)
