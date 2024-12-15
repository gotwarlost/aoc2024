package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"
)

//go:embed input.txt
var input string

type kind int

const (
	empty kind = iota
	wall
	box
	robot
	lbox
	rbox
)

type offset struct {
	row, col int
}

func (o offset) String() string {
	switch {
	case o.row == -1:
		return "^"
	case o.row == 1:
		return "v"
	case o.col == -1:
		return "<"
	default:
		return ">"
	}
}

type point struct {
	row, col int
}

func (p point) withOffset(o offset) point {
	return point{p.row + o.row, p.col + o.col}
}

var (
	left  = offset{0, -1}
	right = offset{0, 1}
	up    = offset{-1, 0}
	down  = offset{1, 0}
)

type grid struct {
	part2      bool
	rows, cols int
	positions  map[point]kind
	robotPos   point
}

func (g *grid) moveSingleBox(p point, o offset) {
	what := g.positions[p]
	switch what {
	case box, lbox:
	default:
		panic(fmt.Sprintf("attempt to move non-box: %v", what))
	}
	delete(g.positions, p)
	moveRight := false
	if what == lbox {
		moveRight = true
		delete(g.positions, p.withOffset(right))
	}
	g.positions[p.withOffset(o)] = what
	if moveRight {
		g.positions[p.withOffset(right).withOffset(o)] = rbox
	}
}

func (g *grid) thingAt(p point) kind {
	t, ok := g.positions[p]
	if !ok {
		if g.robotPos == p {
			return robot
		}
		return empty
	}
	return t
}

func (g *grid) dump(title string) {
	fmt.Println(title)
	for row := 0; row < g.rows; row++ {
		for col := 0; col < g.cols; col++ {
			what := g.thingAt(point{row, col})
			switch what {
			case robot:
				fmt.Print("@")
			case empty:
				fmt.Print(".")
			case wall:
				fmt.Print("#")
			case box:
				fmt.Print("O")
			case lbox:
				fmt.Print("[")
			case rbox:
				fmt.Print("]")
			default:
				panic("booyah")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (g *grid) canMoveBox(current point, o offset) bool {
	newP := current.withOffset(o)
	what := g.thingAt(newP)
	switch what {
	case wall:
		return false
	case empty:
		return true
	case box:
		return g.canMoveBox(newP, o)
	case lbox:
		if o.row != 0 {
			return g.canMoveBox(newP, o) && g.canMoveBox(newP.withOffset(right), o)
		}
		return g.canMoveBox(newP, o)
	case rbox:
		if o.row != 0 {
			return g.canMoveBox(newP, o) && g.canMoveBox(newP.withOffset(left), o)
		}
		return g.canMoveBox(newP, o)
	default:
		panic(fmt.Sprintf("unexpected value: %v", what))
	}
}

func (g *grid) moveBox(current point, o offset) {
	what := g.thingAt(current)
	switch what {
	case box, lbox, rbox:
	default:
		panic(fmt.Sprintf("internal error: %v", what))
	}
	newP := current.withOffset(o)
	newWhat := g.thingAt(newP)
	simpleMove := func() {
		delete(g.positions, current)
		g.positions[newP] = what
	}
	doSimple := func() {
		if newWhat == empty {
			simpleMove()
		} else {
			g.moveBox(newP, o)
			simpleMove()
		}
	}
	if g.part2 && o.row != 0 {
		if what == rbox {
			g.moveBox(current.withOffset(left), o)
			return
		}
		altOffset := right
		altP := newP.withOffset(altOffset)
		altWhat := g.thingAt(altP)

		if newWhat == lbox && altWhat == rbox {
			g.moveBox(newP, o)
		} else {
			if newWhat != empty {
				g.moveBox(newP, o)
			}
			if altWhat != empty {
				g.moveBox(altP, o)
			}
		}
		simpleMove()
		current = current.withOffset(altOffset)
		what = g.thingAt(current)
		newP = altP
		simpleMove()
	} else {
		doSimple()
	}
}

func (g *grid) advance(m offset) {
	nextPos := g.robotPos.withOffset(m)
	thing := g.positions[nextPos]
	switch thing {
	case empty:
		g.robotPos = nextPos
		return
	case wall:
		return
	}

	canDo := g.canMoveBox(nextPos, m)
	canDo2 := true
	if g.part2 && m.row != 0 {
		if g.thingAt(nextPos) == lbox {
			canDo2 = g.canMoveBox(nextPos.withOffset(right), m)
		} else {
			canDo2 = g.canMoveBox(nextPos.withOffset(left), m)
		}
	}
	if canDo && canDo2 {
		g.moveBox(nextPos, m)
		g.robotPos = nextPos
	}
}

func (g *grid) solution() int {
	solution := 0
	for row := 0; row < g.rows; row++ {
		for col := 0; col < g.cols; col++ {
			what := g.thingAt(point{row, col})
			if what == box || what == lbox {
				solution += 100*row + col
			}
		}
	}
	return solution
}

func parse(lines []string, part2 bool) (g *grid, moves []offset) {
	g = &grid{part2: part2}
	g.positions = map[point]kind{}
	isPhase2 := false
	g.rows = 0
	g.cols = len(lines[0])
	if part2 {
		g.cols *= 2
	}
	offsetMap := map[int32]offset{
		'<': left,
		'>': right,
		'^': up,
		'v': down,
	}
	for i, l := range lines {
		if l == "" {
			isPhase2 = true
		}
		if !isPhase2 {
			g.rows++
			for j, ch := range l {
				switch ch {
				case '#':
					if part2 {
						g.positions[point{i, 2 * j}] = wall
						g.positions[point{i, 2*j + 1}] = wall
					} else {
						g.positions[point{i, j}] = wall
					}
				case 'O':
					if part2 {
						g.positions[point{i, 2 * j}] = lbox
						g.positions[point{i, 2*j + 1}] = rbox
					} else {
						g.positions[point{i, j}] = box
					}
				case '@':
					if part2 {
						g.robotPos = point{i, 2 * j}
					} else {
						g.robotPos = point{i, j}
					}
				}
			}
			continue
		}
		for _, ch := range l {
			moves = append(moves, offsetMap[ch])
		}
	}
	return g, moves
}

func run(part2 bool) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	g, moves := parse(lines, part2)
	g.dump("initial state")

	for _, m := range moves {
		g.advance(m)
	}
	g.dump("end state")
	log.Println("SOLUTION:", g.solution())
}

func main() {
	log.Println("================= PART 1 ==============")
	run(false)
	log.Println("================= PART 2 ==============")
	run(true)
}
