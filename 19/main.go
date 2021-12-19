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

type scanner struct {
	B []beacon // beacons
	N bool     // normalized?
}
type beacon [3]int

func AnswerPart1(input string) (int, error) {
	sections := strings.Split(input, "\n\n")
	scanners := make([]*scanner, len(sections))

	for i, s := range sections {
		lines := strings.Split(s, "\n")
		scanners[i] = &scanner{
			B: make([]beacon, len(lines)-1),
		}
		for j, l := range lines[1:] {
			p := strings.Split(l, ",")
			scanners[i].B[j][0], _ = strconv.Atoi(p[0])
			scanners[i].B[j][1], _ = strconv.Atoi(p[1])
			scanners[i].B[j][2], _ = strconv.Atoi(p[2])
		}
	}

	trueBeacons := map[[3]int]struct{}{}
	toSearch := []*scanner{scanners[0]}

	scanners[0].N = true
	for _, b := range scanners[0].B {
		trueBeacons[b] = struct{}{}
	}

	for len(toSearch) > 0 {
		b := toSearch[0] // For each normalized scanner as a "base"
		toSearch = toSearch[1:]
		for _, c := range scanners {
			if c.N { // Compare against only un-normalized scanners
				continue
			}

			counts := map[[4]int]int{}
			for _, beaconC := range c.B {
				r := rotations(beaconC)
				for _, beaconB := range b.B {
					for o := 0; o < 24; o++ { // Try each orientation for the "comparison"
						counts[diff(beaconB, r[o], o)]++
					}
				}
			}
			for k, v := range counts {
				if v >= 12 {
					c.N = true
					for i, beacon := range c.B {
						r := rotations(beacon)[k[3]]
						c.B[i][0] = r[0] + k[0]
						c.B[i][1] = r[1] + k[1]
						c.B[i][2] = r[2] + k[2]
						trueBeacons[c.B[i]] = struct{}{}
					}
					toSearch = append(toSearch, c)

					break
				}
			}
		}
	}

	return len(trueBeacons), nil
}

func AnswerV1(input string) (int, error) {
	sections := strings.Split(input, "\n\n")
	scanners := make([]*scanner, len(sections))

	for i, s := range sections {
		lines := strings.Split(s, "\n")
		scanners[i] = &scanner{
			B: make([]beacon, len(lines)-1),
		}
		for j, l := range lines[1:] {
			p := strings.Split(l, ",")
			scanners[i].B[j][0], _ = strconv.Atoi(p[0])
			scanners[i].B[j][1], _ = strconv.Atoi(p[1])
			scanners[i].B[j][2], _ = strconv.Atoi(p[2])
		}
	}

	scannerLocs := map[[3]int]struct{}{}
	toSearch := []*scanner{scanners[0]}

	scanners[0].N = true
	scannerLocs[[3]int{0, 0, 0}] = struct{}{}

	for len(toSearch) > 0 {
		b := toSearch[0] // For each normalized scanner as a "base"
		toSearch = toSearch[1:]
		for _, c := range scanners {
			if c.N { // Compare against only un-normalized scanners
				continue
			}

			counts := map[[4]int]int{}
			for _, beaconC := range c.B {
				r := rotations(beaconC)
				for _, beaconB := range b.B {
					for o := 0; o < 24; o++ { // Try each orientation for the "comparison"
						counts[diff(beaconB, r[o], o)]++
					}
				}
			}
			for k, v := range counts {
				if v >= 12 {
					c.N = true
					for i, beacon := range c.B {
						r := rotations(beacon)[k[3]]
						c.B[i][0] = r[0] + k[0]
						c.B[i][1] = r[1] + k[1]
						c.B[i][2] = r[2] + k[2]
					}
					toSearch = append(toSearch, c)
					scannerLocs[[3]int{k[0], k[1], k[2]}] = struct{}{}

					break
				}
			}
		}
	}

	D := 0
	for a := range scannerLocs {
		for b := range scannerLocs {
			d := distance(a, b)
			if d > D {
				D = d
			}
		}
	}

	return D, nil
}

func diff(a, b [3]int, o int) (c [4]int) {
	c[0] = a[0] - b[0]
	c[1] = a[1] - b[1]
	c[2] = a[2] - b[2]
	c[3] = o
	return
}

func distance(a, b [3]int) int {
	x, y, z := a[0]-b[0], a[1]-b[1], a[2]-b[2]
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}
	if z < 0 {
		z = -z
	}
	return x + y + z
}

func rotations(in [3]int) (out [24][3]int) {
	out[0][0] = in[0]
	out[0][1] = in[1]
	out[0][2] = in[2]

	out[1][0] = in[0]
	out[1][1] = -in[2]
	out[1][2] = in[1]

	out[2][0] = in[0]
	out[2][1] = -in[1]
	out[2][2] = -in[2]

	out[3][0] = in[0]
	out[3][1] = in[2]
	out[3][2] = -in[1]

	out[4][0] = -in[1]
	out[4][1] = in[0]
	out[4][2] = in[2]

	out[5][0] = in[2]
	out[5][1] = in[0]
	out[5][2] = in[1]

	out[6][0] = in[1]
	out[6][1] = in[0]
	out[6][2] = -in[2]

	out[7][0] = -in[2]
	out[7][1] = in[0]
	out[7][2] = -in[1]

	out[8][0] = -in[0]
	out[8][1] = -in[1]
	out[8][2] = in[2]

	out[9][0] = -in[0]
	out[9][1] = -in[2]
	out[9][2] = -in[1]

	out[10][0] = -in[0]
	out[10][1] = in[1]
	out[10][2] = -in[2]

	out[11][0] = -in[0]
	out[11][1] = in[2]
	out[11][2] = in[1]

	out[12][0] = in[1]
	out[12][1] = -in[0]
	out[12][2] = in[2]

	out[13][0] = in[2]
	out[13][1] = -in[0]
	out[13][2] = -in[1]

	out[14][0] = -in[1]
	out[14][1] = -in[0]
	out[14][2] = -in[2]

	out[15][0] = -in[2]
	out[15][1] = -in[0]
	out[15][2] = in[1]

	out[16][0] = -in[2]
	out[16][1] = in[1]
	out[16][2] = in[0]

	out[17][0] = in[1]
	out[17][1] = in[2]
	out[17][2] = in[0]

	out[18][0] = in[2]
	out[18][1] = -in[1]
	out[18][2] = in[0]

	out[19][0] = -in[1]
	out[19][1] = -in[2]
	out[19][2] = in[0]

	out[20][0] = -in[2]
	out[20][1] = -in[1]
	out[20][2] = -in[0]

	out[21][0] = -in[1]
	out[21][1] = in[2]
	out[21][2] = -in[0]

	out[22][0] = in[2]
	out[22][1] = in[1]
	out[22][2] = -in[0]

	out[23][0] = in[1]
	out[23][1] = -in[2]
	out[23][2] = -in[0]

	return
}
