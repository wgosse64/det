package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/wgosse/det/internal/scan"
)

// modelFor builds a fully-populated Model around an in-memory tree, avoiding
// the live scan side effect of New() so tests are deterministic.
func modelFor(root *scan.Node, devMode bool) Model {
	m := Model{
		startPath:  root.Path,
		devMode:    devMode,
		spinner:    spinner.New(),
		help:       help.New(),
		keys:       defaultKeys(),
		showHidden: true,
		sortMode:   scan.SortSize,
		root:       root,
		current:    root,
	}
	m.recomputeVisible()
	return m
}

func TestViewRendersChildren(t *testing.T) {
	root := &scan.Node{Path: "/tmp", Name: "tmp", IsDir: true, Size: 300}
	a := &scan.Node{Path: "/tmp/a", Name: "a.bin", Size: 200, Parent: root}
	b := &scan.Node{Path: "/tmp/b", Name: "b.bin", Size: 100, Parent: root}
	root.Children = []*scan.Node{a, b}

	m := modelFor(root, false)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = updated.(Model)

	out := m.View()
	if !strings.Contains(out, "DET") {
		t.Errorf("View missing DET title: %q", out)
	}
	if !strings.Contains(out, "a.bin") || !strings.Contains(out, "b.bin") {
		t.Errorf("View missing children: %q", out)
	}
}

func TestDevModeShowsTag(t *testing.T) {
	root := &scan.Node{Path: "/tmp", Name: "tmp", IsDir: true}
	m := modelFor(root, true)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updated.(Model)
	out := m.View()
	if !strings.Contains(out, "DEV") {
		t.Errorf("expected DEV tag in view, got: %q", out)
	}
}

func TestNilCurrentDoesNotCrashView(t *testing.T) {
	m := Model{
		startPath: "/tmp",
		spinner:   spinner.New(),
		help:      help.New(),
		keys:      defaultKeys(),
	}
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updated.(Model)
	_ = m.View() // must not panic even with current == nil
}
