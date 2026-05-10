# det — Disk Exploration Tool

A WinDirStat-style disk-usage explorer for the terminal, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea). Designed for **disk cleaning and maintenance** on macOS and Linux.

- Async tree scan with live progress — the UI is interactive immediately.
- Children sorted by size with color-coded heat bars (green → yellow → orange → red, scaled to share of the parent directory).
- One-key `Move to Trash` with confirmation; safer than `rm` and reversible from the OS.
- Optional block-visualizer mode for an at-a-glance treemap-style view.
- `--dev` mode that visibly disables deletions for safe demos.

## Install

```sh
go install github.com/wgosse/det@latest
```

Or from a local clone:

```sh
go build -o det .
```

## Usage

```sh
det                    # scan the current directory
det ~/Downloads        # scan a specific path
det --dev /tmp         # dev mode: trash actions become no-ops
```

## Keys

| Key | Action |
|---|---|
| `↑`/`k`, `↓`/`j` | Move cursor |
| `→`/`l`/`Enter` | Descend into directory |
| `←`/`h`/`Esc` | Go up |
| `g` / `G` | Top / bottom |
| `d` | Move selected to trash (confirm) |
| `o` | Reveal in Finder / xdg-open |
| `y` | Copy absolute path |
| `r` | Re-scan current view |
| `s` | Cycle sort: size → name → mtime |
| `.` | Toggle hidden files |
| `v` | Toggle block visualizer |
| `?` | Toggle help |
| `q` / `Ctrl-C` | Quit |

## Trash behavior

- **macOS**: routes through Finder via `osascript`, so trashed items appear in `~/.Trash` and can be restored with ⌘Z or "Put Back".
- **Linux**: uses the freedesktop.org Trash spec (`~/.local/share/Trash` or `$XDG_DATA_HOME/Trash`).
- **`--dev`**: nothing is actually moved; the status line shows `DEV: would trash …` and the header carries a red `[DEV — DELETIONS DISABLED]` tag.
