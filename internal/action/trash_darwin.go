//go:build darwin

package action

import (
	"fmt"
	"os/exec"
	"strings"
)

// On macOS, ask Finder to move the file to the Trash so it lands in ~/.Trash
// and is reversible from the Finder UI.
func trashOS(path string) error {
	escaped := strings.ReplaceAll(path, `"`, `\"`)
	script := fmt.Sprintf(`tell application "Finder" to delete POSIX file "%s"`, escaped)
	cmd := exec.Command("osascript", "-e", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("osascript: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
