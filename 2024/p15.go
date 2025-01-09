//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"os"
)

type P struct {
	x int
	y int
}

func parse() ([][]rune, []rune, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanMap := true

	var m [][]rune
	var moves []rune
	for {
		if !scanner.Scan() {
			return m, moves, scanner.Err()
		}

		line := scanner.Text()
		if line == "" {
			scanMap = false
			continue
		}

		var row []rune
		for _, c := range line {
			row = append(row, c)
		}
		if scanMap {
			m = append(m, row)
		} else {
			moves = append(moves, row...)
		}
	}
}

func robotP(m [][]rune) P {
	for y := range m {
		for x := range m[y] {
			if m[y][x] == '@' {
				return P{x, y}
			}
		}
	}
	return P{0, 0}
}

var moveMap map[rune]func(p P) P = map[rune]func(p P) P{
	'^': func(p P) P { return P{p.x, p.y - 1} },
	'>': func(p P) P { return P{p.x + 1, p.y} },
	'<': func(p P) P { return P{p.x - 1, p.y} },
	'v': func(p P) P { return P{p.x, p.y + 1} },
}

func plan(m [][]rune, p P, move rune, moves map[P]rune) *map[P]rune {
	el := m[p.y][p.x]
	switch el {
	case '#':
		return nil
	case '.':
		return &moves
	case '[', ']':
		if (move == '^' || move == 'v') {
			var start P
			if el == ']' {
				start = P{p.x - 1, p.y}
			} else {
				start = P{p.x + 1, p.y}
			}
			_, ok := moves[start]
			// if ok, we are already in the branch
			if !ok {
				moves[start] = '.'
				newMoves := plan(m, start, move, moves)
				if newMoves == nil {
					return nil
				}
			}
		}
	}
	newP := moveMap[move](p)
	moves[newP] = el

	return plan(m, newP, move, moves)
}

func makeMove(m [][]rune, init P, move rune) P {
	planned := plan(m, init, move, map[P]rune{init: '.'})
	if planned == nil {
		return init
	}

	for p, el := range *planned {
		m[p.y][p.x] = el
	}

	return moveMap[move](init)
}

func makeMoves(m [][]rune, init P, moves []rune) {
	p := init
	for _, move := range moves {
		p = makeMove(m, p, move)
	}
}

func sumCoordinates(m [][]rune) int {
	total := 0
	for y := range m {
		for x := range m[y] {
			if m[y][x] == 'O' || m[y][x] == '[' {
				total += y*100 + x
			}
		}
	}
	return total
}

func widenMap(m [][]rune) [][]rune {
	var newM [][]rune
	for y := range m {
		var row []rune
		for x := range m[y] {
			switch m[y][x] {
			case '#':
				row = append(row, '#', '#')
			case 'O':
				row = append(row, '[', ']')
			case '.':
				row = append(row, '.', '.')
			case '@':
				row = append(row, '@', '.')
			}
		}
		newM = append(newM, row)
	}
	return newM
}

func main() {
	m, moves, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	w := widenMap(m)

	makeMoves(m, robotP(m), moves)
	makeMoves(w, robotP(w), moves)
	fmt.Println(sumCoordinates(m))
	fmt.Println(sumCoordinates(w))
}
