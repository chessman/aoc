//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"os"
)

type Lock []int
type Key []int

func parse() ([]Lock, []Key, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var locks []Lock
	var keys []Key
	for {
		var cur [][]rune
		for range 7 {
			scanner.Scan()
			s := scanner.Text()
			cur = append(cur, []rune(s))
		}

		keyOrLock := make([]int, len(cur[0]))
		for x := range cur[0] {
			keyOrLock[x] = -1
			for y := range cur {
				if cur[y][x] == '#' {
					keyOrLock[x]++
				}
			}
		}

		if cur[0][0] == '#' {
			locks = append(locks, Lock(keyOrLock))
		} else {
			keys = append(keys, Key(keyOrLock))
		}

		if !scanner.Scan() {
			return locks, keys, scanner.Err()
		}
	}
}

func fit(lock Lock, key Key) bool {
	for i := range lock {
		if lock[i] + key[i] > 5 {
			return false
		}
	}
	return true
}

func pairCount(locks []Lock, keys[]Key) int {
	total := 0
	for _, l := range locks {
		for _, k := range keys {
			if fit(l, k) {
				total++
			}
		}
	}
	return total
}

func main() {
	locks, keys, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	// 11291

	fmt.Println(pairCount(locks, keys))
}
