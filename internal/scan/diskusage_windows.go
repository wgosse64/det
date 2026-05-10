//go:build windows

package scan

import "io/fs"

// On Windows we don't have a portable cheap way to get allocated size, so we
// fall back to the apparent file length. Windows isn't a supported target for
// the trash action either, but keeping the file lets the package compile.
func diskBytes(info fs.FileInfo) int64 {
	if info == nil {
		return 0
	}
	return info.Size()
}
