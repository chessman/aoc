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

var numKeypad map[rune]P = map[rune]P{
	'7': {0, 0},
	'8': {1, 0},
	'9': {2, 0},
	'4': {0, 1},
	'5': {1, 1},
	'6': {2, 1},
	'1': {0, 2},
	'2': {1, 2},
	'3': {2, 2},
	'0': {1, 3},
	'A': {2, 3},
}

var dirKeypad map[rune]P = map[rune]P{
	'^': {1, 0},
	'A': {2, 0},
	'<': {0, 1},
	'v': {1, 1},
	'>': {2, 1},
}

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func calcMoves(start, finish, forbidden P) [][]rune {
	xdiff := finish.x - start.x
	ydiff := finish.y - start.y

	if start == finish {
		return [][]rune{{'A'}}
	}

	var res [][]rune
	if xdiff < 0 {
		newStart := P{start.x - 1, start.y}
		if newStart != forbidden {
			subres := calcMoves(newStart, finish, forbidden)
			for _, sr := range subres {
				res = append(res, append([]rune{'<'}, sr...))
			}
		}
	}

	if xdiff > 0 {
		newStart := P{start.x + 1, start.y}
		if newStart != forbidden {
			subres := calcMoves(newStart, finish, forbidden)
			for _, sr := range subres {
				res = append(res, append([]rune{'>'}, sr...))
			}
		}
	}

	if ydiff < 0 {
		newStart := P{start.x, start.y - 1}
		if newStart != forbidden {
			subres := calcMoves(newStart, finish, forbidden)
			for _, sr := range subres {
				res = append(res, append([]rune{'^'}, sr...))
			}
		}
	}

	if ydiff > 0 {
		newStart := P{start.x, start.y + 1}
		if newStart != forbidden {
			subres := calcMoves(newStart, finish, forbidden)
			for _, sr := range subres {
				res = append(res, append([]rune{'v'}, sr...))
			}
		}
	}

	return res
}

func numMoves(num1 rune, num2 rune) [][]rune {
	start := numKeypad[num1]
	finish := numKeypad[num2]
	forbidden := P{0, 3}

	return calcMoves(start, finish, forbidden)
}

func dirMoves(dir1 rune, dir2 rune) [][]rune {
	start := dirKeypad[dir1]
	finish := dirKeypad[dir2]
	forbidden := P{0, 0}

	return calcMoves(start, finish, forbidden)
}

type C struct {
	robot  int
	start  rune
	finish rune
}

func getRobotOneCost(robot, totalRobots int, start, finish rune, cache map[C]int) int {
	if res, ok := cache[C{robot, start, finish}]; ok {
		return res
	}
	genMoves := dirMoves
	if robot == 0 {
		genMoves = numMoves
	}
	res := math.MaxInt
	for _, m := range genMoves(start, finish) {
		if robot == totalRobots-1 {
			if res > len(m) {
				res = len(m)
			}
		} else {
			cost := getRobotCost(robot+1, totalRobots, m, cache)
			if cost < res {
				res = cost
			}
		}
	}

	cache[C{robot, start, finish}] = res
	return res
}

func getRobotCost(robot, totalRobots int, moves []rune, cache map[C]int) int {
	total := 0
	moves = append([]rune{'A'}, moves...)
	for i := range moves[:len(moves)-1] {
		j := i + 1
		total += getRobotOneCost(robot, totalRobots, moves[i], moves[j], cache)
	}

	return total
}

func parse() ([][]rune, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var codes [][]rune
	for {
		if !scanner.Scan() {
			return codes, scanner.Err()
		}

		codes = append(codes, []rune(scanner.Text()))
	}
}

func codeScore(totalRobots int, code []rune, cache map[C]int) int {
	seqLen := getRobotCost(0, totalRobots, code, cache)
	var num int
	fmt.Sscanf(string(code), "%d", &num)
	return seqLen * num
}

func codeScores(totalRobots int, codes [][]rune, cache map[C]int) int {
	total := 0
	for _, code := range codes {
		total += codeScore(totalRobots, code, cache)
	}
	return total
}

func main() {
	codes, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(codeScores(3, codes, make(map[C]int)))
	fmt.Println(codeScores(26, codes, make(map[C]int)))
}
