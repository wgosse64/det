package tui

import (
	"time"

	"github.com/wgosse/det/internal/scan"
)

// scanProgressMsg is delivered each time the scanner pushes a progress tick.
type scanProgressMsg scan.Progress

// scanDoneMsg signals the scan goroutine has finished.
type scanDoneMsg struct{ root *scan.Node }

// statusExpireMsg clears the transient status line at a specific time.
type statusExpireMsg struct{ at time.Time }

// rescanMsg starts a fresh scan of the given path (used by the rescan key).
type rescanMsg struct{ path string }
