//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Coord struct {
	x int
	y int
}

func parse() ([][]int, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	var m [][]int

	for scanner.Scan() {
		line := scanner.Text()
		var row []int
		for _, char := range line {
			i, err := strconv.Atoi(string(char))
			if err != nil {
				return nil, err
			}
			row = append(row, i)
		}
		m = append(m, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

func findStarts(m [][]int) []Coord {
	var coords []Coord
	for y, row := range m {
		for x, v := range row {
			if v == 0 {
				coords = append(coords, Coord{x: x, y: y})
			}
		}
	}
	return coords
}

func climb(m [][]int, coords []Coord) []Coord {
	var res []Coord

	for _, p := range coords {
		v := m[p.y][p.x]
		if v == 9 {
			res = append(res, p)
		} else {
			var candidates []Coord
			if p.x - 1 >= 0 && m[p.y][p.x - 1] == v + 1 {
				candidates = append(candidates, Coord{x: p.x - 1, y: p.y})
			}

			if p.x + 1 < len(m[0]) && m[p.y][p.x + 1] == v + 1 {
				candidates = append(candidates, Coord{x: p.x + 1, y: p.y})
			}

			if p.y - 1 >= 0 && m[p.y - 1][p.x] == v + 1 {
				candidates = append(candidates, Coord{x: p.x, y: p.y - 1})
			}

			if p.y + 1 < len(m) && m[p.y + 1][p.x] == v + 1 {
				candidates = append(candidates, Coord{x: p.x, y: p.y + 1})
			}

			res = append(res, climb(m, candidates)...)
		}
	}

	return res
}

func countTrailheads(m [][]int, starts []Coord) int {
	var res int
	for _, start := range starts {
		trailheads := climb(m, []Coord{start})
		unique := map[Coord]struct{}{}

		for _, t := range trailheads {
			unique[t] = struct{}{}
		}

		res += len(unique)
	}
	return res
}

func countTrailheadRatings(m [][]int, starts []Coord) int {
	trailheads := climb(m, starts)
	return len(trailheads)
}

func main() {
	m, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	starts := findStarts(m)
	trailheads := countTrailheads(m, starts)
	ratings := countTrailheadRatings(m, starts)

	fmt.Println(trailheads, ratings)
}
