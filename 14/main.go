package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	answer, err := Answer(strings.TrimSpace(string(input)))
	fmt.Printf("answer: %v\n", answer)
	return err
}

var Answer = AnswerV1

func AnswerPart1(input string) (int, error) {
	template := []uint8{}
	rules := map[uint16]uint8{}

	printTemplate := func() {
		for _, c := range template {
			fmt.Printf("%c", c+'A')
		}
		fmt.Println()
	}
	_ = printTemplate

	lines := strings.Split(input, "\n")
	for _, c := range []byte(lines[0]) {
		template = append(template, uint8(c-'A'))
	}
	for _, l := range lines[2:] {
		p := strings.Split(l, " -> ")
		k := (uint16(p[0][0]-'A') << 8) | uint16(p[0][1]-'A')
		v := uint8(p[1][0] - 'A')
		rules[k] = v
	}

	for step := 0; step < 10; step++ {
		n := make([]uint8, 0, len(template)*2)
		for i := 0; i < len(template)-1; i++ {
			n = append(n, template[i])
			k := (uint16(template[i]) << 8) | uint16(template[i+1])
			n = append(n, rules[k])
		}
		n = append(n, template[len(template)-1])
		template = n
	}

	var counts [26]int
	for _, c := range template {
		counts[c]++
	}

	max, min := 0, math.MaxInt
	for _, c := range counts[:] {
		if c > max {
			max = c
		}
		if c < min && c > 0 {
			min = c
		}
	}

	return max - min, nil
}

func AnswerV1(input string) (int, error) {
	rules := map[string]string{}
	counts := map[string]int{}

	lines := strings.Split(input, "\n")
	for i, l := 0, len(lines[0])-1; i < l; i++ {
		counts[lines[0][i:i+2]]++
	}
	for _, l := range lines[2:] {
		p := strings.Split(l, " -> ")
		rules[p[0]] = p[1]
	}

	for step := 0; step < 40; step++ {
		n := map[string]int{}
		for k, v := range counts {
			i := rules[k]
			a, b := k[:1]+i, i+k[1:]
			n[a] += v
			n[b] += v
		}
		counts = n
	}

	icounts := map[string]int{}
	for k, v := range counts {
		icounts[k[:1]] += v
		icounts[k[1:]] += v
	}
	for k, v := range icounts {
		icounts[k] = (v + 1) / 2
	}

	fmt.Println(icounts)

	max, min := 0, math.MaxInt
	for _, v := range icounts {
		if v > max {
			max = v
		}
		if v < min && v > 0 {
			min = v
		}
	}

	return max - min, nil
}
