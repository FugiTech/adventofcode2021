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
	sections := strings.Split(input, "\n\n")
	lines := strings.Split(sections[1], "\n")
	H, W := len(lines), len(lines[0])
	MaxSteps := 50

	enhancement := [512]int{}
	image := [300][300]int{}

	printTest := func(step int) {
		for y := 100 - step; y < 100+H+step; y++ {
			for x := 100 - step; x < 100+W+step; x++ {
				if image[y][x] == 1 {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}
	_ = printTest

	for i, c := range []byte(sections[0]) {
		if c == '#' {
			enhancement[i] = 1
		}
	}

	for y, l := range lines {
		for x, c := range []byte(l) {
			if c == '#' {
				image[y+100][x+100] = 1
			}
		}
	}

	//printTest(0)
	for step := 1; step <= MaxSteps; step++ {
		n := [300][300]int{}
		for y := 1; y < 299; y++ {
			for x := 1; x < 299; x++ {
				k := ((image[y-1][x-1] << 8) |
					(image[y-1][x] << 7) |
					(image[y-1][x+1] << 6) |
					(image[y][x-1] << 5) |
					(image[y][x] << 4) |
					(image[y][x+1] << 3) |
					(image[y+1][x-1] << 2) |
					(image[y+1][x] << 1) |
					(image[y+1][x+1]))
				n[y][x] = enhancement[k]
			}
		}
		image = n
		//printTest(step)
	}

	count := 0
	for _, r := range image[100-MaxSteps : 100+H+MaxSteps] {
		for _, c := range r[100-MaxSteps : 100+W+MaxSteps] {
			count += c
		}
	}

	return count, nil
}

func AnswerV1(input string) (int, error) {

	return 0, nil
}
