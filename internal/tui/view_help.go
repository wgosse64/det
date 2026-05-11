package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Help styles read from the active theme on every call so theme changes are
// picked up immediately on the next render.
func helpTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme.HelpTitleFg).
		Background(currentTheme.HelpTitleBg).
		Bold(true).
		Padding(0, 2)
}

func helpSectionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme.HelpSectionFg).
		Bold(true).
		Underline(true).
		MarginTop(1)
}

func helpKeyStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme.HelpKeyFg).
		Bold(true)
}

func helpDescStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(currentTheme.FileFg)
}

func helpDimStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme.DimFg).
		Italic(true)
}

func helpFooterStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(currentTheme.DimFg).
		MarginTop(1)
}

type helpEntry struct {
	keys, desc string
	long       string // optional one-line elaboration
}

type helpSection struct {
	title   string
	entries []helpEntry
}

var helpSections = []helpSection{
	{
		title: "Navigation",
		entries: []helpEntry{
			{"↑/k, ↓/j", "Move cursor up / down", ""},
			{"g, G", "Jump to top / bottom of the list", ""},
			{"→/l, Enter", "Descend into the highlighted directory", ""},
			{"←/h", "Go up to the parent directory", ""},
		},
	},
	{
		title: "Cleaning Actions",
		entries: []helpEntry{
			{"d", "Move highlighted item to the system Trash",
				"Opens a confirmation. On macOS items go to ~/.Trash via Finder; on Linux they go to ~/.local/share/Trash per the freedesktop.org spec."},
			{"o", "Reveal in Finder (macOS) or open folder with xdg-open (Linux)", ""},
			{"y", "Yank — copy the item's absolute path to the system clipboard", ""},
			{"r", "Re-scan the current directory",
				"Useful after deleting things outside det, or to refresh sizes after large changes."},
		},
	},
	{
		title: "View / Filter",
		entries: []helpEntry{
			{"v", "Toggle the block visualizer",
				"Replaces the tree with a colored block grid where each child's area is proportional to its share of the directory."},
			{"m", "Toggle the treemap panel",
				"Opens a second, smaller panel below the navigator showing a squarified treemap of the current directory's children. Click a tile to select it."},
			{"s", "Cycle the sort order",
				"size (largest first) → name (alphabetical) → mtime (most recently modified first)."},
			{".", "Toggle visibility of dotfiles / hidden directories",
				"They are always counted toward sizes; this only hides them from the row list."},
			{"t", "Cycle the color theme",
				"Currently: default (cool→hot heatmap) and dec-amber (vintage DEC amber CRT, with an orange→yellow→white gradient on the size bars)."},
		},
	},
	{
		title: "Misc",
		entries: []helpEntry{
			{"?", "Toggle this help screen", ""},
			{"q, Ctrl-C", "Quit", ""},
		},
	},
}

func (m Model) helpView() string {
	width := m.width
	if width < 20 {
		width = 20
	}

	title := helpTitleStyle().Render("Disk Exploration Tool — Help")
	subtitle := helpDimStyle().Render("A WinDirStat-style disk explorer for terminal disk cleanup.")
	if m.devMode {
		subtitle += "  " + devTagStyle.Render("DEV — DELETIONS DISABLED")
	}

	keyColW := keyColumnWidth(helpSections)
	if keyColW > width/3 {
		keyColW = width / 3
	}

	var b strings.Builder
	b.WriteString(title)
	b.WriteByte('\n')
	b.WriteString(subtitle)
	b.WriteByte('\n')

	for _, section := range helpSections {
		b.WriteByte('\n')
		b.WriteString(helpSectionStyle().Render(section.title))
		b.WriteByte('\n')
		for _, e := range section.entries {
			keyCell := helpKeyStyle().Render(padPlain(e.keys, keyColW))
			line := "  " + keyCell + "  " + helpDescStyle().Render(e.desc)
			b.WriteString(line)
			b.WriteByte('\n')
			if e.long != "" {
				wrapped := wrapText(e.long, width-keyColW-6)
				for _, ln := range wrapped {
					b.WriteString(strings.Repeat(" ", keyColW+4))
					b.WriteString(helpDimStyle().Render(ln))
					b.WriteByte('\n')
				}
			}
		}
	}

	b.WriteString(helpFooterStyle().Render("Press ? or Esc to close   ·   q to quit"))

	out := b.String()
	// Pad to full screen height.
	got := lineCount(out)
	if got < m.height {
		out += strings.Repeat("\n", m.height-got)
	}
	return out
}

func keyColumnWidth(sections []helpSection) int {
	max := 0
	for _, s := range sections {
		for _, e := range s.entries {
			if w := lipgloss.Width(e.keys); w > max {
				max = w
			}
		}
	}
	return max
}

func padPlain(s string, w int) string {
	cur := lipgloss.Width(s)
	if cur >= w {
		return s
	}
	return s + strings.Repeat(" ", w-cur)
}

// wrapText breaks s into lines no wider than w runes, splitting on word
// boundaries. It does not handle ANSI styling — pass plain text only.
func wrapText(s string, w int) []string {
	if w <= 0 {
		return []string{s}
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	cur := words[0]
	for _, word := range words[1:] {
		if len(cur)+1+len(word) > w {
			lines = append(lines, cur)
			cur = word
			continue
		}
		cur += " " + word
	}
	lines = append(lines, cur)
	return lines
}
