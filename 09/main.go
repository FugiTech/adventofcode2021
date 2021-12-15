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
	m := make([][]int, len(lines))

	for idx, l := range lines {
		ml := make([]int, len(l))
		for cidx, c := range []byte(l) {
			ml[cidx] = int(c - '0')
		}
		m[idx] = ml
	}

	risk := 0
	for r, lr := 0, len(m); r < lr; r++ {
		for c, lc := 0, len(m[r]); c < lc; c++ {
			v := m[r][c]
			if (c == 0 || v < m[r][c-1]) && (r == 0 || v < m[r-1][c]) && (c == lc-1 || v < m[r][c+1]) && (r == lr-1 || v < m[r+1][c]) {
				risk += v + 1
			}
		}
	}

	return risk, nil
}

func AnswerV1(input string) (int, error) {
	lines := strings.Split(input, "\n")
	heights := make([][]int, len(lines))
	basins := make([][]int, len(lines))

	for idx, l := range lines {
		ml := make([]int, len(l))
		for cidx, c := range []byte(l) {
			ml[cidx] = int(c - '0')
		}
		heights[idx] = ml
		basins[idx] = make([]int, len(l))
	}

	basinIdx := 1
	for r, lr := 0, len(heights); r < lr; r++ {
		for c, lc := 0, len(heights[r]); c < lc; c++ {
			v := heights[r][c]
			if (c == 0 || v < heights[r][c-1]) && (r == 0 || v < heights[r-1][c]) && (c == lc-1 || v < heights[r][c+1]) && (r == lr-1 || v < heights[r+1][c]) {
				basins[r][c] = basinIdx
				basinIdx++
			}
		}
	}

	goal := 9 * len(heights) * len(heights[0])
	for sum := 0; sum != goal; {
		sum = 0
		for r, lr := 0, len(heights); r < lr; r++ {
			for c, lc := 0, len(heights[r]); c < lc; c++ {
				if heights[r][c] == 9 {
					// Not part of a basin
				} else if basins[r][c] != 0 {
					heights[r][c] = 9
				} else if r != 0 && basins[r-1][c] != 0 {
					basins[r][c] = basins[r-1][c]
					heights[r][c] = 9
				} else if c != 0 && basins[r][c-1] != 0 {
					basins[r][c] = basins[r][c-1]
					heights[r][c] = 9
				} else if r != lr-1 && basins[r+1][c] != 0 {
					basins[r][c] = basins[r+1][c]
					heights[r][c] = 9
				} else if c != lc-1 && basins[r][c+1] != 0 {
					basins[r][c] = basins[r][c+1]
					heights[r][c] = 9
				}
				sum += heights[r][c]
			}
		}
	}

	basinSizes := make([]int, basinIdx)
	for r, lr := 0, len(basins); r < lr; r++ {
		for c, lc := 0, len(basins[r]); c < lc; c++ {
			basinSizes[basins[r][c]]++
		}
	}
	basinSizes[0] = 0

	sort.Sort(sort.Reverse(sort.IntSlice(basinSizes)))

	return basinSizes[0] * basinSizes[1] * basinSizes[2], nil
}
