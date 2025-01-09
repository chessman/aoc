//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func parse() ([]string, []string, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()

	patterns := strings.Split(scanner.Text(), ", ")

	scanner.Scan()

	var designs []string
	for {
		if !scanner.Scan() {
			return patterns, designs, scanner.Err()
		}

		designs = append(designs, scanner.Text())
	}
}

func countWays(patterns []string, design string, cache map[string]int) int {
	if count, ok := cache[design]; ok {
		return count
	}

	var count int
	for _, pat := range patterns {
		after, found := strings.CutPrefix(design, pat)
		if !found {
			continue
		}

		if len(after) == 0 {
			count++
			continue
		}

		count += countWays(patterns, after, cache)
	}

	cache[design] = count
	return count
}

func countPossibleDesigns(patterns []string, designs []string) int {
	count := 0
	cache := make(map[string]int)
	for _, design := range designs {
		if countWays(patterns, design, cache) > 0 {
			count++
		}
	}
	return count
}

func countPossibleWays(patterns []string, designs []string) int {
	count := 0
	cache := make(map[string]int)
	for _, design := range designs {
		count += countWays(patterns, design, cache)
	}
	return count
}

func main() {
	patterns, designs, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(countPossibleDesigns(patterns, designs))
	fmt.Println(countPossibleWays(patterns, designs))
}
