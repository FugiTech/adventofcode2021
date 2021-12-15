package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
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
	lines := strings.Split(input, "\n")
	score := 0

LineLoop:
	for _, l := range lines {
		stack := make([]byte, 0, 100)
		for _, c := range []byte(l) {
			// fmt.Printf("%s %c\n", stack, c)
			switch c {
			case '(', '[', '{', '<':
				stack = append(stack, c)
			case ')':
				if stack[len(stack)-1] != '(' {
					score += 3
					continue LineLoop
				}
				stack = stack[:len(stack)-1]
			case ']':
				if stack[len(stack)-1] != '[' {
					score += 57
					continue LineLoop
				}
				stack = stack[:len(stack)-1]
			case '}':
				if stack[len(stack)-1] != '{' {
					score += 1197
					continue LineLoop
				}
				stack = stack[:len(stack)-1]
			case '>':
				if stack[len(stack)-1] != '<' {
					score += 25137
					continue LineLoop
				}
				stack = stack[:len(stack)-1]
			default:
				return 0, fmt.Errorf("Invalid character: %c", c)
			}
		}
	}

	return score, nil
}

func AnswerV1(input string) (int, error) {
	lines := strings.Split(input, "\n")
	scores := make([]int, 0, 1000)

LineLoop:
	for _, l := range lines {
		stack := make([]byte, 0, 100)
		for _, c := range []byte(l) {
			// fmt.Printf("%s %c\n", stack, c)
			switch c {
			case '(', '[', '{', '<':
				stack = append(stack, c)
			case ')':
				if stack[len(stack)-1] != '(' {
					continue LineLoop
				}
				stack = stack[:len(stack)-1]
			case ']':
				if stack[len(stack)-1] != '[' {
					continue LineLoop
				}
				stack = stack[:len(stack)-1]
			case '}':
				if stack[len(stack)-1] != '{' {
					continue LineLoop
				}
				stack = stack[:len(stack)-1]
			case '>':
				if stack[len(stack)-1] != '<' {
					continue LineLoop
				}
				stack = stack[:len(stack)-1]
			default:
				return 0, fmt.Errorf("Invalid character: %c", c)
			}
		}

		score := 0
		for i := len(stack) - 1; i >= 0; i-- {
			switch stack[i] {
			case '(':
				score = 5*score + 1
			case '[':
				score = 5*score + 2
			case '{':
				score = 5*score + 3
			case '<':
				score = 5*score + 4
			}
		}
		scores = append(scores, score)
	}

	sort.Ints(scores)

	return scores[len(scores)/2], nil
}
