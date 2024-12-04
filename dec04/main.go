package main

import (
	_ "embed"
	"log"
	"strings"
)

//go:embed input.txt
var input string

var test = "XMAS"

var lines = strings.Split(strings.TrimSpace(input), "\n")
var maxCols = len(lines[0])
var maxRows = len(lines)

func isChristmas(lines []string, row, col int, rowOffset, colOffset int) int {
	for testIndex := 0; testIndex < len(test); testIndex++ {
		if row < 0 || row >= maxRows || col < 0 || col >= maxCols {
			return 0
		}
		ch := test[testIndex]
		if lines[row][col] != ch {
			return 0
		}
		row += rowOffset
		col += colOffset
	}
	return 1
}

func isMAS(l []string, i int, j int) int {
	ch := l[i][j]
	if ch != 'A' {
		return 0
	}
	if !((l[i-1][j-1] == 'M' && l[i+1][j+1] == 'S') ||
		(l[i-1][j-1] == 'S' && l[i+1][j+1] == 'M')) {
		return 0
	}
	if !((l[i-1][j+1] == 'M' && l[i+1][j-1] == 'S') ||
		(l[i-1][j+1] == 'S' && l[i+1][j-1] == 'M')) {
		return 0
	}
	return 1
}

func main() {
	count := 0
	for i := 0; i < maxRows; i++ {
		for j := 0; j < maxCols; j++ {
			count += isChristmas(lines, i, j, 1, 0)
			count += isChristmas(lines, i, j, 1, 1)
			count += isChristmas(lines, i, j, 1, -1)
			count += isChristmas(lines, i, j, -1, 0)
			count += isChristmas(lines, i, j, -1, -1)
			count += isChristmas(lines, i, j, -1, 1)
			count += isChristmas(lines, i, j, 0, -1)
			count += isChristmas(lines, i, j, 0, 1)
		}
	}
	log.Println("COUNT=", count)

	count = 0
	for i := 1; i < maxRows-1; i++ {
		for j := 1; j < maxCols-1; j++ {
			count += isMAS(lines, i, j)
		}
	}
	log.Println("COUNT MAS=", count)
}
