package main

import (
	_ "embed"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type register int

const (
	regA register = iota
	regB
	regC
)

type opcode int

const (
	adv opcode = iota
	bxl
	bst
	jnz
	bxc
	out
	bdv
	cdv
)

func (o opcode) String() string {
	switch o {
	case adv:
		return "adv"
	case bdv:
		return "bdv"
	case cdv:
		return "cdv"
	case bxl:
		return "bxl"
	case bst:
		return "bst"
	case bxc:
		return "bxc"
	case jnz:
		return "jnz"
	case out:
		return "out"
	}
	return "unk"
}

func comboOperand(val int, registers []int) int {
	switch val {
	case 0, 1, 2, 3:
		return val
	case 4:
		return registers[regA]
	case 5:
		return registers[regB]
	case 6:
		return registers[regC]
	default:
		panic(fmt.Sprintf("Invalid combo operand: %d", val))
	}
}

var opDiv = func(operand int, registers []int, storeReg register) (r result) {
	numerator := registers[regA]
	denominator := int(math.Pow(2, float64(comboOperand(operand, registers))))
	registers[storeReg] = numerator / denominator
	return
}

type result struct {
	out *int
	ip  *int
}

type operator func(operand int, registers []int) result

var ops = map[opcode]operator{
	adv: func(operand int, registers []int) result {
		return opDiv(operand, registers, regA)
	},
	bdv: func(operand int, registers []int) result {
		return opDiv(operand, registers, regB)
	},
	cdv: func(operand int, registers []int) result {
		return opDiv(operand, registers, regC)
	},
	bxl: func(operand int, registers []int) (r result) {
		registers[regB] = registers[regB] ^ operand
		return
	},
	bxc: func(operand int, registers []int) (r result) {
		registers[regB] = registers[regB] ^ registers[regC]
		return
	},
	bst: func(operand int, registers []int) (r result) {
		registers[regB] = comboOperand(operand, registers) % 8
		return
	},
	jnz: func(operand int, registers []int) (r result) {
		if registers[regA] == 0 {
			return r
		}
		return result{ip: &operand}
	},
	out: func(operand int, registers []int) result {
		out := comboOperand(operand, registers) % 8
		return result{out: &out}
	},
}

type puzzle struct {
	registers    []int
	instructions []int
	instStr      string
}

func toNum(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func parse(s string) *puzzle {
	ret := &puzzle{
		registers: []int{0, 0, 0},
	}
	trimNum := func(s string, prefix string) int {
		s = strings.TrimPrefix(s, prefix)
		return toNum(s)
	}
	ra, rb, rc, prog := "Register A: ", "Register B: ", "Register C: ", "Program: "
	for _, l := range strings.Split(strings.TrimSpace(s), "\n") {
		switch {
		case strings.HasPrefix(l, ra):
			ret.registers[regA] = trimNum(l, ra)
		case strings.HasPrefix(l, rb):
			ret.registers[regB] = trimNum(l, rb)
		case strings.HasPrefix(l, rc):
			ret.registers[regC] = trimNum(l, rc)
		case strings.HasPrefix(l, prog):
			l = strings.TrimPrefix(l, prog)
			ret.instStr = l
			parts := strings.Split(l, ",")
			if len(parts)%2 != 0 {
				panic("odd parts")
			}
			for i := 0; i < len(parts); i++ {
				ret.instructions = append(ret.instructions, toNum(parts[i]))
			}
		}
	}
	return ret
}

func (p *puzzle) solve(part2 bool) string {
	ip := 0
	var output []string
	var intout []int
	for ip < len(p.instructions) {
		code := opcode(p.instructions[ip])
		if code > 7 {
			if part2 {
				return ""
			}
			panic("opcode out of bounds")
		}
		operand := p.instructions[ip+1]
		res := ops[code](operand, p.registers)
		if res.out != nil {
			x := *res.out
			intout = append(intout, x)
			if part2 {
				if intout[len(intout)-1] != p.instructions[len(intout)-1] {
					return ""
				}
			}
			output = append(output, fmt.Sprint(x))
		}
		if res.ip != nil {
			ip = *res.ip
			if ip < 0 || ip > len(p.instructions)-1 {
				if part2 {
					return ""
				}
				panic(fmt.Sprintf("bad ip: %d", ip))
			}
		} else {
			ip += 2
		}
	}
	return strings.Join(output, ",")
}

func main() {
	puz := parse(input)
	output := puz.solve(false)
	log.Println("OUTPUT:", output)
}