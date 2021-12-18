package main

import (
	"fmt"
	"io/ioutil"
	"math"
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
	p := strings.Split(input[15:], ", y=")
	p1 := strings.Split(p[0], "..")
	p2 := strings.Split(p[1], "..")
	x1, _ := strconv.Atoi(p1[0])
	x2, _ := strconv.Atoi(p1[1])
	y1, _ := strconv.Atoi(p2[0])
	y2, _ := strconv.Atoi(p2[1])

	X := int(math.Ceil(SolveForX(x1)))

	Y, MaxY := 0, 0
	for i := 0; i < 100; i++ {
		y, maxy, d := 0, 0, i
		for step := 0; y > y2; step++ {
			y += d
			d--
			if y > maxy {
				maxy = y
			}
			if y <= y2 && y >= y1 && step >= X-1 {
				Y = i
				MaxY = maxy
				break
			}
		}
	}

	_, _ = x2, Y

	return MaxY, nil
}

func AnswerV1(input string) (int, error) {
	p := strings.Split(input[15:], ", y=")
	p1 := strings.Split(p[0], "..")
	p2 := strings.Split(p[1], "..")
	x1, _ := strconv.Atoi(p1[0])
	x2, _ := strconv.Atoi(p1[1])
	y1, _ := strconv.Atoi(p2[1])
	y2, _ := strconv.Atoi(p2[0])

	valid := 0
	for ix := -10000; ix < 10000; ix++ {
		for iy := -10000; iy < 10000; iy++ {
			x, y, dx, dy := 0, 0, ix, iy
			for step := 0; x <= x2 && y >= y2; step++ {
				x += dx
				y += dy
				dy--
				if dx < 0 {
					dx++
				} else if dx > 0 {
					dx--
				}
				if x1 <= x && x <= x2 && y1 >= y && y >= y2 {
					valid++
					break
				}
			}

		}
	}

	return valid, nil
}

func SolveForX(target int) float64 {
	return (math.Sqrt(8*float64(target)+1) - 1) / 2
}
