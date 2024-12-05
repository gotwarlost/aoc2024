package main

import (
	_ "embed"
	"log"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func midNumber(parts []int) int {
	return parts[len(parts)/2]
}

func isBadPair(ordering map[int][]int, left, right int) bool {
	after := ordering[right]
	for _, a := range after {
		if a == left {
			return true
		}
	}
	return false
}

func isFixNeeded(parts []int, ordering map[int][]int, autoFix bool) bool {
	for i := 0; i < len(parts)-1; i++ {
		for j := i + 1; j < len(parts); j++ {
			if isBadPair(ordering, parts[i], parts[j]) {
				if autoFix {
					parts[i], parts[j] = parts[j], parts[i]
				}
				return true
			}
		}
	}
	return false
}

func fixOrdering(parts []int, ordering map[int][]int) {
	fixNeeded := true
	for fixNeeded {
		fixNeeded = isFixNeeded(parts, ordering, true)
	}
}

func toNum(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func main() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	ordering := map[int][]int{}
	var updates [][]int

	orderProcess := true
	for _, l := range lines {
		if orderProcess && l == "" {
			orderProcess = false
			continue
		}
		if orderProcess {
			parts := strings.Split(l, "|")
			if len(parts) != 2 {
				panic("bad ordering: " + l)
			}
			l, r := toNum(parts[0]), toNum(parts[1])
			ordering[l] = append(ordering[l], r)
		} else {
			strs := strings.Split(l, ",")
			var nums []int
			for _, str := range strs {
				nums = append(nums, toNum(str))
			}
			updates = append(updates, nums)
		}
	}

	var badOrdering [][]int
	total := 0
	for _, update := range updates {
		if isFixNeeded(update, ordering, false) {
			badOrdering = append(badOrdering, update)
			continue
		}
		total += midNumber(update)
	}
	log.Println("mid num total:", total)

	total = 0
	for _, update := range badOrdering {
		fixOrdering(update, ordering)
		total += midNumber(update)
	}
	log.Println("fixed mid num total:", total)
}
