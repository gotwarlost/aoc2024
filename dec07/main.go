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

func (e expr) testPartial(partialResult int, rest []int) bool {
	if len(rest) == 0 {
		return e.expected == partialResult
	}
	if partialResult > e.expected {
		return false
	}
	r := rest[0]
	remaining := rest[1:]

	newResult := partialResult + r
	if e.testPartial(newResult, remaining) {
		return true
	}
	newResult = partialResult * r
	if e.testPartial(newResult, remaining) {
		return true
	}
	if !e.concat {
		return false
	}
	var err error
	newResult, err = strconv.Atoi(fmt.Sprintf("%d%d", partialResult, r))
	if err != nil {
		panic(err)
	}
	return e.testPartial(newResult, remaining)
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
