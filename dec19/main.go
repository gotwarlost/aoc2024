package main

import (
	_ "embed"
	"log"
	"sort"
	"strings"
)

var (
	//go:embed base.txt
	base string
	//go:embed input.txt
	input string
)

type puzzle struct {
	stripes        []string
	patterns       []string
	substringCount map[string]int
}

func parse(s string) *puzzle {
	pparts := strings.Split(strings.TrimSpace(s), "\n\n")
	stripeParts := strings.Split(pparts[0], ",")
	ret := &puzzle{substringCount: map[string]int{}}
	for _, sp := range stripeParts {
		ret.stripes = append(ret.stripes, strings.TrimSpace(sp))
	}
	for _, line := range strings.Split(pparts[1], "\n") {
		ret.patterns = append(ret.patterns, line)
	}
	return ret
}

func (p *puzzle) candidateStripes(input string) []string {
	var ret []string
	for _, s := range p.stripes {
		if strings.HasPrefix(input, s) {
			ret = append(ret, s)
		}
	}
	sort.Slice(ret, func(i, j int) bool {
		return len(ret[i]) > len(ret[j])
	})
	return ret
}

func (p *puzzle) countPossibilities(input string) (ret int) {
	if count, ok := p.substringCount[input]; ok {
		return count
	}
	defer func() {
		p.substringCount[input] = ret
	}()
	ret = 0
	candidates := p.candidateStripes(input)
	for _, c := range candidates {
		next := input[len(c):]
		if len(next) == 0 {
			ret++
			continue
		}
		children := p.countPossibilities(next)
		ret += children
	}
	return ret
}

func (p *puzzle) solve() (totalPossibe int, numPossibilities int) {
	for _, pat := range p.patterns {
		count := p.countPossibilities(pat)
		if count > 0 {
			totalPossibe++
		}
		numPossibilities += count
	}
	return
}

func (p *puzzle) countPossible() int {
	ret := 0
	for _, pat := range p.patterns {
		if p.countPossibilities(pat) > 0 {
			ret++
		}
	}
	return ret
}

func main() {
	puz := parse(input)
	t, a := puz.solve()
	log.Println("POSSIBLE:", t)
	log.Println("ALL POSSIBILITIES:", a)
}
