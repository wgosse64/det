package scan

import (
	"sort"
	"strings"
	"time"
)

type SortMode int

const (
	SortSize SortMode = iota
	SortName
	SortMTime
)

func (s SortMode) String() string {
	switch s {
	case SortSize:
		return "size"
	case SortName:
		return "name"
	case SortMTime:
		return "mtime"
	}
	return "?"
}

type Node struct {
	Path     string
	Name     string
	Size     int64
	IsDir    bool
	ModTime  time.Time
	Parent   *Node
	Children []*Node
	Err      error
}

func (n *Node) IsHidden() bool {
	return strings.HasPrefix(n.Name, ".") && n.Name != "." && n.Name != ".."
}

func (n *Node) PercentOfParent() float64 {
	if n.Parent == nil || n.Parent.Size <= 0 {
		return 0
	}
	return float64(n.Size) / float64(n.Parent.Size)
}

// FilteredChildren returns children sorted and optionally hiding dotfiles.
func (n *Node) FilteredChildren(mode SortMode, includeHidden bool) []*Node {
	if n == nil || n.Children == nil {
		return nil
	}
	out := make([]*Node, 0, len(n.Children))
	for _, c := range n.Children {
		if !includeHidden && c.IsHidden() {
			continue
		}
		out = append(out, c)
	}
	sortNodes(out, mode)
	return out
}

func sortNodes(nodes []*Node, mode SortMode) {
	switch mode {
	case SortSize:
		sort.SliceStable(nodes, func(i, j int) bool {
			if nodes[i].Size != nodes[j].Size {
				return nodes[i].Size > nodes[j].Size
			}
			return nodes[i].Name < nodes[j].Name
		})
	case SortName:
		sort.SliceStable(nodes, func(i, j int) bool {
			return strings.ToLower(nodes[i].Name) < strings.ToLower(nodes[j].Name)
		})
	case SortMTime:
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].ModTime.After(nodes[j].ModTime)
		})
	}
}

// RemoveChild detaches a child from this node and updates the size of every
// ancestor by subtracting the child's size.
func (n *Node) RemoveChild(child *Node) {
	if n == nil || child == nil {
		return
	}
	for i, c := range n.Children {
		if c == child {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			break
		}
	}
	for a := n; a != nil; a = a.Parent {
		a.Size -= child.Size
	}
}

// CountFiles returns the total file count under this node (excluding directories).
func (n *Node) CountFiles() int {
	if n == nil {
		return 0
	}
	if !n.IsDir {
		return 1
	}
	total := 0
	for _, c := range n.Children {
		total += c.CountFiles()
	}
	return total
}
