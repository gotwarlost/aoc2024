package main

import (
	"fmt"
	"sort"
	"strings"
)

type operator string

const (
	XOR operator = "XOR"
	OR  operator = "OR"
	AND operator = "AND"
)

func (o operator) apply(leftValue, rightValue int) int {
	switch o {
	case XOR:
		return leftValue ^ rightValue
	case OR:
		return leftValue | rightValue
	case AND:
		return leftValue & rightValue
	default:
		panic("invalid operator:" + o)
	}
}

func (o operator) String() string { return string(o) }

func newOperator(s string) operator {
	switch s {
	case string(XOR):
		return XOR
	case string(AND):
		return AND
	case string(OR):
		return OR
	default:
		panic("invalid operator:" + s)
	}
}

type gate struct {
	left, right string
	op          operator
	out         string
}

func (g gate) String() string {
	return fmt.Sprintf("%s %s %s", g.left, g.op, g.right)
}

type signature string

const (
	sigUnknown   signature = "unknown"
	sigOutput    signature = "output"
	sigCarry     signature = "carry"
	sigOutXor    signature = "outXOR"
	sigOutAnd    signature = "outAND"
	sigLastCarry signature = "lastCarry"
)

func (s signature) String() string { return string(s) }

type varOps struct {
	name    string
	input   gate
	outputs []gate
}

func (v *varOps) init() {
	sort.Slice(v.outputs, func(i, j int) bool {
		return v.outputs[i].op < v.outputs[j].op
	})
}

func (v *varOps) sigString() string {
	strs := []string{v.input.op.String() + ":in"}
	for _, o := range v.outputs {
		strs = append(strs, o.op.String())
	}
	return strings.Join(strs, " ")
}

func (v *varOps) signature() signature {
	switch v.sigString() {
	case "XOR:in AND XOR":
		return sigOutXor
	case "AND:in OR":
		return sigOutAnd
	case "XOR:in":
		return sigOutput
	case "OR:in AND XOR":
		return sigCarry
	case "OR:in":
		return sigLastCarry
	}
	return sigUnknown
}

func (v *varOps) matchesSignature(sig signature) bool {
	return sig == v.signature()
}
