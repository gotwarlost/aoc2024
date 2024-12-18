package main

import (
	"container/heap"
	_ "embed"
	"log"
	"math"
	"strconv"
	"strings"
)

var (
	//go:embed base.txt
	base string
	//go:embed input.txt
	input string
)

type in struct {
	content  string
	gridSize int
	numBytes int
}

var (
	testInput = in{
		content:  base,
		gridSize: 7,
		numBytes: 12,
	}
	realInput = in{
		content:  input,
		gridSize: 71,
		numBytes: 1024,
	}
)

type point struct {
	row, col int
}

func (p point) add(row, col int) point {
	return point{p.row + row, p.col + col}
}

type grid struct {
	rows, cols int
	walls      map[point]bool
}

func (g *grid) addWall(p point) {
	g.walls[p] = true
}

func (g *grid) inGrid(p point) bool {
	return p.row >= 0 && p.row < g.rows && p.col >= 0 && p.col < g.cols
}

func (g *grid) isEmpty(p point) bool {
	return g.inGrid(p) && !g.walls[p]
}

func (g *grid) possibleNextPlaces(p point) []point {
	var ret []point
	add := func(row, col int) {
		c := p.add(row, col)
		if g.isEmpty(c) {
			ret = append(ret, c)
		}
	}
	add(-1, 0)
	add(1, 0)
	add(0, -1)
	add(0, 1)
	return ret
}

func toNum(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func toPoint(s string) point {
	parts := strings.SplitN(s, ",", 2)
	return point{row: toNum(parts[1]), col: toNum(parts[0])}
}

func parse(inp in) (ret *grid, rest []point) {
	ret = &grid{
		rows:  inp.gridSize,
		cols:  inp.gridSize,
		walls: map[point]bool{},
	}
	lines := strings.Split(strings.TrimSpace(inp.content), "\n")
	for i, line := range lines {
		if i >= inp.numBytes {
			rest = append(rest, toPoint(line))
			continue
		}
		ret.walls[toPoint(line)] = true
	}
	return ret, rest
}

type item struct {
	pt    point
	score int
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

func (g *grid) solve() (int, bool) {
	q := &queue{
		items: []item{{pt: point{0, 0}, score: 0}},
	}
	endPos := point{g.rows - 1, g.cols - 1}
	heap.Init(q)
	bestScore := math.MaxInt
	minScores := map[point]int{}
	for q.Len() > 0 {
		head := q.Pop().(item)
		if head.pt == endPos {
			if bestScore > head.score {
				bestScore = head.score
			}
			continue
		}
		for _, c := range g.possibleNextPlaces(head.pt) {
			nextScore := head.score + 1
			current, ok := minScores[c]
			if !ok || current > nextScore {
				minScores[c] = nextScore
				q.Push(item{
					pt:    c,
					score: nextScore,
					prev:  &head,
				})
			}
		}
	}
	return bestScore, bestScore != math.MaxInt
}

func main() {
	g, rest := parse(realInput)
	s, found := g.solve()
	if !found {
		panic("no solution for part 1")
	}
	log.Println("BEST SCORE:", s)
	for _, p := range rest {
		// super inefficient to run the algo over and over instead of tracking if a wall was
		// not on a known-best path. But it works in the sense of "minutes"
		g.addWall(p)
		_, hasSoln := g.solve()
		if !hasSoln {
			log.Printf("POINT OF NO SOLUTION: %d,%d\n", p.col, p.row)
			break
		}
	}
}
