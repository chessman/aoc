//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coord struct {
	x int
	y int
}

type Machine struct {
	buttonA Coord
	buttonB Coord
	prize   Coord
}

func parse() ([]Machine, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	var machines []Machine

	for {
		var m Machine
		if !scanner.Scan() {
			return machines, scanner.Err()
		}
		
		_, err = fmt.Sscanf(scanner.Text(), "Button A: X%d, Y%d" , &m.buttonA.x, &m.buttonA.y)
		if err != nil {
			return nil, err
		}
		
		if !scanner.Scan() {
			return nil, scanner.Err()
		}

		_, err = fmt.Sscanf(scanner.Text(), "Button B: X%d, Y%d" , &m.buttonB.x, &m.buttonB.y)
		if err != nil {
			return nil, err
		}

		if !scanner.Scan() {
			return nil, scanner.Err()
		}

		_, err = fmt.Sscanf(scanner.Text(), "Prize: X=%d, Y=%d" , &m.prize.x, &m.prize.y)
		if err != nil {
			return nil, err
		}

		machines = append(machines, m)

                if !scanner.Scan() {
			return machines, nil
		}
	}
}

func minTokens(m Machine) int {
	d := m.buttonA.x * m.buttonB.y - m.buttonA.y * m.buttonB.x

	if d == 0 {
		return 0
	}

	n1 := (m.prize.y * m.buttonB.x - m.prize.x * m.buttonB.y) / (-d)
	n2 := (m.prize.y * m.buttonA.x - m.prize.x * m.buttonA.y) / d

	if n1 == 0 && n2 == 0 {
		return min(m.prize.x / m.buttonA.x * 3, m.prize.x / m.buttonB.x)
	}

	if n1 * m.buttonA.x + n2 * m.buttonB.x != m.prize.x || n1 * m.buttonA.y + n2 * m.buttonB.y != m.prize.y {
		return 0
	}

	return n1 * 3 + n2
}

func minTokensAll(machines []Machine) int {
	total := 0
	for _, m := range machines {
		total += minTokens(m)
	}
	return total
}

func main() {
	machines, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(minTokensAll(machines))

	var newMachines []Machine
	for _, m := range machines {
		m.prize.x += 10000000000000
		m.prize.y += 10000000000000
		newMachines = append(newMachines, m)
	}

	fmt.Println(minTokensAll(newMachines))
}
