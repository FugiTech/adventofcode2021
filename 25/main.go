package main

import (
	"fmt"
	"io/ioutil"
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

var Answer = AnswerPart1

func AnswerPart1(input string) (int, error) {
	lines := strings.Split(input, "\n")
	board := make([][]uint8, len(lines))

	for j, l := range lines {
		row := make([]uint8, len(l))
		for i, c := range []byte(l) {
			if c == '>' {
				row[i] = 1
			} else if c == 'v' {
				row[i] = 2
			}
		}
		board[j] = row
	}

	steps := 0
	for moved := 1; moved > 0; steps++ {
		moved = 0

		var toMove [][4]int
		for y, row := range board {
			for x, v := range row {
				tx := (x + 1) % len(row)
				if v == 1 && board[y][tx] == 0 {
					toMove = append(toMove, [4]int{x, y, tx, y})
				}
			}
		}
		for _, v := range toMove {
			x, y, tx, ty := v[0], v[1], v[2], v[3]
			board[ty][tx] = board[y][x]
			board[y][x] = 0
		}
		moved += len(toMove)
		toMove = nil

		for y, row := range board {
			for x, v := range row {
				ty := (y + 1) % len(board)
				if v == 2 && board[ty][x] == 0 {
					toMove = append(toMove, [4]int{x, y, x, ty})
				}
			}
		}
		for _, v := range toMove {
			x, y, tx, ty := v[0], v[1], v[2], v[3]
			board[ty][tx] = board[y][x]
			board[y][x] = 0
		}
		moved += len(toMove)
		toMove = nil
	}

	return steps, nil
}

func AnswerV1(input string) (int, error) {

	return 0, nil
}
