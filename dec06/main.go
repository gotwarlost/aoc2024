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

type direction struct {
	row, col int
}

var turns = map[direction]direction{
	{-1, 0}: {0, 1},
	{0, 1}:  {1, 0},
	{1, 0}:  {0, -1},
	{0, -1}: {-1, 0},
}

var up = direction{-1, 0}

func walk(rows, cols int, startPos point, obstructions map[point]bool) (visited map[point][]direction, loop bool) {
	withinLimits := func(x point) bool {
		return x.row >= 0 && x.row < rows && x.col >= 0 && x.col < cols
	}
	dir := up
	pos := startPos
	visited = map[point][]direction{
		startPos: {up},
	}

	for {
		next := point{pos.row + dir.row, pos.col + dir.col}
		if !withinLimits(next) {
			return visited, false
		}
		if obstructions[next] {
			dir = turns[dir]
			continue
		}
		if dirs, ok := visited[next]; ok {
			for _, d := range dirs {
				if d == dir {
					return visited, true
				}
			}
		}
		visited[next] = append(visited[next], dir)
		pos = next
	}
}

func withObstruction(current map[point]bool, p point) map[point]bool {
	ret := map[point]bool{}
	for k := range current {
		ret[k] = true
	}
	ret[p] = true
	return ret
}

func main() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	startPos := point{-1, -1}
	obstructions := map[point]bool{}
	for row, line := range lines {
		for col, ch := range line {
			switch ch {
			case '#':
				obstructions[point{row, col}] = true
			case '^':
				startPos = point{row, col}
			}
		}
	}
	cols := len(lines[0])
	rows := len(lines)

	visited, l := walk(rows, cols, startPos, obstructions)
	if l {
		panic("unexpected loop")
	}
	log.Println("UNIQ CELLS:", len(visited))

	count := 0

	for p := range visited {
		if p == startPos {
			continue
		}
		_, loop := walk(rows, cols, startPos, withObstruction(obstructions, p))
		if loop {
			count++
		}
	}
	log.Println("LOOPS:", count)
}
