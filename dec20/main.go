package main

import (
	"container/heap"
	_ "embed"
	"fmt"
	"log"
	"math"
	"slices"
	"strings"
)

//go:embed input.txt
var input string

type point struct {
	row, col int
}

func (p point) add(row, col int) point {
	return point{p.row + row, p.col + col}
}

func (p point) manhattanDistance(other point) int {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	return abs(other.row-p.row) + abs(other.col-p.col)
}

type item struct {
	pt    point
	score int
	prev  *item
}

func (i *item) solution() solution {
	var ret []point
	node := i
	for node != nil {
		ret = append(ret, node.pt)
		node = node.prev
	}
	slices.Reverse(ret)
	return ret
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
	soln       *solution
}

func (m *maze) hasWallAt(p point) bool {
	return m.walls[p]
}

func (m *maze) possibleNextPlaces(p point) []point {
	var ret []point
	add := func(row, col int) {
		c := p.add(row, col)
		if !m.hasWallAt(c) {
			ret = append(ret, c)
		}
	}
	add(-1, 0)
	add(1, 0)
	add(0, -1)
	add(0, 1)
	return ret
}

func (m *maze) blockingWalls(p point) []point {
	var ret []point
	add := func(row, col int) {
		c := p.add(row, col)
		if m.hasWallAt(c) {
			ret = append(ret, c)
		}
	}
	add(-1, 0)
	add(1, 0)
	add(0, -1)
	add(0, 1)
	return ret
}

func (m *maze) solve() solution {
	q := &queue{
		items: []item{{pt: m.startPos, score: 0}},
	}
	endPos := m.endPos
	heap.Init(q)
	bestScore := math.MaxInt
	var sol solution
	minScores := map[point]int{}
	for q.Len() > 0 {
		head := q.Pop().(item)
		if head.pt == endPos {
			if bestScore > head.score {
				bestScore = head.score
				sol = head.solution()
			}
			continue
		}
		for _, c := range m.possibleNextPlaces(head.pt) {
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
	if bestScore == math.MaxInt {
		panic("no solution")
	}
	return sol
}

func (m *maze) dump(s solution) {
	visited := map[point]bool{}
	for _, c := range s {
		visited[c] = true
	}
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			p := point{i, j}
			_, ok := visited[p]
			switch {
			case p == m.startPos:
				fmt.Print("S")
			case p == m.endPos:
				fmt.Print("E")
			case ok:
				fmt.Print("\u2588")
			case m.hasWallAt(p):
				fmt.Print("|")
			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
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

type solution []point

type saving struct {
	pt1, pt2 point
	saving   int
}

func (s solution) savings(maxCheats int) map[int]int {
	var savings []saving
	savingsBySaving := map[int]int{}
	addSaving := func(p1, p2 point, saved int) {
		savings = append(savings, saving{
			pt1:    p1,
			pt2:    p2,
			saving: saved,
		})
		savingsBySaving[saved]++
	}
	for i, p := range s {
		for j := i + 1; j < len(s); j++ {
			candidate := s[j]
			distance := j - i
			cheatDistance := p.manhattanDistance(candidate)
			if cheatDistance <= maxCheats && cheatDistance < distance {
				addSaving(p, candidate, distance-cheatDistance)
			}
		}
	}
	/*
		var keys []int
		for k := range savingsBySaving {
			keys = append(keys, k)
			sort.Ints(keys)
		}
		for _, k := range keys {
			log.Println(k, ":", savingsBySaving[k])
		}
	*/
	return savingsBySaving
}

func main() {
	m := newMaze(input)
	s := m.solve()
	m.dump(s)

	gt100 := func(maxCheats int) {
		sbys := s.savings(maxCheats)
		counter := 0
		for k, v := range sbys {
			if k >= 100 {
				counter += v
			}
		}
		log.Println("CHEATS SAVING AT LEAST HUNDRED PS:", counter)
	}
	gt100(2)
	gt100(20)
}
