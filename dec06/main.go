package main

import (
	_ "embed"
	"log"
	"strings"
)

//go:embed input.txt
var input string

type coord struct {
	row, col int
}

var turns = map[coord]coord{
	{-1, 0}: {0, 1},
	{0, 1}:  {1, 0},
	{1, 0}:  {0, -1},
	{0, -1}: {-1, 0},
}

var up = coord{-1, 0}

func addObstruction(current map[coord]bool, p coord) map[coord]bool {
	ret := map[coord]bool{}
	for k := range current {
		ret[k] = true
	}
	ret[p] = true
	return ret
}

func walk(rows, cols int, startPos coord, obstructions map[coord]bool) (visited map[coord][]coord, loop bool) {
	withinLimits := func(x coord) bool {
		return x.row >= 0 && x.row < rows && x.col >= 0 && x.col < cols
	}
	dir := up
	pos := startPos
	visited = map[coord][]coord{
		startPos: {up},
	}

	count := 0
	path := []coord{pos}
	for {
		count++
		next := coord{pos.row + dir.row, pos.col + dir.col}
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
			visited[next] = append(visited[next], dir)
		} else {
			visited[next] = []coord{dir}
		}
		path = append(path, next)
		pos = next
	}
}

func main() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	startPos := coord{-1, -1}
	obstructions := map[coord]bool{}
	for row, line := range lines {
		for col, ch := range line {
			switch ch {
			case '#':
				obstructions[coord{row, col}] = true
			case '^':
				startPos = coord{row, col}
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

	for point := range visited {
		if point == startPos {
			continue
		}
		_, loop := walk(rows, cols, startPos, addObstruction(obstructions, point))
		if loop {
			count++
		}
	}
	log.Println("LOOPS:", count)
}
