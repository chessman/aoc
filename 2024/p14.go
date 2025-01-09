//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	MaxX = 101
	MaxY = 103
)

type P struct {
	x int
	y int
}

type Robot struct {
	p P
	v P
}

func parse() ([]Robot, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	var robots []Robot

	for {
		var r Robot
		if !scanner.Scan() {
			return robots, scanner.Err()
		}
		
		_, err = fmt.Sscanf(scanner.Text(), "p=%d,%d v=%d,%d" , &r.p.x, &r.p.y, &r.v.x, &r.v.y)
		if err != nil {
			return nil, err
		}
		
		robots = append(robots, r)
	}
}

func simulate(steps int, robots []Robot) []Robot {
	var newRobots []Robot
	for _, r := range robots {
		r.p.x = (r.p.x + r.v.x * steps) % MaxX
		r.p.y = (r.p.y + r.v.y * steps) % MaxY
		if r.p.x < 0 { r.p.x += MaxX }
		if r.p.y < 0 { r.p.y += MaxY }
		newRobots = append(newRobots, r)
	}
	return newRobots
}

func simulateUntilDistinct(robots []Robot) int {
	for c := 0; c < 2000000; c++ {
		rs := simulate(c, robots)
		if allDistinct(rs) {
			return c
		}
	}
	return 0
}

func allDistinct(robots []Robot) bool {
	m := make([][]bool, MaxY, MaxY)
	for i := 0; i < MaxY; i++ {
		m[i] = make([]bool, MaxX, MaxX)
	}

	for _, r := range robots {
		if m[r.p.y][r.p.x] {
			return false
		}
		m[r.p.y][r.p.x] = true
	}
	return true
}

func safetyFactor(robots []Robot) int {
	var topLeft, topRight, bottomLeft, bottomRight int
	for _, r := range robots {
		if r.p.x < MaxX / 2 && r.p.y < MaxY / 2 {
			topLeft++
		}
		if r.p.x < MaxX / 2 && r.p.y > MaxY / 2 {
			bottomLeft++
		}
		if r.p.x > MaxX / 2 && r.p.y < MaxY / 2 {
			topRight++
		}
		if r.p.x > MaxX / 2 && r.p.y > MaxY / 2 {
			bottomRight++
		}
	}
	return topLeft * bottomLeft * topRight * bottomRight
}

func printRobots(robots []Robot) {
	m := make([][]bool, MaxY, MaxY)
	for i := 0; i < MaxY; i++ {
		m[i] = make([]bool, MaxX, MaxX)
	}

	for _, r := range robots {
		m[r.p.y][r.p.x] = true
	}
	for y := 0; y < MaxY; y++ {
		for x := 0; x < MaxX; x++ {
			if m[y][x] {
				fmt.Print("*")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func main() {
	robots, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	printRobots(simulate(6620, robots))
	fmt.Println(safetyFactor(simulate(100, robots)))
	fmt.Println(simulateUntilDistinct(robots))
}
