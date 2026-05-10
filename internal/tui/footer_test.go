package tui

import (
	"strings"
	"testing"
)

func TestWrapHintsPrefers6x2Grid(t *testing.T) {
	out := wrapHints(shortHints(), 200)
	rows := strings.Count(out, "\n") + 1
	if rows != 2 {
		t.Errorf("wide terminal should produce a 6×2 grid (2 rows), got %d:\n%s", rows, out)
	}
}

func TestWrapHintsAddsRowsWhenNarrow(t *testing.T) {
	out := wrapHints(shortHints(), 30)
	rows := strings.Count(out, "\n") + 1
	if rows < 3 {
		t.Errorf("narrow terminal should use more rows, got %d:\n%s", rows, out)
	}
}

func TestWrapHintsEachItemPresent(t *testing.T) {
	out := wrapHints(shortHints(), 30)
	for _, h := range shortHints() {
		if !strings.Contains(out, h.desc) {
			t.Errorf("hint %q missing from wrapped footer:\n%s", h.desc, out)
		}
	}
}
