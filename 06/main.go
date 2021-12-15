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

var Answer = AnswerV5

func AnswerPart1(input string) (int, error) {
	var fish [10]int

	inNums := strings.Split(input, ",")
	for _, n := range inNums {
		v, _ := strconv.Atoi(n)
		fish[v]++
	}

	for i := 0; i < 80; i++ {
		fish[9] = fish[0]
		fish[7] += fish[0]
		for j := 0; j < 9; j++ {
			fish[j] = fish[j+1]
		}
	}

	count := 0
	for i := 0; i < 9; i++ {
		count += fish[i]
	}

	return count, nil
}

func AnswerV1(input string) (int, error) {
	var fish [10]int

	inNums := strings.Split(input, ",")
	for _, n := range inNums {
		v, _ := strconv.Atoi(n)
		fish[v]++
	}

	for i := 0; i < 256; i++ {
		fish[9] = fish[0]
		fish[7] += fish[0]
		for j := 0; j < 9; j++ {
			fish[j] = fish[j+1]
		}
	}

	count := 0
	for i := 0; i < 9; i++ {
		count += fish[i]
	}

	return count, nil
}

func AnswerV2(input string) (int, error) {
	var fish [10]int

	for i, l := 0, len(input); i < l; i += 2 {
		fish[input[i]-'0']++
	}

	for i := 0; i < 256; i++ {
		fish[9] = fish[0]
		fish[7] += fish[0]
		for j := 0; j < 9; j++ {
			fish[j] = fish[j+1]
		}
	}

	count := 0
	for i := 0; i < 9; i++ {
		count += fish[i]
	}

	return count, nil
}

func AnswerV3(input string) (int, error) {
	var adults [7]int
	var babies [9]int

	for i, l := 0, len(input); i < l; i += 2 {
		adults[input[i]-'0']++
	}

	for i := 0; i < 256; i++ {
		a, b := i%7, i%9
		adults[a] += babies[b]
		babies[b] = adults[a]
	}

	count := 0
	for i := 0; i < 7; i++ {
		count += adults[i]
	}
	for i := 0; i < 9; i++ {
		count += babies[i]
	}

	return count, nil
}

func AnswerV4(input string) (int, error) {
	var adults [7]int
	var babies [9]int

	for i, l := 0, len(input); i < l; i += 2 {
		adults[input[i]-'0']++
	}

	for i, a, b := 0, int8(0), int8(0); i < 256; i, a, b = i+1, a+1, b+1 {
		a = ((a - 7) >> 7) & a
		b = ((b - 9) >> 7) & b
		c := adults[a] + babies[b]
		adults[a] = c
		babies[b] = c
	}

	count := 0
	for i := 0; i < 7; i++ {
		count += adults[i]
	}
	for i := 0; i < 9; i++ {
		count += babies[i]
	}

	return count, nil
}

func AnswerV5(input string) (int, error) {
	var adults [7]int
	var babies [9]int

	for i, l := 0, len(input); i < l; i += 2 {
		adults[input[i]-'0']++
	}

	for i, a, b := 0, int8(0), int8(0); i < 256; i, a, b = i+1, a+1, b+1 {
		a = ((a - 7) >> 7) & a
		b = ((b - 9) >> 7) & b
		adults[a] += babies[b]
		babies[b] = adults[a]
	}

	count := adults[0] + adults[1] + adults[2] + adults[3] + adults[4] + adults[5] + adults[6] +
		babies[0] + babies[1] + babies[2] + babies[3] + babies[4] + babies[5] + babies[6] + babies[7] + babies[8]

	return count, nil
}
