package main

import (
	"container/heap"
	_ "embed"
	"fmt"
	"log"
	"math"
	"strings"
)

//go:embed input.txt
var input string

type point struct {
	row, col int
}

func (p point) add(other point) point {
	return point{p.row + other.row, p.col + other.col}
}

type direction int

const (
	_ direction = iota
	dirNorth
	dirSouth
	dirWest
	dirEast
)

type dirCost struct {
	dir  direction
	cost int
}

func (d direction) candidates() []dirCost {
	switch d {
	case dirNorth:
		return []dirCost{{dirNorth, 1}, {dirEast, 1001}, {dirWest, 1001}}
	case dirSouth:
		return []dirCost{{dirSouth, 1}, {dirEast, 1001}, {dirWest, 1001}}
	case dirWest:
		return []dirCost{{dirWest, 1}, {dirNorth, 1001}, {dirSouth, 1001}}
	case dirEast:
		return []dirCost{{dirEast, 1}, {dirNorth, 1001}, {dirSouth, 1001}}
	default:
		panic("invalid direction")
	}
}

func (d direction) offset() point {
	switch d {
	case dirNorth:
		return point{-1, 0}
	case dirSouth:
		return point{1, 0}
	case dirWest:
		return point{0, -1}
	case dirEast:
		return point{0, 1}
	default:
		panic("invalid direction")
	}
}

func (d direction) String() string {
	switch d {
	case dirNorth:
		return "\u25b2"
	case dirSouth:
		return "\u25bc"
	case dirWest:
		return "\u25c0"
	case dirEast:
		return "\u25b6"
	default:
		panic("invalid direction")
	}
}

type item struct {
	pt    point
	score int
	dir   direction
	prev  *item
}

type queue struct {
	items []item
}

func (q *queue) Len() int           { return len(q.items) }
func (q *queue) Less(i, j int) bool { return q.items[i].score < q.items[j].score }
func (q *queue) Swap(i, j int)      { q.items[i], q.items[j] = q.items[j], q.items[i] }
func (q *queue) Push(x any)         { q.items = append(q.items, x.(item)) }
func (q *queue) Pop() any {
	ret := q.items[q.Len()-1]
	q.items = q.items[:q.Len()-1]
	return ret
}

type maze struct {
	rows, cols int
	points     map[point]bool
	walls      map[point]bool
	startPos   point
	endPos     point
}

func (m *maze) hasWallAt(p point) bool {
	return m.walls[p]
}

func (m *maze) dump(visited map[point]direction) {
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			p := point{i, j}
			d, ok := visited[p]
			switch {
			case p == m.startPos:
				fmt.Print("S")
			case p == m.endPos:
				fmt.Print("E")
			case ok:
				fmt.Print(d)
			case m.hasWallAt(p):
				fmt.Print("|")
			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

type result struct {
	visited map[point]bool
	routes  []map[point]direction
}

func (r *result) add(it item) {
	node := &it
	route := map[point]direction{}
	for node != nil {
		r.visited[node.pt] = true
		route[node.pt] = node.dir
		node = node.prev
	}
	r.routes = append(r.routes, route)
}

func (m *maze) solve() {
	q := &queue{
		items: []item{{pt: m.startPos, dir: dirEast, score: 0}},
	}
	heap.Init(q)
	bestScore := math.MaxInt
	bestPoints := map[int]*result{}
	minScores := map[pointDir]int{}
	for q.Len() > 0 {
		head := q.Pop().(item)
		if head.pt == m.endPos {
			if bestScore > head.score {
				bestScore = head.score
			}
			res := bestPoints[head.score]
			if res == nil {
				res = &result{
					visited: map[point]bool{},
				}
				bestPoints[head.score] = res
			}
			res.add(head)
			continue
		}
		for _, c := range head.dir.candidates() {
			nextPos := head.pt.add(c.dir.offset())
			if !m.points[nextPos] {
				continue
			}
			nextScore := head.score + c.cost
			pd := pointDir{nextPos, c.dir}
			current, ok := minScores[pd]
			if !ok || current >= nextScore {
				minScores[pd] = nextScore
				q.Push(item{
					pt:    nextPos,
					score: nextScore,
					dir:   c.dir,
					prev:  &head,
				})
			}
		}
	}
	log.Println("BEST SCORE:", bestScore)
	res := bestPoints[bestScore]
	log.Println("POINT COUNT:", len(res.visited))
	m.dump(res.routes[0])
}

func newMaze(s string) *maze {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	ret := &maze{
		rows:   len(lines),
		cols:   len(lines[0]),
		walls:  make(map[point]bool),
		points: make(map[point]bool),
	}
	for i, l := range lines {
		for j, ch := range l {
			if ch == '#' {
				ret.walls[point{i, j}] = true
				continue
			}
			pt := point{i, j}
			ret.points[pt] = true
			switch ch {
			case 'S':
				ret.startPos = pt
			case 'E':
				ret.endPos = pt
			}
		}
	}
	return ret
}

type pointDir struct {
	pt  point
	dir direction
}

func main() {
	m := newMaze(input)
	m.solve()
}
