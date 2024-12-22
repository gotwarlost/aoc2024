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

func mix(first, second int) int {
	return first ^ second
}

func prune(n int) int {
	return n % 16777216
}

func nextSecret(secret int) int {
	t1 := prune(mix(secret, secret*64))
	t2 := prune(mix(t1, t1/32))
	t3 := prune(mix(t2, t2*2048))
	return t3
}

func toNum(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

type change struct {
	payoff int
	delta  int
}

type partition struct {
	a, b, c, d int
}

type partitionSum struct {
	p   partition
	sum int
}

func main() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var nums []int
	for _, line := range lines {
		nums = append(nums, toNum(line))
	}
	secret2K := func(s int) (int, []change) {
		var ret []change
		for i := 0; i < 2000; i++ {
			orig := s
			s = nextSecret(s)
			oldPayoff := orig % 10
			payoff := s % 10
			c := change{payoff: payoff, delta: payoff - oldPayoff}
			ret = append(ret, c)
		}
		return s, ret
	}
	sum := 0

	partitionSuperset := map[partition]bool{}
	var partitionsByBuyer []map[partition]int

	computePartitions := func(changes []change) map[partition]int {
		a, b, c, d := -1, changes[0].delta, changes[1].delta, changes[2].delta
		ret := map[partition]int{}
		for i := 3; i < len(changes); i++ {
			a, b, c = b, c, d
			d = changes[i].delta
			p := partition{a, b, c, d}
			if _, ok := ret[p]; ok {
				continue
			}
			ret[p] = changes[i].payoff
			partitionSuperset[p] = true
		}
		return ret
	}
	for _, n := range nums {
		val, changes := secret2K(n)
		sum += val
		partitionsByBuyer = append(partitionsByBuyer, computePartitions(changes))
	}
	log.Println("SUM:", sum)

	var psets []partitionSum
	for p := range partitionSuperset {
		sum := 0
		for _, pcount := range partitionsByBuyer {
			sum += pcount[p]
		}
		psets = append(psets, partitionSum{p: p, sum: sum})
	}
	sort.Slice(psets, func(i, j int) bool {
		return psets[i].sum > psets[j].sum
	})
	log.Println(psets[0])
}
