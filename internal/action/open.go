package action

import (
	"os/exec"
	"path/filepath"
	"runtime"
)

// Reveal opens the system file manager pointing at path. On macOS, this
// reveals the file in Finder; on Linux, xdg-open is used to open the
// containing directory.
func Reveal(path string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", "-R", path).Run()
	case "linux":
		return exec.Command("xdg-open", filepath.Dir(path)).Run()
	default:
		return exec.Command("xdg-open", filepath.Dir(path)).Run()
	}
}
