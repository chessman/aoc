//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Dir int

const (
	North Dir = iota
	East
	South
	West
)

type P struct {
	x   int
	y   int
	dir Dir
}

func (p P) withoutDir() P { p.dir = North; return p }

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

func startP(m [][]rune) P { return P{1, len(m) - 2, East} }
func endP(m [][]rune) P   { return P{len(m[0]) - 2, 1, East} }

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

var stepMap = map[Dir]func(P) P{
	North: func(p P) P { p.y--; p.dir = North; return p },
	South: func(p P) P { p.y++; p.dir = South; return p },
	East:  func(p P) P { p.x++; p.dir = East; return p },
	West:  func(p P) P { p.x--; p.dir = West; return p },
}

var allDirections = []Dir{East, West, North, South}

func directionsFrom(d Dir) []Dir {
	var dirs []Dir
	opposite := Dir((int(d) + 2) % len(allDirections))
	for _, d2 := range allDirections {
		if d2 != opposite {
			dirs = append(dirs, d2)
		}
	}

	return dirs
}

func buildScoreMap(m [][]rune, cur P, scores [][]int) [][]int {
	var possibleMoves []P
	for _, dir := range directionsFrom(cur.dir) {
		possibleMoves = append(possibleMoves, stepMap[dir](cur))
	}

	for _, p := range possibleMoves {
		if m[p.y][p.x] == '#' {
			continue
		}

		var newScore int
		if p.dir == cur.dir {
			newScore = min(scores[p.y][p.x], scores[cur.y][cur.x]+1)
		} else {
			newScore = min(scores[p.y][p.x], scores[cur.y][cur.x]+1001)
		}

		if newScore < scores[p.y][p.x] {
			scores[p.y][p.x] = newScore
			buildScoreMap(m, p, scores)
		}
	}

	return scores
}

func findAllTiles(m [][]rune, cur P, scores [][]int, paths map[P]struct{}, maxDist int) map[P]struct{} {
	var possibleMoves []P
	for _, dir := range directionsFrom(cur.dir) {
		possibleMoves = append(possibleMoves, stepMap[dir](cur))
	}

	for _, p := range possibleMoves {
		if m[p.y][p.x] == '#' {
			continue
		}

		var newMaxDist int
		if p.dir == cur.dir {
			newMaxDist = maxDist - 1
		} else {
			newMaxDist = maxDist - 1001
		}

		if scores[p.y][p.x] <= newMaxDist {
			paths[p.withoutDir()] = struct{}{}
			findAllTiles(m, p, scores, paths, newMaxDist)
		}
	}

	return paths
}

func printMap(m [][]rune, paths []P) {
	for _, p := range paths {
		m[p.y][p.x] = 'O'
	}
	for y := range m {
		for x := range m[y] {
			fmt.Printf("%c", m[y][x])
		}
		fmt.Println()
	}
}

func numberOfTiles(m [][]rune, scoreMap [][]int) int {
	end := endP(m)
	end.dir = South

	total := 0
	allPaths := findAllTiles(m, end, scoreMap, map[P]struct{}{end.withoutDir(): {}}, scoreMap[end.y][end.x])
	total += len(allPaths)

	end.dir = West
	allPaths = findAllTiles(m, end, scoreMap, map[P]struct{}{end.withoutDir(): {}}, scoreMap[end.y][end.x])
	total += len(allPaths) - 1

	return total
}

func main() {
	m, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	start := startP(m)
	scoreMap := buildScoreMap(m, start, initScores(m))

	total := numberOfTiles(m, scoreMap)
	fmt.Println(total)
}
