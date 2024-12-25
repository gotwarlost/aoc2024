package main

import (
	_ "embed"
	"log"
	"strings"
)

//go:embed input.txt
var input string

func toLayout(s string) []int {
	vals := strings.Split(s, "")
	if len(vals) != 5 {
		panic("bad line:" + s)
	}
	var ret []int
	for _, val := range vals {
		if val == "#" {
			ret = append(ret, 1)
		} else {
			ret = append(ret, 0)
		}
	}
	return ret
}

func fits(lock, key []int) int {
	for i := 0; i < len(lock); i++ {
		if lock[i]+key[i] > 5 {
			return 0
		}
	}
	return 1
}

func main() {
	var locks, keys [][]int
	lines := strings.Split(strings.TrimSpace(input), "\n")
	for i := 0; i < len(lines); i += 8 {
		layout := make([]int, 5)
		var isLock bool
		switch {
		case lines[i] == "#####":
			isLock = true
		case lines[i+6] == "#####":
			isLock = false
		default:
			panic("bad header or trailer:" + lines[i])
		}

		for j := 1; j < 6; j++ {
			index := i + j
			if !isLock {
				index = i + 6 - j
			}
			vals := toLayout(lines[index])
			for i, v := range vals {
				layout[i] += v
			}
		}
		if isLock {
			locks = append(locks, layout)
		} else {
			keys = append(keys, layout)
		}
	}

	combinations := 0
	for lock := 0; lock < len(locks); lock++ {
		for key := 0; key < len(keys); key++ {
			combinations += fits(locks[lock], keys[key])
		}
	}
	log.Println("combinations:", combinations)
}
