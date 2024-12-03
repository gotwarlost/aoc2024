package main

import (
	_ "embed"
	"log"
	"regexp"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

var (
	mulRE  = regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	fullRE = regexp.MustCompile(`do\(\)|don't\(\)|mul\((\d+),(\d+)\)`)
)

func main() {
	matches := mulRE.FindAllStringSubmatch(input, -1)
	result := 0
	for _, m := range matches {
		n1, _ := strconv.Atoi(m[1])
		n2, _ := strconv.Atoi(m[2])
		result += n1 * n2
	}
	log.Println("RESULT:", result)

	matches = fullRE.FindAllStringSubmatch(input, -1)
	result = 0
	enabled := true
	for _, m := range matches {
		switch {
		case strings.HasPrefix(m[0], "mul"):
			if enabled {
				n1, _ := strconv.Atoi(m[1])
				n2, _ := strconv.Atoi(m[2])
				result += n1 * n2
			}
		case strings.HasPrefix(m[0], "don't"):
			enabled = false
		default:
			enabled = true
		}
	}
	log.Println("FULL RESULT:", result)
}
