package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/wgosse/det/internal/action"
	"github.com/wgosse/det/internal/scan"
)

type confirmState struct {
	target *scan.Node
}

type Model struct {
	root    *scan.Node
	current *scan.Node
	cursor  int
	visible []*scan.Node

	width, height int
	scrollOffset  int

	scanning bool
	spinner  spinner.Model
	progress scan.Progress

	progressCh <-chan scan.Progress
	doneCh     <-chan *scan.Node
	cancelScan context.CancelFunc

	showHidden bool
	sortMode   scan.SortMode
	visualizer bool

	devMode bool

	confirm      *confirmState
	status       string
	statusExpiry time.Time

	keys keyMap
	help help.Model

	startPath string
}

func New(startPath string, devMode bool) Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	h := help.New()
	h.ShowAll = false

	m := Model{
		startPath:  startPath,
		devMode:    devMode,
		spinner:    sp,
		help:       h,
		keys:       defaultKeys(),
		showHidden: true,
		sortMode:   scan.SortSize,
	}
	m.startScan(startPath)
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.waitForScan())
}

// startScan begins (or restarts) a scan and wires the channels into the model.
// Caller must ensure they are working with the canonical model value (not a
// copy) — invoked from New and from the rescan path inside Update.
func (m *Model) startScan(path string) {
	if m.cancelScan != nil {
		m.cancelScan()
	}
	ctx, cancel := context.WithCancel(context.Background())
	m.cancelScan = cancel
	root, progressCh, doneCh := scan.Scan(ctx, path)
	m.root = root
	m.current = root
	m.cursor = 0
	m.scrollOffset = 0
	m.scanning = true
	m.progressCh = progressCh
	m.doneCh = doneCh
	m.recomputeVisible()
}

func (m Model) waitForScan() tea.Cmd {
	progressCh := m.progressCh
	doneCh := m.doneCh
	return func() tea.Msg {
		select {
		case p, ok := <-progressCh:
			if !ok {
				return nil
			}
			return scanProgressMsg(p)
		case root := <-doneCh:
			return scanDoneMsg{root: root}
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.ensureCursorVisible()
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if m.scanning {
			return m, cmd
		}
		return m, nil

	case scanProgressMsg:
		m.progress = scan.Progress(msg)
		m.recomputeVisible()
		return m, m.waitForScan()

	case scanDoneMsg:
		m.scanning = false
		if msg.root != nil {
			m.root = msg.root
			if m.current == nil || m.current.Path == "" {
				m.current = msg.root
			}
		}
		m.recomputeVisible()
		return m, nil

	case statusExpireMsg:
		if !m.statusExpiry.IsZero() && !time.Now().Before(m.statusExpiry) {
			m.status = ""
			m.statusExpiry = time.Time{}
		}
		return m, nil

	case rescanMsg:
		m.startScan(msg.path)
		return m, tea.Batch(m.spinner.Tick, m.waitForScan())

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.help.ShowAll {
		// While the help screen is open, only Help, Quit, or Back keys do
		// anything — all of them dismiss it.
		switch {
		case key.Matches(msg, m.keys.Quit):
			if m.cancelScan != nil {
				m.cancelScan()
			}
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help), key.Matches(msg, m.keys.Back):
			m.help.ShowAll = false
		}
		return m, nil
	}

	if m.confirm != nil {
		switch {
		case key.Matches(msg, m.keys.ConfirmYes):
			return m.executeTrash()
		case key.Matches(msg, m.keys.ConfirmNo):
			m.confirm = nil
			return m, nil
		}
		return m, nil
	}

	switch {
	case key.Matches(msg, m.keys.Quit):
		if m.cancelScan != nil {
			m.cancelScan()
		}
		return m, tea.Quit
	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = true
		return m, nil
	case key.Matches(msg, m.keys.Up):
		if m.cursor > 0 {
			m.cursor--
			m.ensureCursorVisible()
		}
	case key.Matches(msg, m.keys.Down):
		if m.cursor < len(m.visible)-1 {
			m.cursor++
			m.ensureCursorVisible()
		}
	case key.Matches(msg, m.keys.Top):
		m.cursor = 0
		m.scrollOffset = 0
	case key.Matches(msg, m.keys.Bottom):
		if len(m.visible) > 0 {
			m.cursor = len(m.visible) - 1
			m.ensureCursorVisible()
		}
	case key.Matches(msg, m.keys.Enter):
		if sel := m.selected(); sel != nil && sel.IsDir {
			m.current = sel
			m.cursor = 0
			m.scrollOffset = 0
			m.recomputeVisible()
		}
	case key.Matches(msg, m.keys.Back):
		if m.current != nil && m.current.Parent != nil {
			parent := m.current.Parent
			prev := m.current
			m.current = parent
			m.recomputeVisible()
			// place cursor on the dir we came from
			for i, c := range m.visible {
				if c == prev {
					m.cursor = i
					break
				}
			}
			m.ensureCursorVisible()
		}
	case key.Matches(msg, m.keys.Sort):
		m.sortMode = (m.sortMode + 1) % 3
		m.recomputeVisible()
		m.setStatus(fmt.Sprintf("Sort: %s", m.sortMode))
	case key.Matches(msg, m.keys.Hidden):
		m.showHidden = !m.showHidden
		m.recomputeVisible()
		if m.showHidden {
			m.setStatus("Hidden files: shown")
		} else {
			m.setStatus("Hidden files: hidden")
		}
	case key.Matches(msg, m.keys.Visualizer):
		m.visualizer = !m.visualizer
	case key.Matches(msg, m.keys.Open):
		if sel := m.selected(); sel != nil {
			if err := action.Reveal(sel.Path); err != nil {
				m.setStatus("Open failed: " + err.Error())
			} else {
				m.setStatus("Opened in file manager")
			}
		}
	case key.Matches(msg, m.keys.Yank):
		if sel := m.selected(); sel != nil {
			if err := action.CopyToClipboard(sel.Path); err != nil {
				m.setStatus("Copy failed: " + err.Error())
			} else {
				m.setStatus("Copied: " + sel.Path)
			}
		}
	case key.Matches(msg, m.keys.Rescan):
		path := m.startPath
		if m.current != nil {
			path = m.current.Path
		}
		return m, func() tea.Msg { return rescanMsg{path: path} }
	case key.Matches(msg, m.keys.Trash):
		if sel := m.selected(); sel != nil {
			m.confirm = &confirmState{target: sel}
		}
	}

	if m.status != "" {
		return m, m.scheduleStatusExpire()
	}
	return m, nil
}

func (m Model) executeTrash() (tea.Model, tea.Cmd) {
	c := m.confirm
	m.confirm = nil
	if c == nil || c.target == nil {
		return m, nil
	}
	target := c.target
	if err := action.Trash(target.Path, m.devMode); err != nil {
		m.setStatus("Trash failed: " + err.Error())
		return m, m.scheduleStatusExpire()
	}
	if m.devMode {
		m.setStatus("DEV: would trash " + target.Name + " (no-op)")
	} else {
		// Detach from tree on success and pull cursor back if it was last.
		parent := target.Parent
		if parent != nil {
			parent.RemoveChild(target)
		}
		m.recomputeVisible()
		if m.cursor >= len(m.visible) && m.cursor > 0 {
			m.cursor = len(m.visible) - 1
		}
		m.setStatus("Trashed: " + target.Name)
	}
	return m, m.scheduleStatusExpire()
}

func (m Model) selected() *scan.Node {
	if m.cursor < 0 || m.cursor >= len(m.visible) {
		return nil
	}
	return m.visible[m.cursor]
}

func (m *Model) recomputeVisible() {
	if m.current == nil {
		m.visible = nil
		m.cursor = 0
		return
	}
	m.visible = m.current.FilteredChildren(m.sortMode, m.showHidden)
	if m.cursor >= len(m.visible) {
		m.cursor = len(m.visible) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *Model) ensureCursorVisible() {
	rows := m.treeRows()
	if rows <= 0 {
		return
	}
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
	if m.cursor >= m.scrollOffset+rows {
		m.scrollOffset = m.cursor - rows + 1
	}
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

// treeRows returns how many rows are available for the tree body. It accounts
// for the actual rendered chrome height (header, status, possibly-wrapped
// footer) so the tree always fits the screen.
func (m Model) treeRows() int {
	chrome := lineCount(m.headerView()) + lineCount(m.footerView())
	if m.status != "" {
		chrome++
	}
	r := m.height - chrome
	if r < 1 {
		r = 1
	}
	return r
}

func lineCount(s string) int {
	if s == "" {
		return 0
	}
	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			n++
		}
	}
	return n
}

func (m *Model) setStatus(s string) {
	m.status = s
	m.statusExpiry = time.Now().Add(3 * time.Second)
}

func (m Model) scheduleStatusExpire() tea.Cmd {
	at := m.statusExpiry
	return tea.Tick(3*time.Second, func(_ time.Time) tea.Msg {
		return statusExpireMsg{at: at}
	})
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "starting…"
	}

	if m.help.ShowAll {
		return m.helpView()
	}

	var sections []string
	sections = append(sections, m.headerView())

	rows := m.treeRows()
	var body string
	if m.visualizer {
		body = m.blocksView(rows)
	} else {
		body = m.treeView(rows)
	}
	sections = append(sections, body)

	if s := m.statusLine(); s != "" {
		sections = append(sections, s)
	}
	sections = append(sections, m.footerView())

	out := strings.Join(sections, "\n")
	if m.confirm != nil {
		// Overlay the confirmation centered.
		return out + "\n" + m.confirmView()
	}
	return out
}
