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

func toNum(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

type puzzle struct {
	initialValues map[string]int
	values        map[string]int
	vars          map[string]*varOps
	gatesByExpr   map[string]gate
	varOps        map[string][]gate
	subst         map[string]string
}

func makeExpr(in1 string, op operator, in2 string) string {
	if in1 > in2 {
		in1, in2 = in2, in1
	}
	return fmt.Sprintf("%s %s %s", in1, op, in2)
}

func newPuzzle() *puzzle {
	ret := &puzzle{
		initialValues: map[string]int{},
		values:        map[string]int{},
		vars:          map[string]*varOps{},
		gatesByExpr:   map[string]gate{},
	}

	lines := strings.Split(strings.TrimSpace(input), "\n")
	assign := true
	for _, l := range lines {
		if assign {
			if l == "" {
				assign = false
				continue
			}
			parts := strings.Split(l, ": ")
			ret.initialValues[parts[0]] = toNum(parts[1])
			continue
		}
		parts := strings.Split(l, " ")
		if len(parts) != 5 {
			panic("bad line: " + l)
		}
		g := gate{
			left:  parts[0],
			op:    newOperator(parts[1]),
			right: parts[2],
			out:   parts[4],
		}
		if g.left > g.right {
			g.left, g.right = g.right, g.left
		}

		getVarOps := func(name string) *varOps {
			vo := ret.vars[name]
			if vo == nil {
				vo = &varOps{}
				ret.vars[name] = vo
			}
			return vo
		}

		ret.gatesByExpr[makeExpr(g.left, g.op, g.right)] = g
		getVarOps(g.out).input = g
		vo := getVarOps(g.left)
		vo.outputs = append(vo.outputs, g)
		vo = getVarOps(g.right)
		vo.outputs = append(vo.outputs, g)
	}
	for _, vo := range ret.vars {
		vo.init()
	}
	return ret
}

func (z *puzzle) init() {
	z.values = map[string]int{}
	for k, v := range z.initialValues {
		z.values[k] = v
	}
	z.subst = map[string]string{}
}

func (z *puzzle) getValueOf(x string) (ret int) {
	v, ok := z.values[x]
	if ok {
		return v
	}
	defer func() {
		z.values[x] = ret
	}()
	vo, ok := z.vars[x]
	if !ok {
		panic("internal error, no variable:" + x)
	}
	g := vo.input
	l := z.getValueOf(g.left)
	r := z.getValueOf(g.right)
	return g.op.apply(l, r)
}

func (z *puzzle) valueFromBits(prefix string, num int) int {
	output := 0
	for i := num - 1; i >= 0; i-- {
		variable := fmt.Sprintf("%s%02d", prefix, i)
		v := z.getValueOf(variable)
		output |= v
		if i > 0 {
			output = output << 1
		}
	}
	return output
}

func (z *puzzle) part1() (output int, numZs int) {
	z.init()
	var zs []string
	for k := range z.vars {
		if strings.HasPrefix(k, "z") {
			zs = append(zs, k)
		}
	}
	sort.Strings(zs)
	return z.valueFromBits("z", len(zs)), len(zs)
}

func nthBit(num int, n int) int {
	out := (1 << n) ^ num
	if out != 0 {
		return 1
	}
	return 0
}

type bitAdder struct {
	in1, in2    string // left, right input variables
	inputCarry  string // carry input, may be blank for the first bit
	outXOR      string // in1 ^ in2
	outAND      string // in1 & in2
	outSum      string // outXOR ^ inputCarry
	carryAND    string // outXOR & inputCarry
	outputCarry string // outAND | carryAND
}

func (z *puzzle) findExprOutput(in1 string, op operator, in2 string) (string, bool) {
	if s1, ok := z.subst[in1]; ok {
		in1 = s1
	}
	if s2, ok := z.subst[in2]; ok {
		in2 = s2
	}
	e := makeExpr(in1, op, in2)
	g, ok := z.gatesByExpr[e]
	if !ok {
		return "", false
	}
	return g.out, true
}

func (z *puzzle) getVarOps(name string) *varOps {
	vo := z.vars[name]
	if vo == nil {
		panic("no ops for:" + name)
	}
	return vo
}

func (z *puzzle) newAdder(i int, carry string) (ret bitAdder) {
	//log.Println("adder:", i, carry)
	defer func() {
		//log.Printf("adder: %+v", ret)
	}()

	in1, in2 := fmt.Sprintf("x%02d", i), fmt.Sprintf("y%02d", i)
	b := bitAdder{
		in1:        in1,
		in2:        in2,
		inputCarry: carry,
	}
	b.outXOR, _ = z.findExprOutput(in1, XOR, in2)
	if !z.getVarOps(b.outXOR).matchesSignature(sigOutXor) {
		log.Println("out XOR signature mismatch for:", i, b.outXOR)
	}
	b.outAND, _ = z.findExprOutput(in1, AND, in2)
	if !z.getVarOps(b.outAND).matchesSignature(sigOutAnd) {
		log.Println("out AND signature mismatch for:", i, b.outAND)
	}

	// no input carry for the first
	if i == 0 {
		b.outSum = b.outXOR
		b.carryAND = ""
		b.outputCarry = b.outAND
	} else {
		var ok bool
		b.outSum, ok = z.findExprOutput(b.outXOR, XOR, b.inputCarry)
		if !ok {
			log.Printf("no out sum for %d: %s", i, makeExpr(b.outXOR, XOR, b.inputCarry))
		} else {
			if !z.getVarOps(b.outSum).matchesSignature(sigOutput) {
				log.Printf("out sum signature mismatch for %d: %s", i, b.outSum)
			}
		}
		b.carryAND, ok = z.findExprOutput(b.outXOR, AND, b.inputCarry)
		if !ok {
			log.Printf("no carry-and for %d: %s", i, makeExpr(b.outXOR, AND, b.inputCarry))
		} else {
			if !z.getVarOps(b.carryAND).matchesSignature(sigOutAnd) {
				log.Printf("carry-and signature mismatch for %d: %s", i, b.carryAND)
			}
		}
		b.outputCarry, ok = z.findExprOutput(b.outAND, OR, b.carryAND)
		if !ok {
			log.Printf("no output carry for %d: %s", i, makeExpr(b.outAND, OR, b.carryAND))
		} else {
			if !z.getVarOps(b.outputCarry).matchesSignature(sigCarry) {
				log.Printf("carry signature mismatch for %d: %s", i, b.outputCarry)
			}
		}
	}
	if b.outSum != fmt.Sprintf("z%02d", i) {
		log.Printf("unexpected output for %d: %s", i, b.outSum)
	}
	return b
}

func (z *puzzle) part2(numZs int) {
	z.init()
	x := z.valueFromBits("x", numZs-1)
	y := z.valueFromBits("y", numZs-1)
	actual := z.valueFromBits("z", numZs)
	expected := x + y

	log.Printf("Expected: %d %b", expected, expected)
	log.Printf("  Actual: %d %b", actual, actual)
	log.Printf("    diff: %14d %046b", actual^expected, actual^expected)

	var carry string
	for i := 0; i < numZs-1; i++ {
		adder := z.newAdder(i, carry)
		carry = adder.outputCarry
	}
}

func main() {
	puz := newPuzzle()
	out, numZs := puz.part1()
	log.Println("PART 1:", out)
	puz.part2(numZs)
}
