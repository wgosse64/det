package fmtsize

import (
	"strings"

	"github.com/dustin/go-humanize"
)

// Bytes returns a short human-readable size like "4.2 GB".
func Bytes(n int64) string {
	if n < 0 {
		return "—"
	}
	return humanize.Bytes(uint64(n))
}

// Bar renders a fixed-width progress bar where pct is in [0, 1].
// Returned string is exactly width runes wide using █ for filled and ░ for empty.
func Bar(pct float64, width int) string {
	if width <= 0 {
		return ""
	}
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	filled := int(pct*float64(width) + 0.5)
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

// Comma formats an int with thousands separators.
func Comma(n int64) string {
	return humanize.Comma(n)
}
