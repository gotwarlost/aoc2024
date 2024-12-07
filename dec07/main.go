package main

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type expr struct {
	expected int
	parts    []int
	concat   bool
}

func (e expr) withConcat() expr {
	e.concat = true
	return e
}

func (e expr) result() int {
	if e.testPartial(e.parts[0], e.parts[1:]) {
		return e.expected
	}
	return 0
}

func concatenate(x, y int) int {
	r, err := strconv.Atoi(fmt.Sprintf("%d%d", x, y))
	if err != nil {
		panic(err)
	}
	return r
}

func (e expr) testPartial(currentResult int, remaining []int) bool {
	// if nothing left, expected should be equal to actual
	if len(remaining) == 0 {
		return e.expected == currentResult
	}
	// result can only increase so if already too high bail immediately
	if currentResult > e.expected {
		return false
	}

	first, rest := remaining[0], remaining[1:]

	if e.testPartial(currentResult+first, rest) {
		return true
	}
	if e.testPartial(currentResult*first, rest) {
		return true
	}
	if !e.concat {
		return false
	}
	return e.testPartial(concatenate(currentResult, first), rest)
}

func toNum(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func main() {
	var expressions []expr
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range lines {
		rp := strings.Split(line, ":")
		if len(rp) != 2 {
			panic("bad line:" + line)
		}
		parts := strings.Split(strings.TrimSpace(rp[1]), " ")
		e := expr{expected: toNum(rp[0])}
		for _, p := range parts {
			e.parts = append(e.parts, toNum(p))
		}
		expressions = append(expressions, e)
	}

	sum := 0
	for _, e := range expressions {
		sum += e.result()
	}
	log.Println("SUM:", sum)

	sum = 0
	for _, e := range expressions {
		sum += e.withConcat().result()
	}
	log.Println("SUM CONCAT:", sum)
}
