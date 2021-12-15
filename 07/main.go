package main

import (
	"fmt"
	"io/ioutil"
	"math"
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

var Answer = AnswerV1

func AnswerPart1(input string) (int, error) {
	in := make([]int, 0, 1000)

	inNums := strings.Split(input, ",")
	for _, n := range inNums {
		v, _ := strconv.Atoi(n)
		in = append(in, v)
	}

	sort.Sort(sort.IntSlice(in))

	median := in[len(in)/2]
	fuel := 0

	for _, v := range in {
		r := v - median
		if r < 0 {
			r *= -1
		}
		fuel += r
	}

	return fuel, nil
}

func AnswerV1(input string) (int, error) {
	in := make([]int, 0, 1000)
	total := 0

	inNums := strings.Split(input, ",")
	for _, n := range inNums {
		v, _ := strconv.Atoi(n)
		in = append(in, v)
		total += v
	}

	average := int(math.Round(float64(total) / float64(len(in))))
	const guesses = 3
	var answers [guesses]int

	for _, v := range in {
		for k := 0; k < guesses; k++ {
			r := v - (average - guesses/2 + k)
			if r < 0 {
				r *= -1
			}
			for i := 1; i <= r; i++ {
				answers[k] += i
			}
		}
	}

	fmt.Println(total, average, answers)

	min := math.MaxInt
	for k := 0; k < guesses; k++ {
		if answers[k] < min {
			min = answers[k]
		}
	}

	return min, nil
}
