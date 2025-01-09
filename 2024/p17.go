//go:build ignore
package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Proc struct {
	a    int
	b    int
	c    int
	prog []int

	p      int
	stdout []int
}

func (p *Proc) comboOperand(operand int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return p.a
	case 5:
		return p.b
	case 6:
		return p.c
	}

	panic(fmt.Errorf("wrong combo operand value: %d", operand))
}

func (p *Proc) next() {
	p.p += 2
}

func (p *Proc) adv(operand int) {
	v := p.comboOperand(operand)
	if v > 0 {
		p.a = p.a / (2 << (v - 1))
	}
	p.next()
}

func (p *Proc) bxl(operand int) {
	p.b = p.b ^ operand
	p.next()
}

func (p *Proc) bst(operand int) {
	p.b = p.comboOperand(operand) % 8
	p.next()
}

func (p *Proc) jnz(operand int) {
	if p.a == 0 {
		p.next()
	} else {
		p.p = operand
	}
}

func (p *Proc) bxc(_ int) {
	p.b = p.b ^ p.c
	p.next()
}

func (p *Proc) out(operand int) {
	p.stdout = append(p.stdout, p.comboOperand(operand)%8)
	p.next()
}

func (p *Proc) bdv(operand int) {
	v := p.comboOperand(operand)
	if v > 0 {
		p.b = p.a / (2 << (v - 1))
	} else {
		p.b = p.a
	}
	p.next()
}

func (p *Proc) cdv(operand int) {
	v := p.comboOperand(operand)
	if v > 0 {
		p.c = p.a / (2 << (v - 1))
	} else {
		p.c = p.a
	}
	p.next()
}

func (p *Proc) cmd(op, operand int) {
	switch op {
	case 0:
		p.adv(operand)
	case 1:
		p.bxl(operand)
	case 2:
		p.bst(operand)
	case 3:
		p.jnz(operand)
	case 4:
		p.bxc(operand)
	case 5:
		p.out(operand)
	case 6:
		p.bdv(operand)
	case 7:
		p.cdv(operand)
	default:
		panic(fmt.Errorf("cmd: wrong op: %d", op))
	}
}

func (p *Proc) runOne() {
	op := p.prog[p.p]
	operand := p.prog[p.p+1]
	p.cmd(op, operand)
}

func (p *Proc) run() string {
	for {
		if p.p >= len(p.prog) {
			var strs []string
			for _, i := range p.stdout {
				strs = append(strs, strconv.Itoa(i))
			}
			return strings.Join(strs, ",")
		}

		p.runOne()
	}
}

func (p *Proc) quine() int {
	ans := 0
	c := len(p.prog) - 1
	for {
		p2 := *p
		p2.a = ans
		p2.stdout = nil
		for {
			if p2.p >= len(p2.prog) {
				break
			}
			p2.runOne()
			if len(p2.stdout) == len(p.prog)-c {
				if slices.Compare(p.prog[c:], p2.stdout) == 0 {
					c--
					if c == -1 {
						return ans
					}
					ans = ans << 3
				} else {
					break
				}
			}
		}
		ans++
	}
}

func parse() (*Proc, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var proc Proc

	scanner.Scan()
	_, err = fmt.Sscanf(scanner.Text(), "Register A: %d", &proc.a)
	if err != nil {
		return nil, err
	}

	scanner.Scan()
	_, err = fmt.Sscanf(scanner.Text(), "Register B: %d", &proc.b)
	if err != nil {
		return nil, err
	}

	scanner.Scan()
	_, err = fmt.Sscanf(scanner.Text(), "Register C: %d", &proc.c)
	if err != nil {
		return nil, err
	}

	scanner.Scan()
	scanner.Scan()
	var program string
	_, err = fmt.Sscanf(scanner.Text(), "Program: %s", &program)
	if err != nil {
		return nil, err
	}

	for _, s := range strings.Split(program, ",") {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		proc.prog = append(proc.prog, i)
	}

	return &proc, nil
}
func main() {
	proc, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(proc.quine())

	// proc.a=33200148537459
	// fmt.Println(proc.run())

	// a := proc.quine()
	// fmt.Println(a)
}
