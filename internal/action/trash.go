package action

// Trash sends a path to the system trash. On macOS this delegates to
// trashDarwin (osascript → Finder), and on Linux to trashLinux (freedesktop
// spec via hymkor/trash-go). When devMode is true the call is a no-op.
func Trash(path string, devMode bool) error {
	if devMode {
		return nil
	}
	return trashOS(path)
}
