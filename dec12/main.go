package main

import (
	_ "embed"
	"fmt"
	"log"
	"sort"
	"strings"
)

//go:embed input.txt
var input string

type direction int

const (
	_ direction = iota
	top
	bottom
	left
	right
)

// edge is a direction and a value.
// for example: {top, 0} means the top-edge of the second row
type edge struct {
	dir   direction
	index int
}

// offsetPoint is an offset for a point in a specific direction.
// There are exactly 4 offset points.
type offsetPoint struct {
	row, col int
	d        direction
}

// getEdge returns the edge for a point w.r.t to this offset point.
func (o offsetPoint) getEdge(p point) edge {
	if o.d == top || o.d == bottom {
		return edge{dir: o.d, index: p.row + o.row}
	}
	return edge{dir: o.d, index: p.col + o.col}
}

var offsets = []offsetPoint{
	{-1, 0, top},
	{1, 0, bottom},
	{0, -1, left},
	{0, 1, right},
}

type point struct {
	row, col int
}

type cell struct {
	value        string
	pt           point
	regionNumber int
	perimeter    int
	edges        map[edge]bool
}

type grid struct {
	maxRows, maxCols int
	cells            map[point]*cell
}

func (g *grid) inGrid(pt point) bool {
	row := pt.row
	col := pt.col
	return row >= 0 && row < g.maxRows && col >= 0 && col < g.maxCols
}

func (g *grid) assignRegion(current *cell, region int) {
	// already assigned, noop
	if current.regionNumber != 0 {
		return
	}
	current.regionNumber = region
	for _, p := range offsets {
		pt := point{row: current.pt.row + p.row, col: current.pt.col + p.col}
		if !g.inGrid(pt) {
			current.perimeter++
			current.edges[p.getEdge(current.pt)] = true
			continue
		}
		next := g.cells[pt]
		if next.value == current.value {
			g.assignRegion(next, region)
		} else {
			current.perimeter++
			current.edges[p.getEdge(current.pt)] = true
		}
	}
}

type area struct {
	val        string
	region     int
	count      int
	perimeters int
	edges      map[edge][]point
}

// calculateSides calculates the sides for an area. This is done as follows:
// For each edge, sort all points seen for that edge in the appropriate way.
// (i.e. by columns if the edge is top or bottom or by row when it is not)
// add a side when there are non-contiguous points for an edge.
func (a *area) calculateSides() int {
	sides := 0
	for e, pts := range a.edges {
		sort.Slice(pts, func(i, j int) bool {
			if e.dir == top || e.dir == bottom {
				return pts[i].col < pts[j].col
			}
			return pts[i].row < pts[j].row
		})
		prev := -3 // sentinel value that makes the first edge at any index a side
		for _, p := range pts {
			val := p.row
			if e.dir == top || e.dir == bottom {
				val = p.col
			}
			// if not contiguous add a side
			if val != prev+1 {
				sides++
			}
			prev = val
		}
	}
	return sides
}

func main() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	g := &grid{
		cells: map[point]*cell{},
	}
	for row, line := range lines {
		for col, ch := range line {
			pt := point{row: row, col: col}
			g.cells[pt] = &cell{
				value: fmt.Sprintf("%c", ch),
				pt:    pt,
				edges: map[edge]bool{},
			}
		}
	}
	g.maxRows = len(lines)
	g.maxCols = len(lines[0])
	currentRegion := 0

	for i := 0; i < g.maxRows; i++ {
		for j := 0; j < g.maxCols; j++ {
			c := g.cells[point{row: i, col: j}]
			if c.regionNumber != 0 {
				continue
			}
			currentRegion++
			g.assignRegion(c, currentRegion)
		}
	}

	// create areas per region, keyed by region number
	areas := map[int]*area{}
	for i := 0; i < g.maxRows; i++ {
		for j := 0; j < g.maxCols; j++ {
			c := g.cells[point{row: i, col: j}]
			a := areas[c.regionNumber]
			if a == nil {
				a = &area{val: c.value, region: c.regionNumber, edges: map[edge][]point{}}
				areas[c.regionNumber] = a
			}
			a.count += 1
			a.perimeters += c.perimeter
			// accumulate points by edges for calculating sides
			for k := range c.edges {
				a.edges[k] = append(a.edges[k], c.pt)
			}
		}
	}
	out := 0
	out2 := 0
	for _, a := range areas {
		out += a.perimeters * a.count
		out2 += a.calculateSides() * a.count
	}
	log.Println("PERIMETER AREA SUM:", out)
	log.Println("SIDE AREA SUM:", out2)
}
