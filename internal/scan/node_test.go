package scan

import (
	"testing"
	"time"
)

func makeTree() *Node {
	root := &Node{Path: "/r", Name: "r", IsDir: true}
	a := &Node{Path: "/r/a", Name: "a", Size: 30, Parent: root}
	b := &Node{Path: "/r/b", Name: "b", Size: 50, Parent: root, ModTime: time.Now()}
	c := &Node{Path: "/r/.cache", Name: ".cache", Size: 20, Parent: root}
	root.Children = []*Node{a, b, c}
	root.Size = 100
	return root
}

func TestFilteredChildrenSortBySize(t *testing.T) {
	r := makeTree()
	got := r.FilteredChildren(SortSize, true)
	if len(got) != 3 {
		t.Fatalf("want 3 got %d", len(got))
	}
	if got[0].Name != "b" || got[1].Name != "a" || got[2].Name != ".cache" {
		t.Errorf("size order wrong: %s %s %s", got[0].Name, got[1].Name, got[2].Name)
	}
}

func TestFilteredChildrenHidesHidden(t *testing.T) {
	r := makeTree()
	got := r.FilteredChildren(SortSize, false)
	if len(got) != 2 {
		t.Fatalf("want 2 got %d", len(got))
	}
	for _, n := range got {
		if n.IsHidden() {
			t.Errorf("hidden node leaked: %s", n.Name)
		}
	}
}

func TestRemoveChildUpdatesAncestorSizes(t *testing.T) {
	r := makeTree()
	target := r.Children[0]
	r.RemoveChild(target)
	if r.Size != 70 {
		t.Errorf("want size 70 after removing 30, got %d", r.Size)
	}
	if len(r.Children) != 2 {
		t.Errorf("want 2 children after remove, got %d", len(r.Children))
	}
}

func TestPercentOfParent(t *testing.T) {
	r := makeTree()
	if got := r.Children[1].PercentOfParent(); got != 0.5 {
		t.Errorf("want 0.5 got %v", got)
	}
}

func TestIsHidden(t *testing.T) {
	cases := map[string]bool{".cache": true, "cache": false, "..": false, ".": false}
	for name, want := range cases {
		n := &Node{Name: name}
		if got := n.IsHidden(); got != want {
			t.Errorf("%q: want %v got %v", name, want, got)
		}
	}
}
