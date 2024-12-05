package main

import (
	_ "embed"
	"log"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func midNumber(parts []string) int {
	midPoint := len(parts) / 2
	midNum, err := strconv.Atoi(parts[midPoint])
	if err != nil {
		panic(err)
	}
	return midNum
}

func isBadPair(ordering map[string][]string, left, right string) bool {
	after := ordering[right]
	for _, a := range after {
		if a == left {
			return true
		}
	}
	return false
}

func isFixNeeded(parts []string, ordering map[string][]string, autoFix bool) bool {
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

func fixOrdering(parts []string, ordering map[string][]string) {
	fixNeeded := true
	for fixNeeded {
		fixNeeded = isFixNeeded(parts, ordering, true)
	}
}

func main() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	ordering := map[string][]string{}
	var updates [][]string

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
			ordering[parts[0]] = append(ordering[parts[0]], parts[1])
		} else {
			nums := strings.Split(l, ",")
			updates = append(updates, nums)
		}
	}

	var badOrdering [][]string
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
