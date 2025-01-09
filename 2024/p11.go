//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parse() ([]int, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	scanner.Scan()
	line := scanner.Text()
	strs := strings.Split(line, " ")
	var res []int
	for _, s := range strs {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}

func digits(s int) int {
	digs := 1
	for c := s; c > 9; c /= 10 {
		digs++
	}
	return digs
}

type val struct {
	s int
	count int
}

var cache = make(map[val]int)
func blink(s int, count int) int {
	if v, ok := cache[val{s, count}]; ok {
		return v
	} else {
		var res int
		if count == 0 {
			res = 1
		} else {
			if s == 0 {
				res = blink(1, count - 1)
			} else {
				digs := digits(s)
				if (digs % 2 == 0) {
					div := 1
					for range digs / 2 {
						div *= 10
					}
					res = blink(s / div, count - 1) + blink(s % div, count - 1)
				} else {
					res = blink(s * 2024, count - 1)
				}
			}
		}
		cache[val{s, count}] = res
		return res
	}
}

func main() {
	stones, err := parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var sum int
	for _, s := range stones {
		sum += blink(s, 75)
	}

	fmt.Printf("number of stones: %v\n", sum)
}
