package fmtsize

import (
	"strings"
	"testing"
)

func TestBarFull(t *testing.T) {
	got := Bar(1.0, 10)
	if got != strings.Repeat("█", 10) {
		t.Errorf("full bar wrong: %q", got)
	}
}

func TestBarEmpty(t *testing.T) {
	got := Bar(0, 10)
	if got != strings.Repeat("░", 10) {
		t.Errorf("empty bar wrong: %q", got)
	}
}

func TestBarHalf(t *testing.T) {
	got := Bar(0.5, 10)
	want := strings.Repeat("█", 5) + strings.Repeat("░", 5)
	if got != want {
		t.Errorf("half bar wrong: %q", got)
	}
}

func TestBarClampsAbove1(t *testing.T) {
	got := Bar(2.0, 4)
	if got != strings.Repeat("█", 4) {
		t.Errorf("clamp wrong: %q", got)
	}
}

func TestBarZeroWidth(t *testing.T) {
	if Bar(0.5, 0) != "" {
		t.Error("zero width should be empty")
	}
}

func TestBytesNegative(t *testing.T) {
	if Bytes(-1) != "—" {
		t.Errorf("negative byte format wrong")
	}
}
