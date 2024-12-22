package main

import (
	_ "embed"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func numericKeypadValueAt(row, col int) string {
	if col < 0 || col > 2 {
		panic("invalid col")
	}
	switch row {
	case 0:
		return fmt.Sprintf("%d", 7+col)
	case 1:
		return fmt.Sprintf("%d", 4+col)
	case 2:
		return fmt.Sprintf("%d", 1+col)
	case 3:
		switch col {
		case 0:
			return ""
		case 1:
			return "0"
		default:
			return "A"
		}
	default:
		panic("bad row")
	}
}

type point struct {
	row, col int
}

func (p point) manhattanDistance(other point) int {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	return abs(other.row-p.row) + abs(other.col-p.col)
}

type pair struct {
	a, b string
}

func reversePath(e string) string {
	newPath := make([]byte, 0, len(e))
	for i := len(e) - 1; i >= 0; i-- {
		c := e[i]
		switch c {
		case '<':
			newPath = append(newPath, '>')
		case '>':
			newPath = append(newPath, '<')
		case '^':
			newPath = append(newPath, 'v')
		case 'v':
			newPath = append(newPath, '^')
		}
	}
	return string(newPath)
}

func setupKeyboard(valueMap map[point]string, maxRows, maxCols int) map[pair][]string {
	shortestPaths := map[pair][]string{}

	var setShortestPathsBetween func(p1, p2 point) []string
	setShortestPathsBetween = func(p1, p2 point) (ret []string) {
		v1, v2 := valueMap[p1], valueMap[p2]
		if v1 == "" || v2 == "" {
			return nil
		}
		ptPair := pair{v1, v2}
		p, ok := shortestPaths[ptPair]
		if ok {
			return p
		}

		defer func() {
			shortestPaths[ptPair] = ret
			var paths []string
			for _, e := range ret {
				paths = append(paths, reversePath(e))
			}
			shortestPaths[pair{v2, v1}] = paths
		}()

		switch {
		case p1.row == p2.row && p2.col == p1.col+1:
			return []string{">"}
		case p1.row == p2.row && p2.col == p1.col-1:
			return []string{"<"}
		case p1.col == p2.col && p2.row == p1.row+1:
			return []string{"v"}
		case p1.col == p2.col && p2.row == p1.row-1:
			return []string{"^"}
		}

		appendDistance := func(pt point, dir string) []string {
			path := setShortestPathsBetween(pt, p2)
			for _, p0 := range path {
				ret = append(ret, fmt.Sprintf("%s%s", dir, p0))
			}
			return ret
		}
		if p1.col != p2.col {
			p3 := point{p1.row, p1.col + 1}
			dir := ">"
			if p2.col < p1.col {
				p3 = point{p1.row, p1.col - 1}
				dir = "<"
			}
			appendDistance(p3, dir)
		}
		if p1.row != p2.row {
			p4 := point{p1.row + 1, p1.col}
			dir := "v"
			if p2.row < p1.row {
				p4 = point{p1.row - 1, p1.col}
				dir = "^"
			}
			appendDistance(p4, dir)
		}
		return ret
	}

	setShortestPathsFrom := func(startRow, startCol int) {
		ptStart := point{startRow, startCol}
		if valueMap[ptStart] == "" {
			return
		}
		for row := 0; row < maxRows; row++ {
			for col := 0; col < maxCols; col++ {
				if row == startRow && col == startCol {
					continue
				}
				pt := point{row, col}
				nextValue := valueMap[pt]
				if nextValue == "" {
					continue
				}
				setShortestPathsBetween(ptStart, pt)
			}
		}
	}

	for row := 0; row < maxRows; row++ {
		for col := 0; col < maxCols; col++ {
			setShortestPathsFrom(row, col)
		}
	}
	return shortestPaths
}

// debugging stuff
func printShortestPaths(shortestPaths map[pair][]string) {
	var pairs []pair
	for k := range shortestPaths {
		pairs = append(pairs, k)
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].a == pairs[j].a {
			return pairs[i].b < pairs[j].b
		}
		return pairs[i].a < pairs[j].a
	})
	for _, k := range pairs {
		v := shortestPaths[k]
		log.Println(k.a, k.b, strings.Join(v, ", "))
	}
}

type keyboard struct {
	startPoint    point
	pointMap      map[string]point
	valueMap      map[point]string
	shortestPaths map[pair][]string
}

func setup() *puzzle {
	numberKeyboardMap := map[point]string{}
	for row := 0; row < 4; row++ {
		for col := 0; col < 3; col++ {
			v := numericKeypadValueAt(row, col)
			pt := point{row: row, col: col}
			numberKeyboardMap[pt] = v
		}
	}
	numericKeyboard := setupKeyboard(numberKeyboardMap, 4, 3)
	//log.Println("num keyboard, len=", len(numericKeyboard))
	//printShortestPaths(numericKeyboard)

	dirKeyboardMap := map[point]string{
		point{0, 1}: "^",
		point{0, 2}: "A",
		point{1, 0}: "<",
		point{1, 1}: "v",
		point{1, 2}: ">",
	}
	dirKeyboard := setupKeyboard(dirKeyboardMap, 2, 3)
	//log.Println("dir keyboard, len=", len(dirKeyboard))
	//printShortestPaths(dirKeyboard)

	invert := func(m map[point]string) map[string]point {
		ret := map[string]point{}
		for k, v := range m {
			ret[v] = k
		}
		return ret
	}

	return &puzzle{
		num: &keyboard{
			startPoint:    point{3, 2},
			valueMap:      numberKeyboardMap,
			pointMap:      invert(numberKeyboardMap),
			shortestPaths: numericKeyboard,
		},
		dir: &keyboard{
			startPoint:    point{0, 2},
			pointMap:      invert(dirKeyboardMap),
			valueMap:      dirKeyboardMap,
			shortestPaths: dirKeyboard,
		},
	}
}

type puzzle struct {
	num *keyboard
	dir *keyboard
}

func (z *puzzle) traversePointsStep1(values []string) []string {
	head := values[0]
	rest := values[1:]
	paths := z.num.shortestPaths[pair{head, rest[0]}]
	var ret []string
	if len(rest) == 1 {
		for _, p := range paths {
			ret = append(ret, p+"A")
		}
	} else {
		remaining := z.traversePointsStep1(rest)
		for _, p := range paths {
			for _, a := range remaining {
				ret = append(ret, p+"A"+a)
			}
		}
	}
	return ret
}

func (z *puzzle) traversePointsDir(values []string) []string {
	head := values[0]
	rest := values[1:]
	paths := z.dir.shortestPaths[pair{head, rest[0]}]
	if len(paths) == 0 {
		paths = []string{""}
	}
	var ret []string
	if len(rest) == 1 {
		for _, p := range paths {
			ret = append(ret, p+"A")
		}
	} else {
		remaining := z.traversePointsDir(rest)
		for _, p := range paths {
			for _, a := range remaining {
				ret = append(ret, p+"A"+a)
			}
		}
	}
	return ret
}

func sortByLength(values []string) {
	sort.Slice(values, func(i, j int) bool {
		if len(values[i]) == len(values[j]) {
			return values[i] < values[j]
		}
		return len(values[i]) < len(values[j])
	})
}

func toValues(str string) []string {
	values := []string{"A"}
	for i := 0; i < len(str); i++ {
		s := str[i : i+1]
		values = append(values, s)
	}
	return values
}

func debugCandidates(cs []string) {
	for _, c := range cs {
		log.Printf("%3d: %s", len(c), c)
	}
}

func (z *puzzle) findShortestForCode(code string) int {
	// step 1: expand possibilities for code
	values := toValues(code)

	log.Println("CODE:", code)
	step1Candidates := z.traversePointsStep1(values)
	sortByLength(step1Candidates)

	//log.Println("step 1")
	//debugCandidates(step1Candidates)

	//log.Println("step 2")
	var step2Candidates []string
	for _, s := range step1Candidates {
		candidates := z.traversePointsDir(toValues(s))
		step2Candidates = append(step2Candidates, candidates...)
	}
	sortByLength(step2Candidates)
	//debugCandidates(step2Candidates)

	log.Println("step 3")
	var step3Candidates []string
	for _, s := range step2Candidates {
		candidates := z.traversePointsDir(toValues(s))
		step3Candidates = append(step3Candidates, candidates...)
	}
	sortByLength(step3Candidates)
	//debugCandidates(step3Candidates)
	log.Println("SHORTEST:", len(step3Candidates[0]))
	return len(step3Candidates[0])
}

func main() {
	puz := setup()
	codes := strings.Split(strings.TrimSpace(input), "\n")
	sum := 0
	var vals []int
	for _, code := range codes {
		shortestSteps := puz.findShortestForCode(code)
		code0 := strings.TrimSuffix(code, "A")
		n, err := strconv.Atoi(code0)
		if err != nil {
			panic(err)
		}
		val := shortestSteps * n
		vals = append(vals, val)
		sum += val
	}
	log.Println("SUM:", sum)
}
