package main

import (
	_ "embed"
	"log"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func main() {
	lines := strings.Split(input, "\n")
	left := make([]int, 0, len(lines))
	right := make([]int, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		l, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Fatalln(err)
		}
		r, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			log.Fatalln(err)
		}
		left = append(left, l)
		right = append(right, r)
	}
	sort.Ints(left)
	sort.Ints(right)
	distance := 0
	for i := 0; i < len(left); i++ {
		l := left[i]
		r := right[i]
		diff := l - r
		if diff < 0 {
			diff = -diff
		}
		distance += diff
	}
	log.Println("DISTANCE:", distance)

	score := 0
	for _, l := range left {
		count := 0
		for _, r := range right {
			if l == r {
				count++
			}
		}
		score += count * l
	}
	log.Println("SIMILARITY SCORE:", score)
}
