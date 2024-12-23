package main

import (
	_ "embed"
	"gonum.org/v1/gonum/stat/combin"
	"log"
	"sort"
	"strings"
)

//go:embed input.txt
var input string

type pair struct {
	a, b string
}

type trio struct {
	a, b, c string
}

func toCollaboratorMap(pairs []pair) map[string]map[string]bool {
	ret := map[string]map[string]bool{}
	for _, p := range pairs {
		if _, ok := ret[p.a]; !ok {
			ret[p.a] = map[string]bool{}
		}
		if _, ok := ret[p.b]; !ok {
			ret[p.b] = map[string]bool{}
		}
		ret[p.a][p.b] = true
		ret[p.b][p.a] = true
	}
	return ret
}

func ncr(candidates []string, r int, collaborators map[string]map[string]bool) [][]string {
	var ret [][]string
	areConnected := func(people []string) bool {
		for _, p1 := range people {
			for _, p2 := range people {
				if p1 == p2 {
					continue
				}
				if !collaborators[p1][p2] {
					return false
				}
			}
		}
		return true
	}
	n := len(candidates)
	combinations := combin.Combinations(n, r)
	for _, combination := range combinations {
		var people []string
		for _, i := range combination {
			people = append(people, candidates[i])
		}
		if areConnected(people) {
			var combo []string
			for _, p := range people {
				combo = append(combo, p)
			}
			sort.Strings(combo)
			ret = append(ret, combo)
		}
	}
	return ret
}

func collaboratorSets(collaborators map[string]map[string]bool) [][]string {
	var ret [][]string
	for person, c := range collaborators {
		set := []string{person}
		for p := range c {
			set = append(set, p)
		}
		ret = append(ret, set)
	}
	return ret
}

func part1(pairs []pair) {
	trios := map[trio]bool{}
	collaborators := toCollaboratorMap(pairs)
	sets := collaboratorSets(collaborators)
	for _, set := range sets {
		foundTrios := ncr(set, 3, collaborators)
		for _, t := range foundTrios {
			trios[trio{t[0], t[1], t[2]}] = true
		}
	}
	log.Println("TRIPLES:", len(trios))
	var ret2 []trio
	for t := range trios {
		if strings.HasPrefix(t.a, "t") || strings.HasPrefix(t.b, "t") || strings.HasPrefix(t.c, "t") {
			ret2 = append(ret2, t)
		}
	}
	log.Println("T TRIPLES:", len(ret2))
	log.Println("DISTINCT PEOPLE:", len(collaborators))
}

func part2(pairs []pair) {
	collaborators := toCollaboratorMap(pairs)
	sets := collaboratorSets(collaborators)
	maxSizeFound := 0
	var candidateWinner []string

outer:
	for _, candidates := range sets {
		r := len(candidates)
		for r > 2 {
			if r <= maxSizeFound {
				continue outer
			}
			combos := ncr(candidates, r, collaborators)
			if len(combos) > 0 {
				candidateWinner = combos[0]
				maxSizeFound = r
			}
			r--
		}
	}
	log.Println("MSF:", maxSizeFound)
	log.Println(strings.Join(candidateWinner, ","))
}

func main() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var pairs []pair
	for _, line := range lines {
		parts := strings.Split(line, "-")
		if parts[0] > parts[1] {
			parts[0], parts[1] = parts[1], parts[0]
		}
		pairs = append(pairs, pair{parts[0], parts[1]})
	}
	part1(pairs)
	part2(pairs)
}
