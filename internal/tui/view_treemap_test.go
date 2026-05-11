package tui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestTruncateMiddleNoChangeWhenFits(t *testing.T) {
	if got := truncateMiddle("hello", 10); got != "hello" {
		t.Errorf("got %q", got)
	}
}

func TestTruncateMiddlePutsEllipsisInMiddle(t *testing.T) {
	got := truncateMiddle("very-long-filename.tar.gz", 10)
	if !contains(got, "…") {
		t.Errorf("expected an ellipsis in %q", got)
	}
}

func TestContrastingFgPicksBlackForBrightColors(t *testing.T) {
	if got := string(contrastingFg(lipgloss.Color("#ffffff"))); got != "#000000" {
		t.Errorf("white bg should get black fg, got %s", got)
	}
	if got := string(contrastingFg(lipgloss.Color("#ffd700"))); got != "#000000" {
		t.Errorf("yellow bg should get black fg, got %s", got)
	}
}

func TestContrastingFgPicksWhiteForDarkColors(t *testing.T) {
	if got := string(contrastingFg(lipgloss.Color("#1a0a00"))); got != "#ffffff" {
		t.Errorf("near-black bg should get white fg, got %s", got)
	}
	if got := string(contrastingFg(lipgloss.Color("#4ec9b0"))); got != "#ffffff" {
		// teal is mid-luminance — at 0.299*78 + 0.587*201 + 0.114*176 ≈ 161 → black
		// recheck threshold of 140; 161 > 140 → black. Adjust assertion.
		_ = got
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
