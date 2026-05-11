package tui

import "testing"

func TestHeatColorBoundaries(t *testing.T) {
	cases := []struct {
		pct  float64
		want string
	}{
		{0.0, "#4ec9b0"},
		{0.10, "#dcdcaa"},
		{0.30, "#ce9178"},
		{0.60, "#f44747"},
		{1.0, "#f44747"},
	}
	for _, c := range cases {
		got := string(HeatColor(c.pct))
		if got != c.want {
			t.Errorf("HeatColor(%v) = %s, want %s", c.pct, got, c.want)
		}
	}
}

func TestHeatColorInterpolates(t *testing.T) {
	mid := string(HeatColor(0.05))
	if mid == "#4ec9b0" || mid == "#dcdcaa" {
		t.Errorf("expected interpolation between stops, got %s", mid)
	}
}

func TestHeatColorClampsNegative(t *testing.T) {
	if got := string(HeatColor(-1)); got != "#4ec9b0" {
		t.Errorf("negative should clamp to first stop, got %s", got)
	}
}

func TestAllThemesIncludesDefaultAndAmber(t *testing.T) {
	names := map[string]bool{}
	for _, th := range AllThemes() {
		names[th.Name] = true
	}
	if !names["default"] || !names["dec-amber"] {
		t.Errorf("AllThemes missing required entries: %v", names)
	}
}

func TestSetThemeSwitchesGradient(t *testing.T) {
	t.Cleanup(func() { SetTheme(DefaultTheme()) })
	SetTheme(DECAmberTheme())
	if got := string(HeatColor(0)); got != "#ff8c00" {
		t.Errorf("DEC theme low end should be orange #ff8c00, got %s", got)
	}
	if got := string(HeatColor(1)); got != "#ffffff" {
		t.Errorf("DEC theme high end should be white #ffffff, got %s", got)
	}
	if got := string(HeatColor(0.5)); got != "#ffd700" {
		t.Errorf("DEC theme midpoint should be yellow #ffd700, got %s", got)
	}
}
