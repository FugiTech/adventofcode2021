package main

import (
	"container/heap"
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

var Answer = AnswerV2

func AnswerPart1(input string) (int, error) {
	board := [][]int{}

	lines := strings.Split(input, "\n")
	for _, l := range lines {
		b := make([]int, len(l))
		for i, c := range []byte(l) {
			b[i] = int(c - '0')
		}
		board = append(board, b)
	}

	moves := make([]position, 0, 1000)
	moves = append(moves, position{X: 0, Y: 0, Risk: 0})
	visited := map[int]struct{}{}

	addMove := func(x, y, risk int) {
		if x < 0 || x >= len(board[0]) {
			return
		}
		if y < 0 || y >= len(board) {
			return
		}
		k := (y << 16) | x
		if _, ok := visited[k]; ok {
			return
		}
		visited[k] = struct{}{}
		moves = append(moves, position{
			X:    x,
			Y:    y,
			Risk: risk + board[y][x],
		})
	}

	for {
		m := moves[0]
		moves = moves[1:]
		if m.X == len(board[0])-1 && m.Y == len(board)-1 {
			return m.Risk, nil
		}
		addMove(m.X-1, m.Y, m.Risk)
		addMove(m.X+1, m.Y, m.Risk)
		addMove(m.X, m.Y-1, m.Risk)
		addMove(m.X, m.Y+1, m.Risk)
		sort.Slice(moves, func(i, j int) bool {
			return moves[i].Risk < moves[j].Risk
		})
	}
}

type position struct {
	X        int
	Y        int
	Risk     int
	Priority int
	index    int
}

func AnswerV1(input string) (int, error) {
	lines := strings.Split(input, "\n")
	board := make([][]int, len(lines)*5)

	for y, l := range lines {
		for r := 0; r < 5; r++ {
			b := make([]int, len(l)*5)
			for x, c := range []byte(l) {
				v := int(c - '0')
				for rr := 0; rr < 5; rr++ {
					b[x+len(l)*rr] = ((v + r + rr - 1) % 9) + 1
				}
			}
			board[y+len(lines)*r] = b
		}
	}

	moves := make([]position, 0, 1000)
	moves = append(moves, position{X: 0, Y: 0, Risk: 0})
	visited := map[int]struct{}{}

	addMove := func(x, y, risk int) {
		if x < 0 || x >= len(board[0]) {
			return
		}
		if y < 0 || y >= len(board) {
			return
		}
		k := (y << 16) | x
		if _, ok := visited[k]; ok {
			return
		}
		visited[k] = struct{}{}
		moves = append(moves, position{
			X:    x,
			Y:    y,
			Risk: risk + board[y][x],
		})
	}

	for {
		m := moves[0]
		moves = moves[1:]
		if m.X == len(board[0])-1 && m.Y == len(board)-1 {
			return m.Risk, nil
		}
		addMove(m.X-1, m.Y, m.Risk)
		addMove(m.X+1, m.Y, m.Risk)
		addMove(m.X, m.Y-1, m.Risk)
		addMove(m.X, m.Y+1, m.Risk)
		sort.Slice(moves, func(i, j int) bool {
			return moves[i].Risk < moves[j].Risk
		})
	}
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*position

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// min-priority queue actually uses less properly
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*position)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func AnswerV2(input string) (int, error) {
	L := strings.IndexByte(input, '\n')
	board := [500][500]uint8{}

	for y := 0; y < L; y++ {
		for x := 0; x < L; x++ {
			v := uint8(input[y*(L+1)+x] - '0')
			v1 := (v % 9) + 1
			v2 := ((v + 1) % 9) + 1
			v3 := ((v + 2) % 9) + 1
			v4 := ((v + 3) % 9) + 1
			v5 := ((v + 4) % 9) + 1
			v6 := ((v + 5) % 9) + 1
			v7 := ((v + 6) % 9) + 1
			v8 := ((v + 7) % 9) + 1

			board[y+(0*L)][x+(0*L)] = v
			board[y+(0*L)][x+(1*L)] = v1
			board[y+(1*L)][x+(0*L)] = v1
			board[y+(0*L)][x+(2*L)] = v2
			board[y+(1*L)][x+(1*L)] = v2
			board[y+(2*L)][x+(0*L)] = v2
			board[y+(0*L)][x+(3*L)] = v3
			board[y+(1*L)][x+(2*L)] = v3
			board[y+(2*L)][x+(1*L)] = v3
			board[y+(3*L)][x+(0*L)] = v3
			board[y+(0*L)][x+(4*L)] = v4
			board[y+(1*L)][x+(3*L)] = v4
			board[y+(2*L)][x+(2*L)] = v4
			board[y+(3*L)][x+(1*L)] = v4
			board[y+(4*L)][x+(0*L)] = v4
			board[y+(1*L)][x+(4*L)] = v5
			board[y+(2*L)][x+(3*L)] = v5
			board[y+(3*L)][x+(2*L)] = v5
			board[y+(4*L)][x+(1*L)] = v5
			board[y+(2*L)][x+(4*L)] = v6
			board[y+(3*L)][x+(3*L)] = v6
			board[y+(4*L)][x+(2*L)] = v6
			board[y+(3*L)][x+(4*L)] = v7
			board[y+(4*L)][x+(3*L)] = v7
			board[y+(4*L)][x+(4*L)] = v8
		}
	}

	moves := make(PriorityQueue, 0, 1000)
	moves = append(moves, &position{X: 0, Y: 0, Risk: 0})
	heap.Init(&moves)
	visited := map[[2]int]struct{}{}

	addMove := func(x, y, risk int) {
		if x < 0 || x >= 5*L {
			return
		}
		if y < 0 || y >= 5*L {
			return
		}
		k := [2]int{y, x}
		if _, ok := visited[k]; ok {
			return
		}
		visited[k] = struct{}{}
		heap.Push(&moves, &position{
			X:        x,
			Y:        y,
			Risk:     risk + int(board[y][x]),
			Priority: risk + int(board[y][x]) + 10*L - 2 - x - y,
		})
	}

	for {
		m := heap.Pop(&moves).(*position)
		if m.X == 5*L-1 && m.Y == 5*L-1 {
			return m.Risk, nil
		}
		addMove(m.X-1, m.Y, m.Risk)
		addMove(m.X+1, m.Y, m.Risk)
		addMove(m.X, m.Y-1, m.Risk)
		addMove(m.X, m.Y+1, m.Risk)
	}
}
