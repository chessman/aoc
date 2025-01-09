//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
)

func parse() (map[string]map[string]struct{}, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	conns := make(map[string]map[string]struct{})
	for {
		if !scanner.Scan() {
			return conns, scanner.Err()
		}

		s := scanner.Text()
		left := s[0:2]
		right := s[3:5]

		if conns[left] == nil {
			conns[left] = make(map[string]struct{})
		}
		if conns[right] == nil {
			conns[right] = make(map[string]struct{})
		}
		conns[left][right] = struct{}{}
		conns[right][left] = struct{}{}
	}
}

func uniqueSorted(rels map[string]map[string]struct{}) []string {
	uniq := map[string]struct{}{}
	for k := range rels {
		uniq[k] = struct{}{}
	}

	return slices.Sorted(maps.Keys(uniq))
}

func check(c1, c2 string, conns map[string]map[string]struct{}) bool {
	m := conns[c1]
	if m == nil {
		return false
	}
	_, ok := m[c2]
	return ok
}

func countTriplets(conns map[string]map[string]struct{}) int {
	lefts := uniqueSorted(conns)

	count := 0
	for _, el := range lefts {
		if el[0] != 't' {
			continue
		}
		candidates := slices.Collect(maps.Keys(conns[el]))
		slices.Sort(candidates)

		for i := range len(candidates) - 1 {
			if candidates[i][0] == 't' && candidates[i] < el {
				continue
			}
			for j := i + 1; j < len(candidates); j++ {
				if candidates[j][0] == 't' && candidates[j] < el {
					continue
				}
				if check(candidates[i], candidates[j], conns) {
					count++
				}
			}
		}

	}
	return count
}

func intersection(vs []string, edges map[string]struct{}) []string {
	var intersect []string
	for _, v := range vs {
		if _, ok := edges[v]; ok {
			intersect = append(intersect, v)
		}
	}
	return intersect
}

func substract(vs []string, edges map[string]struct{}) []string {
	var res []string
	for _, v := range vs {
		if _, ok := edges[v]; !ok {
			res = append(res, v)
		}
	}
	return res
}

func bronKerbosch(conns map[string]map[string]struct{}, cur []string, todo []string, done []string, cliques *[][]string) {
	if len(todo) == 0 && len(done) == 0 {
		*cliques = append(*cliques, slices.Clone(cur))
		return
	}

	for len(todo) > 0 {
		v := todo[0]
		newTodo := intersection(todo, conns[v])
		newDone := intersection(done, conns[v])
		bronKerbosch(conns, append(cur, v), newTodo, newDone, cliques)
		todo = todo[1:]
		done = append(done, v)
	}
}

func lanPartyPassword(conns map[string]map[string]struct{}) string {
	cliques := make([][]string, 0)
	lefts := uniqueSorted(conns)
	bronKerbosch(conns, nil, lefts, nil, &cliques)
	maxClique := slices.MaxFunc(cliques, func(a []string, b []string) int { return len(a) - len(b) })
	slices.Sort(maxClique)
	return strings.Join(maxClique, ",")
}

func main() {
	conns, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(countTriplets(conns))
	fmt.Println(lanPartyPassword(conns))
}
