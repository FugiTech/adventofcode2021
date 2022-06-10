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
	var octopi [10][10]int
	lines := strings.Split(input, "\n")
	flashes := 0

	for i, l := range lines {
		for j, c := range []byte(l) {
			octopi[i][j] = int(c - '0')
		}
	}

	var inc func(i, j int)
	inc = func(i, j int) {
		if i < 0 || j < 0 || i > 9 || j > 9 {
			return
		}

		octopi[i][j]++
		if octopi[i][j] == 10 {
			flashes++
			inc(i-1, j-1)
			inc(i-1, j)
			inc(i-1, j+1)
			inc(i, j-1)
			inc(i, j+1)
			inc(i+1, j-1)
			inc(i+1, j)
			inc(i+1, j+1)
		}
	}

	for step := 0; step < 100; step++ {
		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				inc(i, j)
			}
		}

		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				if octopi[i][j] > 9 {
					octopi[i][j] = 0
				}
			}
		}
	}

	return flashes, nil
}

func AnswerV1(input string) (int, error) {
	var octopi [10][10]int
	lines := strings.Split(input, "\n")
	flashes := 0

	for i, l := range lines {
		for j, c := range []byte(l) {
			octopi[i][j] = int(c - '0')
		}
	}

	var inc func(i, j int)
	inc = func(i, j int) {
		if i < 0 || j < 0 || i > 9 || j > 9 {
			return
		}

		octopi[i][j]++
		if octopi[i][j] == 10 {
			flashes++
			inc(i-1, j-1)
			inc(i-1, j)
			inc(i-1, j+1)
			inc(i, j-1)
			inc(i, j+1)
			inc(i+1, j-1)
			inc(i+1, j)
			inc(i+1, j+1)
		}
	}

	for step := 1; step < 10000; step++ {
		flashes = 0

		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				inc(i, j)
			}
		}

		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				if octopi[i][j] > 9 {
					octopi[i][j] = 0
				}
			}
		}

		if flashes == 100 {
			return step, nil
		}
	}

	return 0, fmt.Errorf("step greater than 10,000")
}
