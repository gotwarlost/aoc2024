package main

import (
	_ "embed"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

var (
	reA     = regexp.MustCompile(`^Button A: X\+(\d+), Y\+(\d+)`)
	reB     = regexp.MustCompile(`^Button B: X\+(\d+), Y\+(\d+)`)
	rePrize = regexp.MustCompile(`^Prize: X=(\d+), Y=(\d+)`)
)

type offset struct {
	x, y int
}

func (o offset) float64() (x, y float64) {
	return float64(o.x), float64(o.y)
}

type problem struct {
	a     offset
	b     offset
	prize offset
}

type solution struct {
	m, n int64
}

func (s *solution) cost() int64 {
	if s == nil {
		return 0
	}
	return 3*s.m + s.n
}

const epsilon = 1e-6

func isValid(f float64) (int64, bool) {
	if (math.Abs(f-math.Floor(f)) < epsilon) || (math.Abs(f-math.Ceil(f)) < epsilon) {
		return int64(math.Round(f)), true
	}
	return 0, false
}

func equal(f1, f2 float64) bool {
	return math.Abs(f1-f2) < epsilon
}

func (p *problem) solve(prizeOffset int64, constrain100 bool) *solution {
	/*
		x*ax + y*bx = px
		x*ay + y*by = py
		y = (px*ay - py*ax) / (bx*ay - by*ax)
	*/
	debug := func(p string, v any) {
		//log.Println(p,":",v)
	}
	debug("problem:", *p)
	px, py := p.prize.float64()
	px += float64(prizeOffset)
	py += float64(prizeOffset)
	ax, ay := p.a.float64()
	bx, by := p.b.float64()

	//sameEquation := false
	ratio := px / py
	if equal(ax/ay, ratio) && equal(bx/by, ratio) {
		// there is a solution we can find even for this given integer constraints and costs
		// but why bother if the data doesn't have this?
		panic("same equation")
	}

	n0 := (px*ay - py*ax) / (bx*ay - by*ax)
	m0 := (px - n0*bx) / ax

	m, valM := isValid(m0)
	n, valN := isValid(n0)
	if valM && valN {
		s := solution{m: m, n: n}
		if constrain100 {
			if m > 100 || n > 100 {
				debug("no solution", nil)
				return nil
			}
		}
		debug("solution:", s)
		return &s
	}
	debug("no solution", nil)
	return nil
}

func main() {
	lines := strings.Split(strings.TrimSpace(input)+"\n", "\n")
	if len(lines)%4 != 0 {
		log.Fatalln("not a 4-multiple")
	}

	extractOffset := func(re *regexp.Regexp, s string) offset {
		matches := re.FindStringSubmatch(s)
		x, _ := strconv.Atoi(matches[1])
		y, _ := strconv.Atoi(matches[2])
		return offset{x, y}
	}

	var problems []problem
	for i := 0; i < len(lines); i += 4 {
		a := extractOffset(reA, lines[i])
		b := extractOffset(reB, lines[i+1])
		p := extractOffset(rePrize, lines[i+2])
		problems = append(problems, problem{a: a, b: b, prize: p})
	}

	var totalCost1, totalCost2 int64
	for _, p := range problems {
		c := p.solve(0, true).cost()
		totalCost1 += c
		c2 := p.solve(10000000000000, false).cost()
		totalCost2 += c2
	}
	log.Println("TOTAL COST1:", totalCost1)
	log.Println("TOTAL COST2:", totalCost2)
}
