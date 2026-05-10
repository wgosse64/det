package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/wgosse/det/internal/scan"
)

func helpModel(devMode bool) Model {
	root := &scan.Node{Path: "/tmp", Name: "tmp", IsDir: true}
	m := Model{
		startPath:  "/tmp",
		devMode:    devMode,
		spinner:    spinner.New(),
		help:       help.New(),
		keys:       defaultKeys(),
		showHidden: true,
		sortMode:   scan.SortSize,
		root:       root,
		current:    root,
	}
	m.help.ShowAll = true
	return m
}

func TestHelpViewListsAllSections(t *testing.T) {
	m := helpModel(false)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)
	out := m.View()
	for _, want := range []string{"Navigation", "Cleaning Actions", "View / Filter", "Misc", "Press ? or Esc to close"} {
		if !strings.Contains(out, want) {
			t.Errorf("help view missing %q", want)
		}
	}
}

func TestHelpViewShowsDevTagWhenEnabled(t *testing.T) {
	m := helpModel(true)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)
	if !strings.Contains(m.View(), "DEV") {
		t.Errorf("dev mode tag missing from help view")
	}
}

func TestHelpEscClosesHelp(t *testing.T) {
	m := helpModel(false)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(Model)
	if m.help.ShowAll {
		t.Errorf("Esc should close help")
	}
}
