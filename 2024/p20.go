//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type P struct {
	x   int
	y   int
}

func parse() ([][]rune, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var m [][]rune
	for {
		if !scanner.Scan() {
			return m, scanner.Err()
		}

		m = append(m, []rune(scanner.Text()))
	}
}

func startP(m [][]rune) P {
	for y := range m {
		for x, el := range m[y] {
			if el == 'S' {
				return P{x, y}
			}
		}
	}
	return P{0, 0}
}

func endP(m [][]rune) P {
	for y := range m {
		for x, el := range m[y] {
			if el == 'E' {
				return P{x, y}
			}
		}
	}
	return P{0, 0}
}

func initScores(m [][]rune) [][]int {
	start := startP(m)
	scores := make([][]int, len(m))
	for y := 0; y < len(m); y++ {
		scores[y] = make([]int, len(m[0]))
		for x := 0; x < len(m[0]); x++ {
			scores[y][x] = math.MaxInt32
		}
	}
	scores[start.y][start.x] = 0

	return scores
}

func directionsFrom(maxX, maxY int, p P) []P {
	var ps []P
	if p.x-1 >= 0 {
		ps = append(ps, P{p.x - 1, p.y})
	}
	if p.x+1 < maxX {
		ps = append(ps, P{p.x + 1, p.y})
	}
	if p.y-1 >= 0 {
		ps = append(ps, P{p.x, p.y - 1})
	}
	if p.y+1 < maxY {
		ps = append(ps, P{p.x, p.y + 1})
	}
	return ps
}

func buildScoreMap(m [][]rune, cur P, scores [][]int) [][]int {
	for _, p := range directionsFrom(len(m[0]), len(m), cur) {
		if m[p.y][p.x] == '#' {
			continue
		}

		curValue := scores[cur.y][cur.x]
		newScore := min(scores[p.y][p.x], curValue+1)
		if newScore < scores[p.y][p.x] {
			scores[p.y][p.x] = newScore
			buildScoreMap(m, p, scores)
		}
	}

	return scores
}

func getPath(scores [][]int, start, end P) []P {
	if start == end {
		return []P{start}
	}

	var next P
	for _, p := range directionsFrom(len(scores[0]), len(scores), end) {
		if scores[end.y][end.x] - scores[p.y][p.x] == 1 {
			next = p
			break
		}
	}

	return append(getPath(scores, start, next), end)
}

func printMap(m [][]rune, paths []P) {
	for _, p := range paths {
		if p.y == -1 && p.x == -1 {
			continue
		}
		m[p.y][p.x] = 'X'
	}
	for y := range m {
		for x := range m[y] {
			fmt.Printf("%c", m[y][x])
		}
		fmt.Println()
	}
}

func abs(x int) int { if (x < 0) { return -x }; return x }

func countCheats(path []P, maxCheatSize, minBenefit int) int {
	count := 0
	for i := 0; i < len(path); i++ {
		for j := i + minBenefit; j < len(path); j++ {
			path1 := path[i]
			path2 := path[j]
			distance := abs(path1.x - path2.x) + abs(path1.y - path2.y)

			if distance <= maxCheatSize && j - i >= minBenefit + distance {
				count++
			}
		}
	}
	return count
}

func main() {
	m, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	start := startP(m)
	end := endP(m)
	scores := buildScoreMap(m, start, initScores(m))
	fmt.Println(scores[end.y][end.x])
	fmt.Println(start, end)
	path := getPath(scores, start, end)
	fmt.Println(countCheats(path, 2, 100))
	fmt.Println(countCheats(path, 20, 100))
}
