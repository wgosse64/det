package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Enter      key.Binding
	Back       key.Binding
	Top        key.Binding
	Bottom     key.Binding
	Trash      key.Binding
	Open       key.Binding
	Yank       key.Binding
	Rescan     key.Binding
	Sort       key.Binding
	Hidden     key.Binding
	Visualizer key.Binding
	Treemap    key.Binding
	Theme      key.Binding
	Help       key.Binding
	Quit       key.Binding
	ConfirmYes key.Binding
	ConfirmNo  key.Binding
}

func defaultKeys() keyMap {
	return keyMap{
		Up:         key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:       key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Enter:      key.NewBinding(key.WithKeys("right", "l", "enter"), key.WithHelp("→/⏎", "enter")),
		Back:       key.NewBinding(key.WithKeys("left", "h", "esc"), key.WithHelp("←/h", "back")),
		Top:        key.NewBinding(key.WithKeys("g", "home"), key.WithHelp("g", "top")),
		Bottom:     key.NewBinding(key.WithKeys("G", "end"), key.WithHelp("G", "bottom")),
		Trash:      key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "trash")),
		Open:       key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open")),
		Yank:       key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "yank path")),
		Rescan:     key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rescan")),
		Sort:       key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "sort")),
		Hidden:     key.NewBinding(key.WithKeys("."), key.WithHelp(".", "hidden")),
		Visualizer: key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "blocks")),
		Treemap:    key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "treemap")),
		Theme:      key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "theme")),
		Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		ConfirmYes: key.NewBinding(key.WithKeys("y", "Y", "enter")),
		ConfirmNo:  key.NewBinding(key.WithKeys("n", "N", "esc")),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Back, k.Trash, k.Open, k.Yank, k.Visualizer, k.Treemap, k.Theme, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Top, k.Bottom},
		{k.Enter, k.Back, k.Trash, k.Open},
		{k.Yank, k.Rescan, k.Sort, k.Hidden},
		{k.Visualizer, k.Treemap, k.Theme, k.Help, k.Quit},
	}
}
