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

var Answer = AnswerV5

type board struct {
	nums      map[uint8]uint8
	rows      [5]uint8
	cols      [5]uint8
	completed bool
}

func AnswerPart1(input string) (int, error) {
	answers := make([]uint8, 0, 100)
	boards := make([]*board, 0, 100)

	i, l := 0, len(input)
	for v := uint8(0); i < l; i++ {
		c := input[i]
		if c == ',' {
			answers = append(answers, v)
			v = 0
		} else if c == '\n' {
			answers = append(answers, v)
			break
		} else {
			v = (v * 10) + uint8(c-'0')
		}
	}
	for j := 2; i+j < l; j += 76 {
		b := &board{nums: make(map[uint8]uint8, 25)}
		for k := 0; k < 25; k++ {
			o := i + j + 3*k
			v := uint8((k/5)<<4 | (k % 5))
			if input[o] == ' ' {
				b.nums[uint8(input[o+1]-'0')] = v
			} else {
				b.nums[10*uint8(input[o]-'0')+uint8(input[o+1]-'0')] = v
			}
		}
		boards = append(boards, b)
	}

	for _, a := range answers {
		for _, b := range boards {
			if v, ok := b.nums[a]; ok {
				r, c := v>>4, v&0b1111
				b.rows[r] |= 1 << c
				b.cols[c] |= 1 << r
				delete(b.nums, a)
				if b.rows[r] == 31 || b.cols[c] == 31 {
					uncalled := 0
					for k := range b.nums {
						uncalled += int(k)
					}
					return uncalled * int(a), nil
				}
			}
		}
	}

	return 0, fmt.Errorf("no winning board found")
}

func AnswerV1(input string) (int, error) {
	answers := make([]uint8, 0, 100)
	boards := make([]*board, 0, 100)

	i, l := 0, len(input)
	for v := uint8(0); i < l; i++ {
		c := input[i]
		if c == ',' {
			answers = append(answers, v)
			v = 0
		} else if c == '\n' {
			answers = append(answers, v)
			break
		} else {
			v = (v * 10) + uint8(c-'0')
		}
	}
	for j := 2; i+j < l; j += 76 {
		b := &board{nums: make(map[uint8]uint8, 25)}
		for k := 0; k < 25; k++ {
			o := i + j + 3*k
			v := uint8((k/5)<<4 | (k % 5))
			if input[o] == ' ' {
				b.nums[uint8(input[o+1]-'0')] = v
			} else {
				b.nums[10*uint8(input[o]-'0')+uint8(input[o+1]-'0')] = v
			}
		}
		boards = append(boards, b)
	}

	completed := 0
	for _, a := range answers {
		for _, b := range boards {
			if v, ok := b.nums[a]; ok {
				r, c := v>>4, v&0b1111
				b.rows[r] |= 1 << c
				b.cols[c] |= 1 << r
				delete(b.nums, a)
				if !b.completed && (b.rows[r] == 31 || b.cols[c] == 31) {
					b.completed = true
					completed++
					if completed == len(boards) {
						uncalled := 0
						for k := range b.nums {
							uncalled += int(k)
						}
						return uncalled * int(a), nil
					}
				}
			}
		}
	}

	return 0, fmt.Errorf("no winning board found")
}

type boardV2 struct {
	board  [5][5]uint8
	turns  [5][5]uint8
	winsOn uint8
}

func AnswerV2(input string) (int, error) {
	answers := make([]uint8, 100)
	boards := make([]boardV2, 0, 100)

	i, l := 0, len(input)
	for j, v := uint8(1), uint8(0); i < l; i++ {
		c := input[i]
		if c == ',' {
			answers[v] = j
			j++
			v = 0
		} else if c == '\n' {
			answers[v] = j
			break
		} else {
			v = (v * 10) + uint8(c-'0')
		}
	}
	for j := 2; i+j < l; j += 76 {
		b := boardV2{winsOn: 255}
		// Build the board
		for r := 0; r < 5; r++ {
			for c := 0; c < 5; c++ {
				o := i + j + 15*r + 3*c
				v := 10*(uint8(input[o]-'0')&0x0F) + uint8(input[o+1]-'0')
				b.board[r][c] = v
				b.turns[r][c] = answers[v]
			}
		}
		// Check rows to find when we win
		for r := 0; r < 5; r++ {
			t := uint8(0)
			for c := 0; c < 5; c++ {
				if b.turns[r][c] > t {
					t = b.turns[r][c]
				}
			}
			if t < b.winsOn {
				b.winsOn = t
			}
		}
		// Check columns to find when we win
		for c := 0; c < 5; c++ {
			t := uint8(0)
			for r := 0; r < 5; r++ {
				if b.turns[r][c] > t {
					t = b.turns[r][c]
				}
			}
			if t < b.winsOn {
				b.winsOn = t
			}
		}
		boards = append(boards, b)
	}

	idx, winsOn := 0, uint8(0)
	for i, b := range boards {
		if b.winsOn > winsOn {
			idx, winsOn = i, b.winsOn
		}
	}
	b := boards[idx]

	winningMove, remaining := 0, 0
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if b.turns[r][c] > b.winsOn {
				remaining += int(b.board[r][c])
			} else if b.turns[r][c] == b.winsOn {
				winningMove = int(b.board[r][c])
			}
		}
	}
	return winningMove * remaining, nil
}

func AnswerV3(input string) (int, error) {
	answers := make([]uint8, 100)
	boards := make([]boardV2, 0, 100)

	i, l := 0, len(input)
	for j, v := uint8(1), uint8(0); i < l; i++ {
		c := input[i]
		if c == ',' {
			answers[v] = j
			j++
			v = 0
		} else if c == '\n' {
			answers[v] = j
			break
		} else {
			v = (v * 10) + uint8(c-'0')
		}
	}
	for i += 2; i < l; i++ {
		b := boardV2{winsOn: 255}
		// Build the board
		for r := 0; r < 5; r++ {
			for c := 0; c < 5; c++ {
				v := 10*(uint8(input[i]-'0')&0x0F) + uint8(input[i+1]-'0')
				b.board[r][c] = v
				b.turns[r][c] = answers[v]
				i += 3
			}
		}
		// Check rows to find when we win
		for r := 0; r < 5; r++ {
			t := uint8(0)
			for c := 0; c < 5; c++ {
				if b.turns[r][c] > t {
					t = b.turns[r][c]
				}
			}
			if t < b.winsOn {
				b.winsOn = t
			}
		}
		// Check columns to find when we win
		for c := 0; c < 5; c++ {
			t := uint8(0)
			for r := 0; r < 5; r++ {
				if b.turns[r][c] > t {
					t = b.turns[r][c]
				}
			}
			if t < b.winsOn {
				b.winsOn = t
			}
		}
		boards = append(boards, b)
	}

	idx, winsOn := 0, uint8(0)
	for i, b := range boards {
		if b.winsOn > winsOn {
			idx, winsOn = i, b.winsOn
		}
	}
	b := boards[idx]

	winningMove, remaining := 0, 0
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if b.turns[r][c] > b.winsOn {
				remaining += int(b.board[r][c])
			} else if b.turns[r][c] == b.winsOn {
				winningMove = int(b.board[r][c])
			}
		}
	}
	return winningMove * remaining, nil
}

type boardV3 struct {
	board  [5][5]int8
	turns  [5][5]int8
	winsOn int8
}

func AnswerV4(input string) (int, error) {
	answers := make([]int8, 100)
	boards := make([]boardV3, 0, 100)

	i, l := 0, len(input)
	for j, v := int8(1), int8(0); i < l; i++ {
		c := input[i]
		if c == ',' {
			answers[v] = j
			j++
			v = 0
		} else if c == '\n' {
			answers[v] = j
			break
		} else {
			v = (v * 10) + int8(c-'0')
		}
	}
	for i += 2; i < l; i++ {
		b := boardV3{winsOn: 101}
		// Build the board
		for r := 0; r < 5; r++ {
			for c := 0; c < 5; c++ {
				v := 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
				b.board[r][c] = v
				b.turns[r][c] = answers[v]
				i += 3
			}
		}
		// Check rows/cols to find when we win
		for r := 0; r < 5; r++ {
			tr, tc := int8(0), int8(0)
			for c := 0; c < 5; c++ {
				tr = max(tr, b.turns[r][c])
				tc = max(tc, b.turns[c][r])
			}
			b.winsOn = min(b.winsOn, tr)
			b.winsOn = min(b.winsOn, tc)
		}
		boards = append(boards, b)
	}

	idx, winsOn := 0, int8(0)
	for i, b := range boards {
		if b.winsOn > winsOn {
			idx, winsOn = i, b.winsOn
		}
	}
	b := boards[idx]

	winningMove, remaining := 0, 0
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if b.turns[r][c] > b.winsOn {
				remaining += int(b.board[r][c])
			} else if b.turns[r][c] == b.winsOn {
				winningMove = int(b.board[r][c])
			}
		}
	}
	return winningMove * remaining, nil
}

type boardV4 struct {
	board  [25]int8
	turns  [25]int8
	winsOn int8
}

func AnswerV5(input string) (int, error) {
	answers := make([]int8, 100)
	boards := make([]boardV4, 0, 100)

	i, l := 0, len(input)
	for j, v := int8(1), int8(0); i < l; i++ {
		c := input[i]
		if c == ',' {
			answers[v] = j
			j++
			v = 0
		} else if c == '\n' {
			answers[v] = j
			break
		} else {
			v = (v * 10) + int8(c-'0')
		}
	}
	for i += 2; i < l; i++ {
		b := boardV4{winsOn: 101}
		// Build the board (unrolled)
		v := 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[0] = v
		b.turns[0] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[1] = v
		b.turns[1] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[2] = v
		b.turns[2] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[3] = v
		b.turns[3] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[4] = v
		b.turns[4] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[5] = v
		b.turns[5] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[6] = v
		b.turns[6] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[7] = v
		b.turns[7] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[8] = v
		b.turns[8] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[9] = v
		b.turns[9] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[10] = v
		b.turns[10] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[11] = v
		b.turns[11] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[12] = v
		b.turns[12] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[13] = v
		b.turns[13] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[14] = v
		b.turns[14] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[15] = v
		b.turns[15] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[16] = v
		b.turns[16] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[17] = v
		b.turns[17] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[18] = v
		b.turns[18] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[19] = v
		b.turns[19] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[20] = v
		b.turns[20] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[21] = v
		b.turns[21] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[22] = v
		b.turns[22] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[23] = v
		b.turns[23] = answers[v]
		i += 3
		v = 10*(int8(input[i]-'0')&0x0F) + int8(input[i+1]-'0')
		b.board[24] = v
		b.turns[24] = answers[v]
		i += 3
		// Check rows to find when we win (unrolled)
		t := int8(0)
		t = max(t, b.turns[0])
		t = max(t, b.turns[1])
		t = max(t, b.turns[2])
		t = max(t, b.turns[3])
		t = max(t, b.turns[4])
		b.winsOn = min(b.winsOn, t)
		t = 0
		t = max(t, b.turns[5])
		t = max(t, b.turns[6])
		t = max(t, b.turns[7])
		t = max(t, b.turns[8])
		t = max(t, b.turns[9])
		b.winsOn = min(b.winsOn, t)
		t = 0
		t = max(t, b.turns[10])
		t = max(t, b.turns[11])
		t = max(t, b.turns[12])
		t = max(t, b.turns[13])
		t = max(t, b.turns[14])
		b.winsOn = min(b.winsOn, t)
		t = 0
		t = max(t, b.turns[15])
		t = max(t, b.turns[16])
		t = max(t, b.turns[17])
		t = max(t, b.turns[18])
		t = max(t, b.turns[19])
		b.winsOn = min(b.winsOn, t)
		t = 0
		t = max(t, b.turns[20])
		t = max(t, b.turns[21])
		t = max(t, b.turns[22])
		t = max(t, b.turns[23])
		t = max(t, b.turns[24])
		b.winsOn = min(b.winsOn, t)
		// Check cols to find when we win (unrolled)
		t = 0
		t = max(t, b.turns[0])
		t = max(t, b.turns[5])
		t = max(t, b.turns[10])
		t = max(t, b.turns[15])
		t = max(t, b.turns[20])
		b.winsOn = min(b.winsOn, t)
		t = 0
		t = max(t, b.turns[1])
		t = max(t, b.turns[6])
		t = max(t, b.turns[11])
		t = max(t, b.turns[16])
		t = max(t, b.turns[21])
		b.winsOn = min(b.winsOn, t)
		t = 0
		t = max(t, b.turns[2])
		t = max(t, b.turns[7])
		t = max(t, b.turns[12])
		t = max(t, b.turns[17])
		t = max(t, b.turns[22])
		b.winsOn = min(b.winsOn, t)
		t = 0
		t = max(t, b.turns[3])
		t = max(t, b.turns[8])
		t = max(t, b.turns[13])
		t = max(t, b.turns[18])
		t = max(t, b.turns[23])
		b.winsOn = min(b.winsOn, t)
		t = 0
		t = max(t, b.turns[4])
		t = max(t, b.turns[9])
		t = max(t, b.turns[14])
		t = max(t, b.turns[19])
		t = max(t, b.turns[24])
		b.winsOn = min(b.winsOn, t)

		boards = append(boards, b)
	}

	idx, winsOn := 0, int8(0)
	for i, b := range boards {
		if b.winsOn > winsOn {
			idx, winsOn = i, b.winsOn
		}
	}
	b := boards[idx]

	winningMove, remaining := 0, 0
	for k := 0; k < 25; k++ {
		if b.turns[k] > b.winsOn {
			remaining += int(b.board[k])
		} else if b.turns[k] == b.winsOn {
			winningMove = int(b.board[k])
		}
	}
	return winningMove * remaining, nil
}

func min(x, y int8) int8 {
	return y + ((x - y) & ((x - y) >> 7))
}

func max(x, y int8) int8 {
	return x - ((x - y) & ((x - y) >> 7))
}
