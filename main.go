package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/wgosse64/det/internal/tui"
)

func main() {
	devMode := flag.Bool("dev", false, "developer mode: deletions become no-ops and a DEV tag is shown")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Disk Exploration Tool — a TUI disk usage explorer\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [--dev] [path]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "det: %v\n", err)
		os.Exit(2)
	}
	if _, err := os.Stat(abs); err != nil {
		fmt.Fprintf(os.Stderr, "det: %v\n", err)
		os.Exit(2)
	}

	model := tui.New(abs, *devMode)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "det: %v\n", err)
		os.Exit(1)
	}
}
