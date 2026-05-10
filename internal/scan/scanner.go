package scan

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

type Progress struct {
	CurrentPath  string
	FilesScanned int64
	BytesScanned int64
	Done         bool
}

// Scan walks root and returns:
//   - the partially-built root Node (mutating live as the goroutine progresses),
//   - a progress channel emitting throttled updates,
//   - a done channel that yields the finalized root once the walk completes.
//
// Both channels close when scanning finishes or ctx is cancelled.
func Scan(ctx context.Context, root string) (*Node, <-chan Progress, <-chan *Node) {
	abs, err := filepath.Abs(root)
	if err != nil {
		abs = root
	}
	info, statErr := os.Lstat(abs)
	rootNode := &Node{
		Path:  abs,
		Name:  filepath.Base(abs),
		IsDir: statErr == nil && info.IsDir(),
	}
	if statErr != nil {
		rootNode.Err = statErr
	} else {
		rootNode.ModTime = info.ModTime()
		if !info.IsDir() {
			rootNode.Size = diskBytes(info)
		}
	}

	progressCh := make(chan Progress, 16)
	doneCh := make(chan *Node, 1)

	if rootNode.Err != nil || !rootNode.IsDir {
		go func() {
			doneCh <- rootNode
			close(progressCh)
			close(doneCh)
		}()
		return rootNode, progressCh, doneCh
	}

	var (
		mu         sync.Mutex
		filesDone  atomic.Int64
		bytesDone  atomic.Int64
		curPathPtr atomic.Pointer[string]
	)

	// dirIndex maps absolute path → *Node so we can attach children fast.
	dirIndex := map[string]*Node{abs: rootNode}

	// Progress publisher: ~10 Hz.
	progressDone := make(chan struct{})
	go func() {
		defer close(progressCh)
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-progressDone:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				cur := ""
				if p := curPathPtr.Load(); p != nil {
					cur = *p
				}
				select {
				case progressCh <- Progress{
					CurrentPath:  cur,
					FilesScanned: filesDone.Load(),
					BytesScanned: bytesDone.Load(),
				}:
				default:
				}
			}
		}
	}()

	go func() {
		defer func() {
			close(progressDone)
			doneCh <- rootNode
			close(doneCh)
		}()

		_ = filepath.WalkDir(abs, func(path string, d fs.DirEntry, walkErr error) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if path == abs {
				return nil
			}
			parentPath := filepath.Dir(path)
			mu.Lock()
			parent, ok := dirIndex[parentPath]
			mu.Unlock()
			if !ok {
				return nil
			}

			if walkErr != nil {
				node := &Node{
					Path:   path,
					Name:   filepath.Base(path),
					Parent: parent,
					IsDir:  d != nil && d.IsDir(),
					Err:    walkErr,
				}
				mu.Lock()
				parent.Children = append(parent.Children, node)
				if node.IsDir {
					dirIndex[path] = node
				}
				mu.Unlock()
				if node.IsDir {
					return fs.SkipDir
				}
				return nil
			}

			info, infoErr := d.Info()
			node := &Node{
				Path:   path,
				Name:   d.Name(),
				Parent: parent,
				IsDir:  d.IsDir(),
				Err:    infoErr,
			}
			if infoErr == nil {
				node.ModTime = info.ModTime()
				if !d.IsDir() && info.Mode().IsRegular() {
					// On-disk allocation, not apparent size: cloud-only
					// placeholders (Dropbox, iCloud, OneDrive) report 0,
					// matching `du` and what disk cleanup actually frees.
					node.Size = diskBytes(info)
				}
			}
			s := path
			curPathPtr.Store(&s)

			mu.Lock()
			parent.Children = append(parent.Children, node)
			if node.IsDir {
				dirIndex[path] = node
			} else {
				// Bubble file size up through ancestors.
				for a := parent; a != nil; a = a.Parent {
					a.Size += node.Size
				}
			}
			mu.Unlock()

			if !d.IsDir() && info != nil && info.Mode().IsRegular() {
				filesDone.Add(1)
				bytesDone.Add(node.Size)
			}
			return nil
		})
	}()

	return rootNode, progressCh, doneCh
}
