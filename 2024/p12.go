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

func parse() ([][]rune, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	var m [][]rune

	for scanner.Scan() {
		line := scanner.Text()
		var row []rune
		for _, char := range line {
			row = append(row, char)
		}
		m = append(m, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

type Side int

const (
	North Side = iota
	East
	South
	West
)

type Res struct {
	area      int
	perimeter int
	sides     int
}

func (r Res) Plus(r2 Res) Res {
	return Res{r.area + r2.area, r.perimeter + r2.perimeter, r.sides + r2.sides}
}

func (r Res) Price() int {
	return r.area * r.perimeter
}

func (r Res) PriceDiscounted() int {
	return r.area * r.sides
}

func areaAndPerimeter(m [][]rune) []Res {
	visited := make(map[Coord][]Side)
	var results []Res
	for y := range m {
		for x := range m[y] {
			c := Coord{x, y}
			if _, ok := visited[c]; !ok {
				res := areaAndPerimeterOne(m, c, visited)
				if res.sides%2 == 1 {
					//fmt.Printf("%q %v %v\n", m[y][x], res, c)
					res.sides--
				}
				results = append(results, res)
			}
		}
	}
	return results
}

func checkNeighbor(m [][]rune, c Coord, n Coord, side Side, visited map[Coord][]Side) bool {
	if n.y >= 0 && n.y < len(m) && n.x >= 0 && n.x < len(m[0]) && m[c.y][c.x] == m[n.y][n.x] {
		sides := visited[n]
		for _, s := range sides {
			if s == side {
				return true
			}
		}

	}
	return false
}

func checkNeighbors(m [][]rune, c Coord, side Side, visited map[Coord][]Side) bool {
	if side == North || side == South {
		return checkNeighbor(m, c, Coord{c.x - 1, c.y}, side, visited) ||
			checkNeighbor(m, c, Coord{c.x + 1, c.y}, side, visited)
	} else {
		return checkNeighbor(m, c, Coord{c.x, c.y - 1}, side, visited) ||
			checkNeighbor(m, c, Coord{c.x, c.y + 1}, side, visited)
	}
}

func areaAndPerimeterOne(m [][]rune, c Coord, visited map[Coord][]Side) Res {
	v := m[c.y][c.x]
	if _, ok := visited[c]; ok {
		return Res{0, 0, 0}
	}
	var sides []Side
	res := Res{1, 0, 0}
	if !(c.y-1 >= 0 && m[c.y-1][c.x] == v) {
		sides = append(sides, North)
	}

	if !(c.y+1 < len(m) && m[c.y+1][c.x] == v) {
		sides = append(sides, South)
	}

	if !(c.x-1 >= 0 && m[c.y][c.x-1] == v) {
		sides = append(sides, West)
	}

	if !(c.x+1 < len(m[0]) && m[c.y][c.x+1] == v) {
		sides = append(sides, East)
	}

	for _, s := range sides {
		if !checkNeighbors(m, c, s, visited) {
			res.sides++
		}
	}

	visited[c] = sides
	res.perimeter = len(sides)

	if c.y-1 >= 0 && m[c.y-1][c.x] == v {
		res = res.Plus(areaAndPerimeterOne(m, Coord{c.x, c.y - 1}, visited))
	}
	if c.y+1 < len(m) && m[c.y+1][c.x] == v {
		res = res.Plus(areaAndPerimeterOne(m, Coord{c.x, c.y + 1}, visited))
	}
	if c.x-1 >= 0 && m[c.y][c.x-1] == v {
		res = res.Plus(areaAndPerimeterOne(m, Coord{c.x - 1, c.y}, visited))
	}
	if c.x+1 < len(m[0]) && m[c.y][c.x+1] == v {
		res = res.Plus(areaAndPerimeterOne(m, Coord{c.x + 1, c.y}, visited))
	}
	return res
}

func totalPrice(results []Res) int {
	total := 0
	for _, res := range results {
		total += res.Price()
	}
	return total
}

func totalDiscountedPrice(results []Res) int {
	total := 0
	for _, res := range results {
		total += res.PriceDiscounted()
	}
	return total
}

func main() {
	m, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	res := areaAndPerimeter(m)

	//fmt.Println(areaAndPerimeterOne(m, Coord{2, 0}, make(map[Coord][]Side)))

	// fmt.Println(totalPrice(res))
	fmt.Println(totalDiscountedPrice(res))
}
