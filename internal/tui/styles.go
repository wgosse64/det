package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// HeatStop is one anchor point of a color gradient — a fraction in [0, 1]
// and an RGB color.
type HeatStop struct {
	Pct     float64
	R, G, B int
}

// Theme bundles every color decision the TUI makes. Themes are immutable;
// switching themes calls SetTheme(t) which rebuilds the package-level styles.
type Theme struct {
	Name string

	// Window border surrounding the whole TUI.
	Border lipgloss.Color

	// Title chip ("DET").
	TitleFg lipgloss.Color
	TitleBg lipgloss.Color

	PathFg   lipgloss.Color
	TotalsFg lipgloss.Color
	DimFg    lipgloss.Color

	// DEV mode banner.
	DevFg lipgloss.Color
	DevBg lipgloss.Color

	// Selected row background.
	CursorFg lipgloss.Color
	CursorBg lipgloss.Color

	// File / dir name colors in the tree.
	DirFg   lipgloss.Color
	FileFg  lipgloss.Color
	ErrorFg lipgloss.Color

	// Status line (success / info).
	StatusFg lipgloss.Color

	// Confirm dialog border.
	ConfirmFg lipgloss.Color

	// Help screen.
	HelpTitleFg   lipgloss.Color
	HelpTitleBg   lipgloss.Color
	HelpSectionFg lipgloss.Color
	HelpKeyFg     lipgloss.Color

	// Color gradient used for size bars and the block visualizer.
	HeatStops []HeatStop
}

var (
	currentTheme Theme
	heatStops    []HeatStop

	titleStyle        lipgloss.Style
	pathStyle         lipgloss.Style
	totalsStyle       lipgloss.Style
	dimStyle          lipgloss.Style
	devTagStyle       lipgloss.Style
	cursorStyle       lipgloss.Style
	dirNameStyle      lipgloss.Style
	fileNameStyle     lipgloss.Style
	errorStyle        lipgloss.Style
	statusStyle       lipgloss.Style
	confirmBoxStyle   lipgloss.Style
	windowBorderStyle lipgloss.Style
)

func init() { SetTheme(DefaultTheme()) }

// SetTheme installs t as the current theme and rebuilds every package-level
// style. Safe to call from the tea Update loop (no concurrent renders).
func SetTheme(t Theme) {
	currentTheme = t
	heatStops = t.HeatStops

	titleStyle = lipgloss.NewStyle().
		Foreground(t.TitleFg).
		Background(t.TitleBg).
		Bold(true).
		Padding(0, 1)

	pathStyle = lipgloss.NewStyle().
		Foreground(t.PathFg).
		Bold(true)

	totalsStyle = lipgloss.NewStyle().Foreground(t.TotalsFg)
	dimStyle = lipgloss.NewStyle().Foreground(t.DimFg)

	devTagStyle = lipgloss.NewStyle().
		Foreground(t.DevFg).
		Background(t.DevBg).
		Bold(true).
		Padding(0, 1)

	cursorStyle = lipgloss.NewStyle().
		Foreground(t.CursorFg).
		Background(t.CursorBg).
		Bold(true)

	dirNameStyle = lipgloss.NewStyle().
		Foreground(t.DirFg).
		Bold(true)

	fileNameStyle = lipgloss.NewStyle().Foreground(t.FileFg)

	errorStyle = lipgloss.NewStyle().
		Foreground(t.ErrorFg).
		Bold(true)

	statusStyle = lipgloss.NewStyle().Foreground(t.StatusFg)

	confirmBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.ConfirmFg).
		Padding(1, 2)

	windowBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border)
}

// HeatColor returns a lipgloss color for the given fraction in [0, 1] using
// the current theme's gradient stops.
func HeatColor(pct float64) lipgloss.Color {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	for i := 0; i < len(heatStops)-1; i++ {
		a, b := heatStops[i], heatStops[i+1]
		if pct >= a.Pct && pct <= b.Pct {
			t := 0.0
			if b.Pct > a.Pct {
				t = (pct - a.Pct) / (b.Pct - a.Pct)
			}
			r := int(float64(a.R) + t*float64(b.R-a.R))
			g := int(float64(a.G) + t*float64(b.G-a.G))
			bl := int(float64(a.B) + t*float64(b.B-a.B))
			return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, bl))
		}
	}
	last := heatStops[len(heatStops)-1]
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", last.R, last.G, last.B))
}

// DefaultTheme is the cool/teal → red heatmap on a near-black UI.
func DefaultTheme() Theme {
	return Theme{
		Name:          "default",
		Border:        lipgloss.Color("#5f5fd7"),
		TitleFg:       lipgloss.Color("#ffffff"),
		TitleBg:       lipgloss.Color("#5f5fd7"),
		PathFg:        lipgloss.Color("#9cdcfe"),
		TotalsFg:      lipgloss.Color("#dcdcaa"),
		DimFg:         lipgloss.Color("#7f848e"),
		DevFg:         lipgloss.Color("#ffffff"),
		DevBg:         lipgloss.Color("#d70000"),
		CursorFg:      lipgloss.Color("#ffffff"),
		CursorBg:      lipgloss.Color("#3a3d41"),
		DirFg:         lipgloss.Color("#9cdcfe"),
		FileFg:        lipgloss.Color("#d4d4d4"),
		ErrorFg:       lipgloss.Color("#f44747"),
		StatusFg:      lipgloss.Color("#b5cea8"),
		ConfirmFg:     lipgloss.Color("#f44747"),
		HelpTitleFg:   lipgloss.Color("#ffffff"),
		HelpTitleBg:   lipgloss.Color("#5f5fd7"),
		HelpSectionFg: lipgloss.Color("#dcdcaa"),
		HelpKeyFg:     lipgloss.Color("#9cdcfe"),
		HeatStops: []HeatStop{
			{0.0, 0x4e, 0xc9, 0xb0},
			{0.10, 0xdc, 0xdc, 0xaa},
			{0.30, 0xce, 0x91, 0x78},
			{0.60, 0xf4, 0x47, 0x47},
			{1.0, 0xf4, 0x47, 0x47},
		},
	}
}

// DECAmberTheme evokes a vintage Digital Equipment Corp. amber phosphor CRT:
// warm orange chrome with an orange → yellow → white gradient for size bars.
func DECAmberTheme() Theme {
	const (
		ink         = lipgloss.Color("#1a0a00") // near-black for inverse text
		deepAmber   = lipgloss.Color("#cc6600")
		amber       = lipgloss.Color("#ff8c00")
		paleAmber   = lipgloss.Color("#ffb84d")
		yellow      = lipgloss.Color("#ffd700")
		darkBurnish = lipgloss.Color("#2a1500") // dark brown for cursor bg — contrasts with every gradient stop
		danger      = lipgloss.Color("#ff3300")
	)
	return Theme{
		Name:          "dec-amber",
		Border:        amber,
		TitleFg:       ink,
		TitleBg:       amber,
		PathFg:        paleAmber,
		TotalsFg:      yellow,
		DimFg:         deepAmber,
		DevFg:         ink,
		DevBg:         danger,
		CursorFg:      lipgloss.Color("#ffffff"),
		CursorBg:      darkBurnish,
		DirFg:         paleAmber,
		FileFg:        amber,
		ErrorFg:       danger,
		StatusFg:      yellow,
		ConfirmFg:     yellow,
		HelpTitleFg:   ink,
		HelpTitleBg:   amber,
		HelpSectionFg: yellow,
		HelpKeyFg:     paleAmber,
		HeatStops: []HeatStop{
			{0.0, 0xff, 0x8c, 0x00}, // orange
			{0.5, 0xff, 0xd7, 0x00}, // yellow
			{1.0, 0xff, 0xff, 0xff}, // white
		},
	}
}

// AllThemes returns the theme cycle order.
func AllThemes() []Theme {
	return []Theme{DefaultTheme(), DECAmberTheme()}
}
