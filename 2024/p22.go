//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func parse() ([]int, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var secrets []int
	for {
		if !scanner.Scan() {
			return secrets, scanner.Err()
		}

		code, _ := strconv.Atoi(scanner.Text())
		secrets = append(secrets, code)
	}
}

func next(secret int) int {
	step1 := ((secret * 64) ^ secret) % 16777216
	step2 := ((step1 / 32) ^ step1) % 16777216
	step3 := ((step2 * 2048) ^ step2) % 16777216

	return step3
}

func nsecret(secret int, n int) int {
	for range n {
		secret = next(secret)
	}

	return secret
}
func sumSecrets(secrets []int, n int) int {
	total := 0
	for _, s := range secrets {
		total += nsecret(s, n)
	}
	return total
}

func nprices(secret int, n int) []int {
	var prices []int
	for range n {
		prices = append(prices, secret%10)
		secret = next(secret)
	}
	return prices
}

type Seq struct {
	c1 int
	c2 int
	c3 int
	c4 int
}

func makeSeqmap(secret int, n int) map[Seq]int {
	s := make(map[Seq]int)

	prices := nprices(secret, n)
	var changes []int

	for i := range len(prices) - 1 {
		changes = append(changes, prices[i+1]-prices[i])
	}

	for i := range len(prices) - 4 {
		seq := Seq{changes[i], changes[i+1], changes[i+2], changes[i+3]}
		if _, ok := s[seq]; !ok {
			s[seq] = prices[i+4]
		}
	}

	return s
}

func totalBananas(secrets []int, n int) int {
	var seqmaps []map[Seq]int
	for _, s := range secrets {
		seqmaps = append(seqmaps, makeSeqmap(s, n))
	}

	uniqueSeqs := make(map[Seq]struct{})
	for _, seqmap := range seqmaps {
		for seq := range seqmap {
			uniqueSeqs[seq] = struct{}{}
		}
	}

	maxBananas := 0
	for seq := range uniqueSeqs {
		bananas := 0
		for _, seqmap := range seqmaps {
			bananas += seqmap[seq]
		}
		if bananas > maxBananas {
			maxBananas = bananas
		}
	}
	return maxBananas
}

func main() {
	secrets, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(sumSecrets(secrets, 2000))
	fmt.Println(totalBananas(secrets, 2000))
}
