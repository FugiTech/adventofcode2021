package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

var Answer = AnswerV1

func AnswerPart1(input string) (int, error) {
	var paper [1400][1400]bool

	printPaper := func() {
		for y := 0; y < 15; y++ {
			for x := 0; x < 15; x++ {
				if paper[y][x] {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		}
	}
	_ = printPaper

	instructions := strings.Split(input, "\n\n")
	dots := strings.Split(instructions[0], "\n")
	folds := strings.Split(instructions[1], "\n")

	for _, l := range dots {
		p := strings.Split(l, ",")
		x, _ := strconv.Atoi(p[0])
		y, _ := strconv.Atoi(p[1])
		paper[y][x] = true
	}

	for _, l := range folds {
		p := strings.Split(l[11:], "=")
		axis := p[0]
		idx, _ := strconv.Atoi(p[1])

		maxX, maxY := 1400, 1400
		if axis == "x" {
			maxX = idx
		} else {
			maxY = idx
		}

		var n [1400][1400]bool

		for y := 0; y < maxY; y++ {
			for x := 0; x < maxX; x++ {
				n[y][x] = paper[y][x]
			}
		}

		if axis == "x" {
			for y := 0; y < 1400; y++ {
				for d, x := idx-1, idx+1; d >= 0 && x < 1400; d, x = d-1, x+1 {
					n[y][d] = n[y][d] || paper[y][x]
				}
			}
		} else {
			for d, y := idx-1, idx+1; d >= 0 && y < 1400; d, y = d-1, y+1 {
				for x := 0; x < 1400; x++ {
					n[d][x] = n[d][x] || paper[y][x]
				}
			}
		}

		paper = n

		break
	}

	count := 0
	for y := 0; y < 1400; y++ {
		for x := 0; x < 1400; x++ {
			if paper[y][x] {
				count++
			}
		}
	}

	return count, nil
}

func AnswerV1(input string) (int, error) {
	var paper [1400][1400]bool

	printPaper := func() {
		for y := 0; y < 15; y++ {
			for x := 0; x < 150; x++ {
				if paper[y][x] {
					fmt.Print("#")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Println()
		}
	}
	_ = printPaper

	instructions := strings.Split(input, "\n\n")
	dots := strings.Split(instructions[0], "\n")
	folds := strings.Split(instructions[1], "\n")

	for _, l := range dots {
		p := strings.Split(l, ",")
		x, _ := strconv.Atoi(p[0])
		y, _ := strconv.Atoi(p[1])
		paper[y][x] = true
	}

	for _, l := range folds {
		p := strings.Split(l[11:], "=")
		axis := p[0]
		idx, _ := strconv.Atoi(p[1])

		maxX, maxY := 1400, 1400
		if axis == "x" {
			maxX = idx
		} else {
			maxY = idx
		}

		var n [1400][1400]bool

		for y := 0; y < maxY; y++ {
			for x := 0; x < maxX; x++ {
				n[y][x] = paper[y][x]
			}
		}

		if axis == "x" {
			for y := 0; y < 1400; y++ {
				for d, x := idx-1, idx+1; d >= 0 && x < 1400; d, x = d-1, x+1 {
					n[y][d] = n[y][d] || paper[y][x]
				}
			}
		} else {
			for d, y := idx-1, idx+1; d >= 0 && y < 1400; d, y = d-1, y+1 {
				for x := 0; x < 1400; x++ {
					n[d][x] = n[d][x] || paper[y][x]
				}
			}
		}

		paper = n
	}

	count := 0
	for y := 0; y < 1400; y++ {
		for x := 0; x < 1400; x++ {
			if paper[y][x] {
				count++
			}
		}
	}

	printPaper()

	return count, nil
}
