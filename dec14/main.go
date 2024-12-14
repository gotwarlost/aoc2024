package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var re = regexp.MustCompile(`p=(\d+),(\d+) v=(-?\d+),(-?\d+)`)

//go:embed input.txt
var input string

// var gx, gy = 11, 7
var gx, gy = 101, 103

type robot struct {
	x, y   int
	vx, vy int
}

func (r *robot) move(g *grid) {
	x1 := r.x + r.vx
	y1 := r.y + r.vy
	if x1 >= g.cols {
		x1 -= g.cols
	}
	if x1 < 0 {
		x1 += g.cols
	}
	if y1 >= g.rows {
		y1 -= g.rows
	}
	if y1 < 0 {
		y1 += g.rows
	}
	r.x = x1
	r.y = y1
}

type grid struct {
	rows, cols int
}

func (g *grid) inMiddle(x, y int) bool {
	return x == g.cols/2 || y == g.rows/2
}

func (g *grid) solution(robots []*robot) int {
	var q1, q2, q3, q4 int
	middles := 0
	for _, r := range robots {
		if g.inMiddle(r.x, r.y) {
			middles++
			continue
		}
		switch {
		case r.x < g.cols/2:
			if r.y < g.rows/2 {
				q1++
			} else {
				q3++
			}
		default:
			if r.y < g.rows/2 {
				q2++
			} else {
				q4++
			}
		}
	}
	log.Println("Q:", q1, q2, q3, q4, "M=", middles)
	return q1 * q2 * q3 * q4
}

type point struct {
	x, y int
}

func (g *grid) dump(i int, robots []*robot) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "-------------------------- %d  --------------------------------", i)
	m := map[point]int{}
	for _, r := range robots {
		m[point{r.x, r.y}]++
	}
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			n := m[point{j, i}]
			if n == 0 {
				fmt.Fprintf(&b, " ")
			} else {
				fmt.Fprintf(&b, ".")
			}
		}
		fmt.Fprintf(&b, "\n")
	}
	fmt.Println(b.String())
}

func toNum(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func main() {
	g := grid{gy, gx}
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var robots []*robot
	for i, l := range lines {
		matches := re.FindStringSubmatch(l)
		if matches == nil {
			panic(fmt.Sprintf("Booyah %d: %q", i, l))
		}
		r := robot{
			x:  toNum(matches[1]),
			y:  toNum(matches[2]),
			vx: toNum(matches[3]),
			vy: toNum(matches[4]),
		}
		robots = append(robots, &r)
	}

	for i := 0; i < 10000; i++ {
		byRow := map[int]int{}
		byCol := map[int]int{}
		for _, r := range robots {
			r.move(&g)
			byCol[r.x]++
			byRow[r.y]++
		}
		threshold := 30
		rows := 0
		for _, n := range byRow {
			if n >= threshold {
				rows++
			}
		}
		cols := 0
		for _, n := range byCol {
			if n >= threshold {
				cols++
			}
		}
		if rows > 1 && cols > 1 {
			log.Println("Potential solution at:", i+1)
			g.dump(i, robots)
		}
		if i == 99 {
			log.Printf("Solution: %d", g.solution(robots))
		}
	}
}
