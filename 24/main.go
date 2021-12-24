package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
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

var Answer = AnswerPart1

func validMonad(monad int, lines []string) bool {
	input := fmt.Sprintf("%014d", monad)
	inputIdx := 0
	registers := [4]int{}

	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		p := strings.Split(l, " ")
		r := p[1][0] - 'w'
		v := 0
		if len(p) > 2 {
			r2 := int(p[2][0]) - 'w'
			if 0 <= r2 && r2 < 4 {
				v = registers[r2]
			} else {
				v, _ = strconv.Atoi(p[2])
			}
		}
		switch p[0] {
		case "inp":
			registers[r] = int(input[inputIdx] - '0')
			inputIdx++
		case "add":
			registers[r] += v
		case "mul":
			registers[r] *= v
		case "div":
			registers[r] /= v
		case "mod":
			registers[r] = registers[r] % v
		case "eql":
			if registers[r] == v {
				registers[r] = 1
			} else {
				registers[r] = 0
			}
		}
	}

	fmt.Println(input, registers)
	return registers[3] == 0
}

type state struct {
	input     [14]byte
	registers [4]int
}

var maximums = [14]int{26, 676, 17576, 676, 26, 676, 17576, 676, 17576, 456976, 17576, 676, 26, 1}

func runBlock(s state, lines []string, inputIdx int) (state, bool) {
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}

		p := strings.Split(l, " ")
		r := p[1][0] - 'w'
		v := 0
		if len(p) > 2 {
			r2 := int(p[2][0]) - 'w'
			if 0 <= r2 && r2 < 4 {
				v = s.registers[r2]
			} else {
				v, _ = strconv.Atoi(p[2])
			}
		}
		switch p[0] {
		case "inp":
			s.registers[r] = int(s.input[inputIdx])
		case "add":
			s.registers[r] += v
		case "mul":
			s.registers[r] *= v
		case "div":
			s.registers[r] /= v
		case "mod":
			s.registers[r] = s.registers[r] % v
		case "eql":
			if s.registers[r] == v {
				s.registers[r] = 1
			} else {
				s.registers[r] = 0
			}
		}
	}

	return s, s.registers[3] < maximums[inputIdx]
}

func AnswerPart1(input string) (int, error) {
	blocks := strings.Split(input, "\n\n")
	states := []state{
		{},
	}

	for bIdx, block := range blocks {
		lines := strings.Split(block, "\n")
		nextStates := []state{}
		for _, s := range states {
			for i := 1; i < 10; i++ {
				in := s
				in.input[bIdx] = byte(i)
				out, valid := runBlock(in, lines, bIdx)
				if valid {
					nextStates = append(nextStates, out)
				}
			}
		}
		states = nextStates
		log.Printf("% 3d % 7d % 9d", bIdx+1, maximums[bIdx], len(states))
	}

	sort.Slice(states, func(i, j int) bool {
		return states[i].registers[3] < states[j].registers[3]
	})
	for _, s := range states {
		fmt.Println(s)
	}

	return 0, nil
}

func AnswerV1(input string) (int, error) {

	return 0, nil
}
