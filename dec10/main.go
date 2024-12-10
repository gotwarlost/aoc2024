package main

import (
	_ "embed"
	"log"
	"strings"
)

//go:embed input.txt
var input string

type point struct {
	row, col int
}

func (p point) offset(o point) point {
	return point{p.row + o.row, p.col + o.col}
}

type grid struct {
	rows      int
	cols      int
	heads     []point
	points    [][]int
	end       map[point]bool
	numTrails int
}

func (g *grid) valueAt(pt point) int {
	return g.points[pt.row][pt.col]
}

func (g *grid) withinLimits(pt point) bool {
	return pt.row >= 0 && pt.col >= 0 && pt.row < g.rows && pt.col < g.cols
}

var offsets = []point{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
}

func (g *grid) next(current point) []point {
	var ret []point
	val := g.valueAt(current)
	for _, o := range offsets {
		np := current.offset(o)
		if !g.withinLimits(np) {
			continue
		}
		newV := g.valueAt(np)
		if newV == val+1 {
			ret = append(ret, np)
			if val == 8 {
				g.end[np] = true
				g.numTrails++
			}
		}
	}
	return ret
}

func (g *grid) traverse(p point) {
	nextPoints := g.next(p)
	for _, child := range nextPoints {
		g.traverse(child)
	}
}

func (g *grid) calculateScore() int {
	heads := g.heads
	score := 0
	for _, head := range heads {
		g.end = map[point]bool{}
		g.traverse(head)
		score += len(g.end)
	}
	return score
}

func main() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var grid2D [][]int
	var heads []point

	for r, l := range lines {
		var row []int
		for c, ch := range l {
			val := int(ch - '0')
			if val == 0 {
				heads = append(heads, point{row: r, col: c})
			}
			row = append(row, val)
		}
		grid2D = append(grid2D, row)
	}
	g := &grid{
		rows:   len(grid2D),
		cols:   len(grid2D[0]),
		heads:  heads,
		points: grid2D,
		end:    map[point]bool{},
	}
	log.Println("SCORE:", g.calculateScore())
	log.Println("NUM TRAILS:", g.numTrails)
}
