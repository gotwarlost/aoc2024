package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"
)

//go:embed input.txt
var input string

type point struct {
	row, col int
}

type antenna struct {
	name      string
	locations []point
}

func parse(lines []string) (map[string]*antenna, map[point]string) {
	antennas := map[string]*antenna{}
	byPosition := map[point]string{}

	for row, line := range lines {
		for col, ch := range line {
			if ch == '.' {
				continue
			}
			name := fmt.Sprintf("%c", ch)
			a := antennas[name]
			if a == nil {
				a = &antenna{name: name}
				antennas[name] = a
			}
			pt := point{row: row, col: col}
			a.locations = append(a.locations, pt)
			byPosition[pt] = name
		}
	}
	return antennas, byPosition
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func walkAway(p1, p2 point, steps int) (ret1, ret2 point) {
	rowDiff := abs(p1.row - p2.row)
	colDiff := abs(p1.col - p2.col)

	rowDir, colDir := 1, 1
	if p1.row < p2.row {
		rowDir = -1
	}
	if p1.col < p2.col {
		colDir = -1
	}
	return point{p1.row + rowDir*rowDiff*steps, p1.col + colDir*colDiff*steps},
		point{p2.row - rowDir*rowDiff*steps, p2.col - colDir*colDiff*steps}
}

func gcd(a, b int) int {
	a = abs(a)
	b = abs(b)
	for b != 0 {
		temp := b
		b = a % b
		a = temp
	}
	return a
}

func walkToward(p1, p2 point) []point {
	rowDiff := abs(p1.row - p2.row)
	colDiff := abs(p1.col - p2.col)
	g := gcd(rowDiff, colDiff)
	rowDiff /= g
	colDiff /= g

	rowDir, colDir := 1, 1
	if p1.row < p2.row {
		rowDir = -1
	}
	if p1.col < p2.col {
		colDir = -1
	}
	inGrid := func(pt point) bool {
		if pt.row <= min(p1.row, p2.row) || pt.row >= max(p1.row, p2.row) {
			return false
		}
		if pt.col <= min(p1.col, p2.col) || pt.col >= max(p1.col, p2.col) {
			return false
		}
		return true
	}
	var ret []point
	steps := 1
	for {
		current := point{p1.row - rowDir*rowDiff*steps, p1.col - colDir*colDiff*steps}
		if !inGrid(current) {
			break
		}
		ret = append(ret, current)
		steps++
	}
	return ret
}

func main() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	maxCols := len(lines[0])
	maxRows := len(lines)
	antennas, byPosition := parse(lines)

	inMap := func(pt point) bool {
		return pt.row >= 0 && pt.col >= 0 && pt.row < maxRows && pt.col < maxCols
	}

	antinodes := map[point]bool{}

	printMap := func() {
		for row := 0; row < maxRows; row++ {
			for col := 0; col < maxCols; col++ {
				pt := point{row: row, col: col}
				if _, ok := antinodes[pt]; ok {
					if x, ok := byPosition[pt]; ok {
						fmt.Print(x)
					} else {
						fmt.Print("#")
					}
					continue
				}
				fmt.Print(".")
			}
			fmt.Println()
		}
	}

	log.Println("part 1")
	for _, a := range antennas {
		locs := a.locations
		for i := 0; i < len(locs)-1; i++ {
			for j := i + 1; j < len(locs); j++ {
				p1 := locs[i]
				p2 := locs[j]
				a1, a2 := walkAway(p1, p2, 1)
				if inMap(a1) {
					antinodes[a1] = true
				}
				if inMap(a2) {
					antinodes[a2] = true
				}
			}
		}
	}
	log.Println("antinode count:", len(antinodes))

	log.Println("step 2")
	antinodes = map[point]bool{}
	for _, a := range antennas {
		locs := a.locations
		for i := 0; i < len(locs)-1; i++ {
			for j := i + 1; j < len(locs); j++ {
				p1 := locs[i]
				p2 := locs[j]

				steps := 0
				found := true
				for found {
					found = false
					a1, a2 := walkAway(p1, p2, steps)
					if inMap(a1) {
						found = true
						antinodes[a1] = true
					}
					if inMap(a2) {
						found = true
						antinodes[a2] = true
					}
					steps++
				}
				points := walkToward(p1, p2)
				for _, p := range points {
					antinodes[p] = true
				}
			}
		}
	}
	log.Println("antinodes:", len(antinodes))
	printMap()
}
