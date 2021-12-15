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
	board := [][]int{}

	lines := strings.Split(input, "\n")
	for _, l := range lines {
		b := make([]int, len(l))
		for i, c := range []byte(l) {
			b[i] = int(c - '0')
		}
		board = append(board, b)
	}

	moves := make([]position, 0, 1000)
	moves = append(moves, position{0, 0, 0})
	visited := map[int]struct{}{}

	addMove := func(x, y, risk int) {
		if x < 0 || x >= len(board[0]) {
			return
		}
		if y < 0 || y >= len(board) {
			return
		}
		k := (y << 16) | x
		if _, ok := visited[k]; ok {
			return
		}
		visited[k] = struct{}{}
		moves = append(moves, position{x, y, risk + board[y][x]})
	}

	for {
		m := moves[0]
		moves = moves[1:]
		if m.X == len(board[0])-1 && m.Y == len(board)-1 {
			return m.Risk, nil
		}
		addMove(m.X-1, m.Y, m.Risk)
		addMove(m.X+1, m.Y, m.Risk)
		addMove(m.X, m.Y-1, m.Risk)
		addMove(m.X, m.Y+1, m.Risk)
		sort.Slice(moves, func(i, j int) bool {
			return moves[i].Risk < moves[j].Risk
		})
	}
}

type position struct {
	X    int
	Y    int
	Risk int
}

func AnswerV1(input string) (int, error) {
	lines := strings.Split(input, "\n")
	board := make([][]int, len(lines)*5)

	for y, l := range lines {
		for r := 0; r < 5; r++ {
			b := make([]int, len(l)*5)
			for x, c := range []byte(l) {
				v := int(c - '0')
				for rr := 0; rr < 5; rr++ {
					b[x+len(l)*rr] = ((v + r + rr - 1) % 9) + 1
				}
			}
			board[y+len(lines)*r] = b
		}
	}

	moves := make([]position, 0, 1000)
	moves = append(moves, position{0, 0, 0})
	visited := map[int]struct{}{}

	addMove := func(x, y, risk int) {
		if x < 0 || x >= len(board[0]) {
			return
		}
		if y < 0 || y >= len(board) {
			return
		}
		k := (y << 16) | x
		if _, ok := visited[k]; ok {
			return
		}
		visited[k] = struct{}{}
		moves = append(moves, position{x, y, risk + board[y][x]})
	}

	for {
		m := moves[0]
		moves = moves[1:]
		if m.X == len(board[0])-1 && m.Y == len(board)-1 {
			return m.Risk, nil
		}
		addMove(m.X-1, m.Y, m.Risk)
		addMove(m.X+1, m.Y, m.Risk)
		addMove(m.X, m.Y-1, m.Risk)
		addMove(m.X, m.Y+1, m.Risk)
		sort.Slice(moves, func(i, j int) bool {
			return moves[i].Risk < moves[j].Risk
		})
	}
}
