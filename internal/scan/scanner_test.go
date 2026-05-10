package scan

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestScanBuildsTreeAndAggregatesSizes(t *testing.T) {
	dir := t.TempDir()
	const apparent = 1600
	mustWrite(t, filepath.Join(dir, "a.txt"), make([]byte, 1000))
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(sub, "b.txt"), make([]byte, 500))
	mustWrite(t, filepath.Join(dir, ".hidden"), make([]byte, 100))

	root, _, doneCh := Scan(context.Background(), dir)
	select {
	case <-doneCh:
	case <-time.After(5 * time.Second):
		t.Fatal("scan timed out")
	}

	// We measure on-disk allocation (st_blocks * 512), so small files round
	// up to the filesystem block size — root.Size should be at least the
	// apparent sum but not absurdly larger.
	if root.Size < apparent {
		t.Errorf("root size %d below apparent total %d", root.Size, apparent)
	}
	if root.Size > 1<<20 {
		t.Errorf("root size %d unexpectedly large for ~1.6 KB of data", root.Size)
	}
	if !root.IsDir {
		t.Errorf("want root to be a dir")
	}
	if len(root.Children) != 3 {
		t.Errorf("want 3 children, got %d", len(root.Children))
	}
}

func mustWrite(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// Sparse files are the closest stand-in for a Dropbox/iCloud "online-only"
// placeholder: huge apparent size, ~zero on-disk allocation. The scan should
// report the on-disk size, not the logical length.
func TestScanReportsOnDiskSizeForSparseFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sparse.bin")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	// 100 MB hole; on most filesystems no blocks are allocated.
	if err := f.Truncate(100 << 20); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	root, _, doneCh := Scan(context.Background(), dir)
	select {
	case <-doneCh:
	case <-time.After(5 * time.Second):
		t.Fatal("scan timed out")
	}

	if root.Size > 1<<20 {
		t.Errorf("expected sparse file to report tiny on-disk size, got %d bytes", root.Size)
	}
}
