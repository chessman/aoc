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

type Cmd struct {
	op1 string
	op2 string
	cmd string
}

type System struct {
	input map[string]int
	wires map[string]Cmd
}

func (s System) zouts() []string {
	var zouts []string
	for k := range s.wires {
		if k[0] == 'z' {
			zouts = append(zouts, k)
		}
	}

	slices.Sort(zouts)
	return zouts
}

func parse() (System, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return System{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	input := make(map[string]int)
	wires := make(map[string]Cmd)
	parseConsts := true

	for {
		if !scanner.Scan() {
			return System{input, wires}, scanner.Err()
		}

		s := scanner.Text()
		if parseConsts {
			if len(s) == 0 {
				parseConsts = false
				continue
			}

			i, _ := strconv.Atoi(s[5:6])
			input[s[0:3]] = i
		} else {
			var op1, cmd, op2, res string
			fmt.Sscanf(s, "%s %s %s -> %s", &op1, &cmd, &op2, &res)
			wires[res] = Cmd{op1: op1, op2: op2, cmd: cmd}
		}
	}
}

func resolve(output string, system System, cache map[string]int) int {
	if v, ok := system.input[output]; ok {
		return v
	}
	gate := system.wires[output]
	v1 := resolve(gate.op1, system, cache)
	v2 := resolve(gate.op2, system, cache)
	var res int
	switch gate.cmd {
	case "AND":
		res = v1 & v2
	case "OR":
		res = v1 | v2
	case "XOR":
		res = v1 ^ v2
	}

	cache[output] = res
	return res
}

func run(system System) int {
	zouts := system.zouts()
	cache := make(map[string]int)

	var res int
	for i, o := range zouts {
		if resolve(o, system, cache) == 1 {
			res |= 1 << i
		}
	}

	return res
}

func findDiscrepancies(system System) string {
	// andXor(n) = orAnd(n-1) AND xor(n-1)
	// orAnd(n) = andXor(n) OR and(n-1)
	// z(n) = orAnd(n) XOR xor(n)

	xyXors := make([]string, 46)
	xyAnds := make([]string, 46)

	for output, cmd := range system.wires {
		if (cmd.op1[0] == 'x' && cmd.op2[0] == 'y') || (cmd.op1[0] == 'y' && cmd.op2[0] == 'x') {
			i, _ := strconv.Atoi(cmd.op1[1:3])
			switch cmd.cmd {
			case "XOR":
				xyXors[i] = output
			case "AND":
				xyAnds[i] = output
			default:
				panic("not expected")
			}
		}
	}

	var wrong []string

	orAnds := make([]string, 46)
	andXors := make([]string, 46)

	for output, cmd := range system.wires {
		if (cmd.cmd == "OR") {
			i := slices.Index(xyAnds, cmd.op1)
			if i == -1 {
				i = slices.Index(xyAnds, cmd.op2)
			}
			if i != -1 {
				orAnds[i+1] = output
			}
		} else if (cmd.cmd == "AND" && cmd.op1[0] != 'x' && cmd.op1[0] != 'y') {
			i := slices.Index(xyXors, cmd.op1)
			if i == -1 {
				i = slices.Index(xyXors, cmd.op2)
			}
			if i != -1 {
				andXors[i+1] = output
			}
		}
	}

	findTrueZ := func (i int) string {
		xor := xyXors[i]
		for output, cmd := range system.wires {
			if cmd.cmd == "XOR" && (cmd.op1 == xor || cmd.op2 == xor) {
				return output
			}
		}
		return ""
	}

	zouts := system.zouts()
	for i, zout := range zouts {
		cmd := system.wires[zout]
		if cmd.cmd != "XOR" {
			if xyXors[i] != "" {
				wrong = append(wrong, zout)
				wrong = append(wrong, findTrueZ(i))
			}
			continue
		}
		xor := xyXors[i]
		orAnd := orAnds[i]
		if xor == "" || orAnd == "" {
			continue
		}
		if cmd.op1 == orAnd && cmd.op2 != xor {
			wrong = append(wrong, cmd.op2)
			wrong = append(wrong, xor)
		}
		if cmd.op2 == orAnd && cmd.op1 != xor {
			wrong = append(wrong, cmd.op1)
			wrong = append(wrong, xor)
		}
	}

	slices.Sort(wrong)
	return strings.Join(wrong, ",")
}

func main() {
	system, err := parse()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(run(system))
	fmt.Println(findDiscrepancies(system))
}
