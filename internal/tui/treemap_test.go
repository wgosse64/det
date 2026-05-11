package tui

import "testing"

func TestSquarifyAreasAreProportional(t *testing.T) {
	weights := []float64{50, 30, 15, 5}
	rects := squarifyRects(weights, tmRect{x: 0, y: 0, w: 40, h: 20})
	if len(rects) != 4 {
		t.Fatalf("want 4 rects, got %d", len(rects))
	}
	// Largest weight should also produce the largest rect area.
	areas := make([]int, len(rects))
	for i, r := range rects {
		areas[i] = r.w * r.h
	}
	if areas[0] < areas[1] || areas[1] < areas[2] || areas[2] < areas[3] {
		t.Errorf("rect areas are not monotonically decreasing: %v", areas)
	}
}

func TestSquarifyHandlesEmpty(t *testing.T) {
	rects := squarifyRects(nil, tmRect{w: 10, h: 5})
	if len(rects) != 0 {
		t.Errorf("expected 0 rects for empty input, got %d", len(rects))
	}
}

func TestSquarifyZeroWeightsSkipped(t *testing.T) {
	rects := squarifyRects([]float64{10, 0, 5}, tmRect{w: 30, h: 10})
	if len(rects) != 3 {
		t.Fatalf("want 3 rects, got %d", len(rects))
	}
	// The zero-weight item should get a zero-sized (or unset) rect.
	if rects[1].w*rects[1].h != 0 {
		t.Errorf("zero-weight item should have zero area, got rect %v", rects[1])
	}
}

func TestSquarifyFillsTinyArea(t *testing.T) {
	// 2 items in a tiny rect — must not crash and each item should have
	// at least one cell.
	rects := squarifyRects([]float64{1, 1}, tmRect{w: 2, h: 1})
	for i, r := range rects {
		if r.w < 1 || r.h < 1 {
			t.Errorf("rect %d collapsed to %v", i, r)
		}
	}
}
