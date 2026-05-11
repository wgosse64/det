package tui

import (
	"math"
	"sort"
)

type tmRect struct{ x, y, w, h int }

type tmItem struct {
	weight float64
	index  int
}

// squarifyRects packs proportional weights into the given rectangle using the
// squarified treemap algorithm (Bruls, Huijsen & van Wijk, 2000). Returned
// rectangles are in the same order as `weights` and have integer cell
// dimensions that fully cover `area` (modulo per-row rounding).
//
// We compensate for terminal-cell aspect ratio (cells are roughly twice as
// tall as wide) by treating one cell-row as 2 weight-units of height when
// scoring aspect ratios — this produces visually squarish boxes.
func squarifyRects(weights []float64, area tmRect) []tmRect {
	out := make([]tmRect, len(weights))
	if len(weights) == 0 || area.w <= 0 || area.h <= 0 {
		return out
	}

	items := make([]tmItem, 0, len(weights))
	var totalW float64
	for i, w := range weights {
		if w <= 0 {
			continue
		}
		items = append(items, tmItem{weight: w, index: i})
		totalW += w
	}
	if totalW <= 0 || len(items) == 0 {
		return out
	}
	sort.Slice(items, func(a, b int) bool { return items[a].weight > items[b].weight })

	// Cell-aspect compensation: scale rectangle width down by 2 in the
	// math, then scale the resulting widths back up when placing.
	const cellAspect = 2.0

	// Convert weights to areas in (compensated) cell units.
	// Total compensated area = (area.w / cellAspect) * area.h.
	compArea := (float64(area.w) / cellAspect) * float64(area.h)
	areas := make([]float64, len(items))
	for i, it := range items {
		areas[i] = it.weight / totalW * compArea
	}

	// Compensated rect for the algorithm.
	work := tmRect{
		x: 0,
		y: 0,
		w: int(math.Round(float64(area.w) / cellAspect)),
		h: area.h,
	}
	if work.w < 1 {
		work.w = 1
	}

	tmpRects := make([]tmRect, len(items))
	squarifyHelper(items, areas, work, tmpRects)

	// De-compensate widths back to terminal cells.
	for i, r := range tmpRects {
		tmpRects[i] = tmRect{
			x: int(math.Round(float64(r.x) * cellAspect)),
			y: r.y,
			w: int(math.Round(float64(r.w) * cellAspect)),
			h: r.h,
		}
	}

	// Snap rightmost / bottommost to area edges to avoid 1-cell gaps from rounding.
	for i, r := range tmpRects {
		if r.x+r.w > area.w {
			tmpRects[i].w = area.w - r.x
		}
		if r.y+r.h > area.h {
			tmpRects[i].h = area.h - r.y
		}
		if tmpRects[i].w < 1 {
			tmpRects[i].w = 1
		}
		if tmpRects[i].h < 1 {
			tmpRects[i].h = 1
		}
	}

	for k, it := range items {
		out[it.index] = tmpRects[k]
	}
	return out
}

func squarifyHelper(items []tmItem, areas []float64, r tmRect, out []tmRect) {
	if len(items) == 0 || r.w <= 0 || r.h <= 0 {
		return
	}

	// Always lay out one row spanning the shorter side.
	short := r.w
	if r.h < short {
		short = r.h
	}
	if short < 1 {
		short = 1
	}

	rowEnd := 1
	bestRatio := worstAspectRatio(areas[:1], short)

	for rowEnd < len(items) {
		ratio := worstAspectRatio(areas[:rowEnd+1], short)
		if ratio > bestRatio {
			break
		}
		bestRatio = ratio
		rowEnd++
	}

	rowAreas := areas[:rowEnd]
	var rowSum float64
	for _, a := range rowAreas {
		rowSum += a
	}

	if r.w >= r.h {
		// Vertical strip on the left.
		colW := int(math.Round(rowSum / float64(short)))
		if colW < 1 {
			colW = 1
		}
		if colW > r.w {
			colW = r.w
		}
		y := r.y
		remH := r.h
		for i := range items[:rowEnd] {
			h := int(math.Round(rowAreas[i] / float64(colW)))
			if h < 1 {
				h = 1
			}
			if i == rowEnd-1 || h > remH {
				h = remH
			}
			out[i] = tmRect{x: r.x, y: y, w: colW, h: h}
			y += h
			remH -= h
			if remH <= 0 {
				break
			}
		}
		squarifyHelper(items[rowEnd:], areas[rowEnd:],
			tmRect{x: r.x + colW, y: r.y, w: r.w - colW, h: r.h}, out[rowEnd:])
	} else {
		// Horizontal strip on top.
		rowH := int(math.Round(rowSum / float64(short)))
		if rowH < 1 {
			rowH = 1
		}
		if rowH > r.h {
			rowH = r.h
		}
		x := r.x
		remW := r.w
		for i := range items[:rowEnd] {
			w := int(math.Round(rowAreas[i] / float64(rowH)))
			if w < 1 {
				w = 1
			}
			if i == rowEnd-1 || w > remW {
				w = remW
			}
			out[i] = tmRect{x: x, y: r.y, w: w, h: rowH}
			x += w
			remW -= w
			if remW <= 0 {
				break
			}
		}
		squarifyHelper(items[rowEnd:], areas[rowEnd:],
			tmRect{x: r.x, y: r.y + rowH, w: r.w, h: r.h - rowH}, out[rowEnd:])
	}
}

func worstAspectRatio(areas []float64, w int) float64 {
	if len(areas) == 0 || w == 0 {
		return math.Inf(1)
	}
	var s, mx, mn float64
	mx = -math.Inf(1)
	mn = math.Inf(1)
	for _, a := range areas {
		s += a
		if a > mx {
			mx = a
		}
		if a < mn {
			mn = a
		}
	}
	sw := float64(w)
	r1 := (sw * sw * mx) / (s * s)
	r2 := (s * s) / (sw * sw * mn)
	if r1 > r2 {
		return r1
	}
	return r2
}
