package main

import (
	"fmt"
	"io/ioutil"
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

var Answer = AnswerV3

func AnswerPart1(input string) (int, error) {
	lines := strings.Split(input, "\n")
	p1, _ := strconv.Atoi(lines[0][28:])
	p2, _ := strconv.Atoi(lines[1][28:])

	dice, rolls := 0, 0
	roll := func() int {
		dice = (dice % 100) + 1
		rolls++
		return dice
	}

	p1s, p2s := 0, 0
	for p1s < 1000 && p2s < 1000 {
		p1 = ((p1 + roll() + roll() + roll() - 1) % 10) + 1
		p1s += p1
		if p1s >= 1000 {
			break
		}
		p2 = ((p2 + roll() + roll() + roll() - 1) % 10) + 1
		p2s += p2
	}

	loser := p1s
	if p2s < p1s {
		loser = p2s
	}

	return loser * rolls, nil
}

func roll(i [10]int) (o [10]int) {
	o[0] = i[9] + i[8] + i[7]
	o[1] = i[0] + i[9] + i[8]
	o[2] = i[1] + i[0] + i[9]
	o[3] = i[2] + i[1] + i[0]
	o[4] = i[3] + i[2] + i[1]
	o[5] = i[4] + i[3] + i[2]
	o[6] = i[5] + i[4] + i[3]
	o[7] = i[6] + i[5] + i[4]
	o[8] = i[7] + i[6] + i[5]
	o[9] = i[8] + i[7] + i[6]
	return o
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sum(s []int) int {
	c := 0
	for _, v := range s {
		c += v
	}
	return c
}

type state struct {
	P [2]player
}
type player struct {
	L [10]int
	S int
}

func AnswerV1(input string) (int, error) {
	lines := strings.Split(input, "\n")
	p1, _ := strconv.Atoi(lines[0][28:])
	p2, _ := strconv.Atoi(lines[1][28:])

	gameStates := make([]state, 1)
	gameStates[0].P[0].L[p1-1] = 1
	gameStates[0].P[1].L[p2-1] = 1

	active, inactive := 0, 1
	wins := [2]int{}
	for len(gameStates) > 0 {
		next := make([]state, 0, 100)

		for _, S := range gameStates {
			S.P[active].L = roll(S.P[active].L)
			S.P[active].L = roll(S.P[active].L)
			S.P[active].L = roll(S.P[active].L)

			for i, c := range S.P[active].L[:] {
				if c > 0 {
					s := state{}
					s.P[inactive] = S.P[inactive]
					s.P[active].L[i] = c
					s.P[active].S = S.P[active].S + i + 1
					if s.P[active].S >= 21 {
						other := 0
						for _, v := range S.P[inactive].L {
							other += v
						}
						wins[active] += c * other
					} else {
						next = append(next, s)
					}
				}
			}
		}

		gameStates = next
		active = (active + 1) % 2
		inactive = (inactive + 1) % 2
	}

	winner := wins[0]
	if wins[1] > winner {
		winner = wins[1]
	}

	return winner, nil
}

func AnswerV2(input string) (int, error) {
	lines := strings.Split(input, "\n")
	p1, _ := strconv.Atoi(lines[0][28:])
	p2, _ := strconv.Atoi(lines[1][28:])

	var play func(p1l, p1s, p2l, p2s, universes, activePlayer int) (p1w, p2w int)
	play = func(p1l, p1s, p2l, p2s, universes, activePlayer int) (p1w, p2w int) {
		if activePlayer == 0 {
			L := [10]int{}
			L[p1l-1] = universes
			L = roll(L)
			L = roll(L)
			L = roll(L)
			for i, c := range L {
				if c > 0 {
					s := p1s + i + 1
					if s >= 21 {
						p1w += c
					} else {
						a, b := play(i+1, s, p2l, p2s, c, 1)
						p1w += a
						p2w += b
					}
				}
			}
		} else {
			L := [10]int{}
			L[p2l-1] = universes
			L = roll(L)
			L = roll(L)
			L = roll(L)
			for i, c := range L {
				if c > 0 {
					s := p2s + i + 1
					if s >= 21 {
						p2w += c
					} else {
						a, b := play(p1l, p1s, i+1, s, c, 0)
						p1w += a
						p2w += b
					}
				}
			}
		}
		return
	}

	p1w, p2w := play(p1, 0, p2, 0, 1, 0)

	winner := p1w
	if p2w > winner {
		winner = p2w
	}

	return winner, nil
}

func AnswerV3(input string) (int, error) {
	lines := strings.Split(input, "\n")
	p1, _ := strconv.Atoi(lines[0][28:])
	p2, _ := strconv.Atoi(lines[1][28:])

	// GameState is a map of [p1score, p2score] -> [p1positions, p2positions], where each position array is the number of universes in that spot
	type Board [2][10]int
	type GameState map[[2]int]Board
	state := GameState{}
	wins := [2]int{} // number of wins for each player

	// Initialize the first round of state with the starting positions
	b := Board{}
	b[0][p1-1] = 1
	b[1][p2-1] = 1
	state[[2]int{0, 0}] = b

	debug := func(s GameState) {
		type value struct {
			K  [2]int
			P1 [10]int
			P2 [10]int
		}

		keys, p1u, p2u := 0, 0, 0
		x := make([]value, 0, len(s))

		for k, v := range s {
			x = append(x, value{K: k, P1: v[0], P2: v[1]})
			keys++
			for _, vv := range v[0] {
				p1u += vv
			}
			for _, vv := range v[1] {
				p2u += vv
			}
		}

		sort.Slice(x, func(a, b int) bool {
			return x[a].K[0] < x[b].K[0] || (x[a].K[0] == x[b].K[0] && x[a].K[1] < x[b].K[1])
		})

		for _, v := range x {
			fmt.Println(v.K, v.P1, v.P2)
		}
		fmt.Println(keys, p1u, p2u)
		fmt.Println()
	}
	_ = debug

	active, inactive := 0, 1
	for step := 0; step < 50000 && len(state) > 0; step++ {
		// debug(state)
		next := GameState{}

		for scores, positions := range state {
			base := sum(positions[active][:])
			positions[active] = roll(positions[active])
			positions[active] = roll(positions[active])
			positions[active] = roll(positions[active])

			for pos, count := range positions[active] {
				if count == 0 {
					continue
				}

				newScore := scores[active] + pos + 1
				if newScore >= 21 {
					wins[active] += count
					continue
				}

				k := [2]int{}
				k[active] = newScore
				k[inactive] = scores[inactive]

				b := next[k]
				b[active][pos] += count
				for iPos, iCount := range positions[inactive] {
					b[inactive][iPos] += iCount * count / base
				}
				next[k] = b
			}
		}

		state = next
		active = (active + 1) % 2
		inactive = (inactive + 1) % 2
	}

	fmt.Println(wins)
	return max(wins[0], wins[1]), nil
}
