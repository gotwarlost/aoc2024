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

func toNum(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func maybeSplitStone(v int) []int {
	s := fmt.Sprintf("%d", v)
	if len(s)%2 == 0 {
		s1 := s[:len(s)/2]
		s2 := s[len(s)/2:]
		return []int{toNum(s1), toNum(s2)}
	}
	return nil
}

func newStones(v int) []int {
	if v == 0 {
		return []int{1}
	}
	split := maybeSplitStone(v)
	if split != nil {
		return split
	}
	return []int{v * 2024}
}

func advance(stoneCounters map[int]int, nextStones map[int][]int) {
	snapshot := map[int]int{}
	for stone, n := range stoneCounters {
		snapshot[stone] = n
	}
	for stone, n := range snapshot {
		if nextStones[stone] == nil {
			nextStones[stone] = newStones(stone)
		}
		next := nextStones[stone]
		for _, s := range next {
			stoneCounters[s] += n
		}
		stoneCounters[stone] -= n
	}
}

func main() {
	stoneStrings := strings.Split(strings.TrimSpace(input), " ")
	var stones []int
	for _, s := range stoneStrings {
		stones = append(stones, toNum(s))
	}
	blinks := 75
	stoneMap := map[int][]int{}
	stoneCounters := map[int]int{}

	for _, stone := range stones {
		stoneCounters[stone]++
	}

	countStones := func() int {
		counter := 0
		for _, n := range stoneCounters {
			counter += n
		}
		return counter
	}

	for i := 0; i < blinks; i++ {
		advance(stoneCounters, stoneMap)
		if i == 24 {
			log.Println("COUNT 25:", countStones())
		}
	}
	log.Println("COUNT 75:", countStones())
}
