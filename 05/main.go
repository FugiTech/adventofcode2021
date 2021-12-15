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
	var vents [1000][1000]uint8

	lines := strings.Split(input, "\n")
	for _, l := range lines {
		points := strings.Split(l, " -> ")
		p0 := strings.Split(points[0], ",")
		p1 := strings.Split(points[1], ",")
		x1, _ := strconv.Atoi(p0[0])
		y1, _ := strconv.Atoi(p0[1])
		x2, _ := strconv.Atoi(p1[0])
		y2, _ := strconv.Atoi(p1[1])
		if x1 == x2 {
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			for i := y1; i <= y2; i++ {
				vents[i][x1]++
			}
		} else if y1 == y2 {
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			for i := x1; i <= x2; i++ {
				vents[y1][i]++
			}
		}
	}

	count := 0
	for _, r := range vents {
		for _, v := range r {
			if v >= 2 {
				count++
			}
		}
	}

	return count, nil
}

func AnswerV1(input string) (int, error) {
	var vents [1000][1000]uint8

	lines := strings.Split(input, "\n")
	for _, l := range lines {
		points := strings.Split(l, " -> ")
		p0 := strings.Split(points[0], ",")
		p1 := strings.Split(points[1], ",")
		x1, _ := strconv.Atoi(p0[0])
		y1, _ := strconv.Atoi(p0[1])
		x2, _ := strconv.Atoi(p1[0])
		y2, _ := strconv.Atoi(p1[1])

		xD := 0
		if x2 > x1 {
			xD = 1
		} else if x2 < x1 {
			xD = -1
		}
		yD := 0
		if y2 > y1 {
			yD = 1
		} else if y2 < y1 {
			yD = -1
		}

		//fmt.Printf("%d,%d -> %d,%d (%d,%d)\n", x1, y1, x2, y2, xD, yD)
		for x, y := x1, y1; x != x2 || y != y2; x, y = x+xD, y+yD {
			vents[y][x]++
		}
		vents[y2][x2]++
	}

	count := 0
	for _, r := range vents {
		//fmt.Printf("%v\n", r)
		for _, v := range r {
			if v >= 2 {
				count++
			}
		}
	}

	return count, nil
}

func AnswerV2(input string) (int, error) {
	var vents [1000][1000]uint8

	for i, l := 0, len(input); i < l; i++ {
		var x1, y1, x2, y2 int16
		for ; input[i] != ','; i++ {
			x1 = 10*x1 + int16(input[i]-'0')
		}
		for i++; input[i] != ' '; i++ {
			y1 = 10*y1 + int16(input[i]-'0')
		}
		for i += 4; input[i] != ','; i++ {
			x2 = 10*x2 + int16(input[i]-'0')
		}
		for i++; i < l && input[i] != '\n'; i++ {
			y2 = 10*y2 + int16(input[i]-'0')
		}
		xD := normalizeSlow(x2 - x1)
		yD := normalizeSlow(y2 - y1)

		for x, y := x1, y1; x != x2 || y != y2; x, y = x+xD, y+yD {
			vents[y][x]++
		}
		vents[y2][x2]++
	}

	count := 0
	for _, r := range vents {
		for _, v := range r {
			if v >= 2 {
				count++
			}
		}
	}

	return count, nil
}

func AnswerV3(input string) (int, error) {
	var vents [1000][1000]uint8
	count := 0

	for i, l := 0, len(input); i < l; i++ {
		var x1, y1, x2, y2 int16
		for ; input[i] != ','; i++ {
			x1 = 10*x1 + int16(input[i]-'0')
		}
		for i++; input[i] != ' '; i++ {
			y1 = 10*y1 + int16(input[i]-'0')
		}
		for i += 4; input[i] != ','; i++ {
			x2 = 10*x2 + int16(input[i]-'0')
		}
		for i++; i < l && input[i] != '\n'; i++ {
			y2 = 10*y2 + int16(input[i]-'0')
		}
		xD := normalizeSlow(x2 - x1)
		yD := normalizeSlow(y2 - y1)

		for x, y := x1, y1; x != x2 || y != y2; x, y = x+xD, y+yD {
			vents[y][x]++
			if vents[y][x] == 2 {
				count++
			}
		}
		vents[y2][x2]++
		if vents[y2][x2] == 2 {
			count++
		}
	}

	return count, nil
}

func AnswerV4(input string) (int, error) {
	vents := make(map[int32]int32, 32*1024)
	count := 0

	for i, l := 0, len(input); i < l; i++ {
		var x1, y1, x2, y2 int16
		for ; input[i] != ','; i++ {
			x1 = 10*x1 + int16(input[i]-'0')
		}
		for i++; input[i] != ' '; i++ {
			y1 = 10*y1 + int16(input[i]-'0')
		}
		for i += 4; input[i] != ','; i++ {
			x2 = 10*x2 + int16(input[i]-'0')
		}
		for i++; i < l && input[i] != '\n'; i++ {
			y2 = 10*y2 + int16(input[i]-'0')
		}
		xD := normalizeSlow(x2 - x1)
		yD := normalizeSlow(y2 - y1)

		for x, y := x1, y1; x != x2 || y != y2; x, y = x+xD, y+yD {
			k := int32(y)<<16 | int32(x)
			vents[k]++
			if vents[k] == 2 {
				count++
			}
		}
		k := int32(y2)<<16 | int32(x2)
		vents[k]++
		if vents[k] == 2 {
			count++
		}
	}

	return count, nil
}

type line struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

func AnswerV5(input string) (int, error) {
	lines := make([]line, 0, 500)
	intersections := make(map[int]struct{}, 32000)

	for i, j, l := 0, 0, len(input); i < l; i, j = i+1, j+1 {
		p := line{}
		for ; input[i] != ','; i++ {
			p.x1 = 10*p.x1 + int(input[i]-'0')
		}
		for i++; input[i] != ' '; i++ {
			p.y1 = 10*p.y1 + int(input[i]-'0')
		}
		for i += 4; input[i] != ','; i++ {
			p.x2 = 10*p.x2 + int(input[i]-'0')
		}
		for i++; i < l && input[i] != '\n'; i++ {
			p.y2 = 10*p.y2 + int(input[i]-'0')
		}
		lines = append(lines, p)
	}

	for i, l1 := range lines[:] {
		for _, l2 := range lines[i+1:] {
			xd1, yd1 := l1.x1-l1.x2, l1.y1-l1.y2
			xd2, yd2 := l2.x1-l2.x2, l2.y1-l2.y2
			div := xd1*yd2 - xd2*yd1
			if div == 0 { // lines are parallel, but maybe still intersecting
				/*
					xd, yd := normalizeSlow64(l1.x2-l1.x1), normalizeSlow64(l1.y2-l1.y1)
					for x, y := l1.x1, l1.y1; x != l1.x2+xd || y != l1.y2+yd; x, y = x+xd, y+yd {
						if (x >= l2.x1 || x >= l2.x2) && (x <= l2.x1 || x <= l2.x2) && (y >= l2.y1 || y >= l2.y2) && (y <= l2.y1 || y <= l2.y2) {
							fmt.Printf("[P] %d,%d -> %d,%d (%d,%d) %d,%d -> %d,%d\n", l1.x1, l1.y1, l1.x2, l1.y2, x, y, l2.x1, l2.y1, l2.x2, l2.y2)
							intersections[(x<<16)|y] = struct{}{}
						}
					}
				*/
				continue
			}

			d1 := l1.x1*l1.y2 - l1.x2*l1.y1
			d2 := l2.x1*l2.y2 - l2.x2*l2.y1
			x := (d1*xd2 - d2*xd1) / div
			y := (d1*yd2 - d2*yd1) / div

			valid := true
			valid = valid && (x >= l1.x1 || x >= l1.x2) && (x <= l1.x1 || x <= l1.x2)
			valid = valid && (x >= l2.x1 || x >= l2.x2) && (x <= l2.x1 || x <= l2.x2)
			valid = valid && (y >= l1.y1 || y >= l1.y2) && (y <= l1.y1 || y <= l1.y2)
			valid = valid && (y >= l2.y1 || y >= l2.y2) && (y <= l2.y1 || y <= l2.y2)

			if !valid {
				continue
			}
			// fmt.Printf("[I] %d,%d -> %d,%d (%d,%d) %d,%d -> %d,%d\n", l1.x1, l1.y1, l1.x2, l1.y2, x, y, l2.x1, l2.y1, l2.x2, l2.y2)

			intersections[(x<<16)|y] = struct{}{}
		}
	}

	return len(intersections), nil
}

// normalize turns c into 1 (if positive), 0 (if zero), or -1 (if negative)
func normalize(c int16) int16 {
	s := (c >> 14) & 2           // 2 if c is negative, 0 otherwise
	m := (c >> 8) | (c & 0x00FF) // merge first 8 bits with second 8
	m = (m >> 4) | (m & 0x000F)  // merge 4 into 4
	m = (m >> 2) | (m & 0x0003)  // merge 2 into 2
	m = (m >> 1) | (m & 0x0001)  // merge 1 into 1, m is now 1 or 0
	return m - s                 // (1 or 0) - (2 or 0) = 1, 0, or -1
}

func normalizeSlow(c int16) int16 {
	if c > 0 {
		return 1
	} else if c < 0 {
		return -1
	} else {
		return 0
	}
}

func normalizeSlow64(c int) int {
	if c > 0 {
		return 1
	} else if c < 0 {
		return -1
	} else {
		return 0
	}
}
