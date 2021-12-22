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
	cubes := [101][101][101]bool{}
	MIN, MAX := -50, 50

	lines := strings.Split(input, "\n")
	for _, l := range lines {
		s1 := strings.Split(l, " ")
		s2 := strings.Split(s1[1], ",")
		s31 := strings.Split(s2[0][2:], "..")
		s32 := strings.Split(s2[1][2:], "..")
		s33 := strings.Split(s2[2][2:], "..")

		state := s1[0] == "on"
		x1, _ := strconv.Atoi(s31[0])
		x2, _ := strconv.Atoi(s31[1])
		y1, _ := strconv.Atoi(s32[0])
		y2, _ := strconv.Atoi(s32[1])
		z1, _ := strconv.Atoi(s33[0])
		z2, _ := strconv.Atoi(s33[1])
		if x1 <= MIN {
			x1 = MIN
		}
		if y1 <= MIN {
			y1 = MIN
		}
		if z1 <= MIN {
			z1 = MIN
		}

		for x := x1; x <= x2 && x <= MAX; x++ {
			for y := y1; y <= y2 && y <= MAX; y++ {
				for z := z1; z <= z2 && z <= MAX; z++ {
					cubes[x+50][y+50][z+50] = state
				}
			}
		}
	}

	count := 0
	for _, r1 := range cubes {
		for _, r2 := range r1 {
			for _, v := range r2 {
				if v {
					count++
				}
			}
		}
	}

	return count, nil
}

type cube struct {
	x, y, z int
	w, h, d int
}

func AnswerV1(input string) (int, error) {
	cubes := []cube{}

	lines := strings.Split(input, "\n")
	for _, l := range lines {
		s1 := strings.Split(l, " ")
		s2 := strings.Split(s1[1], ",")
		s31 := strings.Split(s2[0][2:], "..")
		s32 := strings.Split(s2[1][2:], "..")
		s33 := strings.Split(s2[2][2:], "..")

		state := s1[0] == "on"
		x1, _ := strconv.Atoi(s31[0])
		x2, _ := strconv.Atoi(s31[1])
		y1, _ := strconv.Atoi(s32[0])
		y2, _ := strconv.Atoi(s32[1])
		z1, _ := strconv.Atoi(s33[0])
		z2, _ := strconv.Atoi(s33[1])

		C := cube{
			x: x1,
			y: y1,
			z: z1,
			w: x2 - x1 + 1,
			h: y2 - y1 + 1,
			d: z2 - z1 + 1,
		}

		newCubes := []cube{}
		for _, c := range cubes {
			newCubes = append(newCubes, splitCube(c, C)...)
		}
		if state {
			newCubes = append(newCubes, C)
		}
		cubes = newCubes
	}

	count := 0
	for _, c := range cubes {
		count += c.w * c.h * c.d
	}

	return count, nil
}

// splitCube returns a list of cubes that are within "a" but NOT within "b"
func splitCube(a, b cube) []cube {
	// If they don't intersect at all, just return the existing cube
	if (a.x+a.w) <= b.x || (b.x+b.w) <= a.x || (a.y+a.h) <= b.y || (b.y+b.h) <= a.y || (a.z+a.d) <= b.z || (b.z+b.d) <= a.z {
		return []cube{a}
	}

	// Otherwise split in "half" up to 6 times to form 1-7 cubes, of which 0-6 get returned and the 1 remainder is the intersection
	r := make([]cube, 0, 6)

	// Split x dimension
	if a.x < b.x {
		// the left side is saved
		n := cube{
			x: a.x,
			y: a.y,
			z: a.z,
			w: b.x - a.x,
			h: a.h,
			d: a.d,
		}
		r = append(r, n)
		a.x = b.x
		a.w -= n.w
	}
	if (a.x + a.w) > (b.x + b.w) {
		// the right side is saved
		n := cube{
			x: b.x + b.w,
			y: a.y,
			z: a.z,
			w: (a.x + a.w) - (b.x + b.w),
			h: a.h,
			d: a.d,
		}
		r = append(r, n)
		a.w -= n.w
	}

	// Split y dimension
	if a.y < b.y {
		// the bottom side is saved
		n := cube{
			x: a.x,
			y: a.y,
			z: a.z,
			w: a.w,
			h: b.y - a.y,
			d: a.d,
		}
		r = append(r, n)
		a.y = b.y
		a.h -= n.h
	}
	if (a.y + a.h) > (b.y + b.h) {
		// the top side is saved
		n := cube{
			x: a.x,
			y: b.y + b.h,
			z: a.z,
			w: a.w,
			h: (a.y + a.h) - (b.y + b.h),
			d: a.d,
		}
		r = append(r, n)
		a.h -= n.h
	}

	// Split z dimension
	if a.z < b.z {
		// the close side is saved
		n := cube{
			x: a.x,
			y: a.y,
			z: a.z,
			w: a.w,
			h: a.h,
			d: b.z - a.z,
		}
		r = append(r, n)
		a.z = b.z
		a.d -= n.d
	}
	if (a.z + a.d) > (b.z + b.d) {
		// the far side is saved
		n := cube{
			x: a.x,
			y: a.y,
			z: b.z + b.d,
			w: a.w,
			h: a.h,
			d: (a.z + a.d) - (b.z + b.d),
		}
		r = append(r, n)
		a.d -= n.d
	}

	return r
}
