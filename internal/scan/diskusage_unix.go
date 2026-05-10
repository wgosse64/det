//go:build !windows

package scan

import (
	"io/fs"
	"syscall"
)

// diskBytes returns the number of bytes a file actually occupies on disk.
//
// This differs from info.Size() (the apparent / logical length) for:
//   - sparse files
//   - Dropbox / iCloud / OneDrive online-only placeholders (Blocks == 0)
//   - HFS+/APFS compressed system files
//
// st_blocks is reported in 512-byte units on every POSIX system. We trust
// that and avoid an extra statfs to read the filesystem block size.
func diskBytes(info fs.FileInfo) int64 {
	if info == nil {
		return 0
	}
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		return int64(st.Blocks) * 512
	}
	return info.Size()
}
