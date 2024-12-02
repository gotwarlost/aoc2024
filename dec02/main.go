package main

import (
	_ "embed"
	"log"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func getLevels(s string) []int {
	parts := strings.Split(s, " ")
	var ret []int
	for _, part := range parts {
		x, err := strconv.Atoi(part)
		if err != nil {
			panic(err)
		}
		ret = append(ret, x)
	}
	return ret
}

func isSafe(levels []int) bool {
	if len(levels) == 0 {
		panic("need at least one level")
	}
	if len(levels) == 1 {
		return true
	}
	increasing := levels[1] > levels[0]
	for i := 1; i < len(levels); i++ {
		diff := levels[i] - levels[i-1]
		if diff == 0 {
			return false
		}
		if increasing && diff < 0 {
			return false
		}
		if !increasing && diff > 0 {
			return false
		}
		if diff < 0 {
			diff = -diff
		}
		if diff > 3 {
			return false
		}
	}
	return true
}

func removeLevel(levels []int, toRemove int) []int {
	var ret []int
	for i := 0; i < len(levels); i++ {
		if i == toRemove {
			continue
		}
		ret = append(ret, levels[i])
	}
	return ret
}

func isSafeWithDampening(levels []int) bool {
	if isSafe(levels) {
		return true
	}
	for remove := 0; remove < len(levels); remove++ {
		truncLevels := removeLevel(levels, remove)
		if isSafe(truncLevels) {
			return true
		}
	}
	return false
}

func main() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	safeCount := 0
	safeWithDampeningCount := 0
	for _, l := range lines {
		levels := getLevels(l)
		if isSafe(levels) {
			safeCount++
		}
		if isSafeWithDampening(levels) {
			safeWithDampeningCount++
		}
	}
	log.Println("SAFE COUNT:", safeCount)
	log.Println("SAFE WITH DAMPENING COUNT:", safeWithDampeningCount)
}
