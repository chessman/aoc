//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type P struct {
	x int
	y int
}

func parse() ([]P, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var ps []P
	for {
		if !scanner.Scan() {
			return ps, scanner.Err()
		}

		var p P
		if _, err := fmt.Sscanf(scanner.Text(), "%d,%d", &p.x, &p.y); err != nil {
			return nil, err
		}

		ps = append(ps, p)
	}
}

func initScores(m [][]int, start P) [][]int {
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

func directionsFrom(m [][]int, p P) []P {
	var ps []P
	if p.x-1 >= 0 {
		ps = append(ps, P{p.x - 1, p.y})
	}
	if p.x+1 < len(m[0]) {
		ps = append(ps, P{p.x + 1, p.y})
	}
	if p.y-1 >= 0 {
		ps = append(ps, P{p.x, p.y - 1})
	}
	if p.y+1 < len(m) {
		ps = append(ps, P{p.x, p.y + 1})
	}
	return ps
}

func buildScoreMap(m [][]int, cur P, scores [][]int, canStep func([][]int, P) bool) [][]int {
	possibleMoves := directionsFrom(m, cur)

	for _, p := range possibleMoves {
		if !canStep(m, p) {
			continue
		}

		newScore := min(scores[p.y][p.x], scores[cur.y][cur.x]+1)
		if newScore < scores[p.y][p.x] {
			scores[p.y][p.x] = newScore
			buildScoreMap(m, p, scores, canStep)
		}
	}

	return scores
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

func createMap(sizeX, sizeY int, corrupted []P) [][]int {
	var m [][]int
	for y := 0; y < sizeY; y++ {
		m = append(m, make([]int, sizeX))
	}
	for i, p := range corrupted {
		m[p.y][p.x] = i + 1
	}
	return m
}

func reachable(m [][]int, corruptionSize int) bool {
	canStep := func(m [][]int, p P) bool { return m[p.y][p.x] == 0 || m[p.y][p.x] > corruptionSize }
	start := P{0, 0}
	scoreMap := buildScoreMap(m, start, initScores(m, start), canStep)
        return scoreMap[len(m) - 1][len(m[0]) - 1] != math.MaxInt32
}

func firstUnreachable(m [][]int, corrupted []P) string {
	minSize := 1024
	maxSize := len(corrupted)

	for {
		if maxSize - minSize <= 1 {
			return fmt.Sprintf("%d,%d", corrupted[maxSize - 1].x, corrupted[maxSize - 1].y)
		}
		curSize := minSize + (maxSize - minSize) / 2
		if reachable(m, curSize) {
			minSize = curSize
		} else {
			maxSize = curSize
		}
	}
}

func main() {
	corrupted, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	sizeX := 71
	sizeY := 71
	start := P{0, 0}
	end := P{sizeX - 1, sizeY - 1}
	corruptionSize := 1024
	m := createMap(sizeX, sizeY, corrupted)
	canStep := func(m [][]int, p P) bool { return m[p.y][p.x] == 0 || m[p.y][p.x] > corruptionSize }
	scoreMap := buildScoreMap(m, start, initScores(m, start), canStep)
	fmt.Println("part 1:", scoreMap[end.y][end.x])
	fmt.Println("part 2:", firstUnreachable(m, corrupted))
}
