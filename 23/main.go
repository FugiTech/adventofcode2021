package main

import (
	"container/heap"
	"fmt"
	"io/ioutil"
	"math"
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

var Answer = AnswerV1

type state struct {
	hallway [11]int
	rooms   [4][4]int
	energy  int
	cost    int

	index  int
	parent *state
}

func (s *state) clone() *state {
	return &state{
		hallway: s.hallway,
		rooms:   s.rooms,
		energy:  s.energy,
		parent:  s,
	}
}

func AnswerPart1(input string) (int, error) {
	energy := [4]int{1, 10, 100, 1000}

	debug := func(s *state) {
		m := map[int]string{
			0: " ",
			1: "A",
			2: "B",
			3: "C",
			4: "D",
		}
		fmt.Println("#############")
		fmt.Print("#")
		for _, c := range s.hallway {
			fmt.Print(m[c])
		}
		fmt.Println("#")
		fmt.Printf("###%s#%s#%s#%s### energy: %d\n", m[s.rooms[0][0]], m[s.rooms[1][0]], m[s.rooms[2][0]], m[s.rooms[3][0]], s.energy)
		fmt.Printf("  #%s#%s#%s#%s#  \n", m[s.rooms[0][1]], m[s.rooms[1][1]], m[s.rooms[2][1]], m[s.rooms[3][1]])
		fmt.Println("  #########  ")
	}
	var recursiveDebug func(*state)
	recursiveDebug = func(s *state) {
		if s.parent != nil {
			recursiveDebug(s.parent)
		}
		debug(s)
	}

	cost := func(s *state) int {
		// Determine the minimum energy required to teleport each amphipod directly to where it needs to go, ignoring collisions
		c := 0
		for x, v := range s.hallway {
			if v == 0 {
				continue
			}
			tx := 2 * v
			c += int(math.Abs(float64(tx-x))) + 1
		}
		for rnum, room := range s.rooms {
			for _, v := range room {
				if rnum+1 == v {
					continue
				}
				x, tx := 2+2*rnum, 2*v
				c += int(math.Abs(float64(tx-x))) + 2
			}
		}
		return c
	}

	actuallyTrulyValid := func(s *state) bool {
		minV := 0
		for _, v := range s.hallway {
			if v == 0 {
				continue
			}
			if v < minV {
				return false
			}
			minV = v
		}
		return true
	}

	validMoves := func(s *state) []*state {
		r := make([]*state, 0, 100)

		// Check if we can move any hallway pods into their room
		for ix, v := range s.hallway {
			// Empty cells can't move
			if v == 0 {
				continue
			}

			// Room needs to be empty or partially filled with the correct amphipod
			if s.rooms[v-1][0] != 0 || (s.rooms[v-1][1] != 0 && s.rooms[v-1][1] != v) {
				continue
			}

			clear := true
			sx, tx := ix, 2*v
			if tx < sx {
				sx, tx = tx, sx
			}
			for x := sx + 1; x < tx; x++ {
				clear = clear && s.hallway[x] == 0
			}
			if !clear {
				continue
			}

			ns := s.clone()
			ns.hallway[ix] = 0

			moves := tx - sx
			if s.rooms[v-1][1] == 0 {
				ns.rooms[v-1][1] = v
				moves += 2
			} else {
				ns.rooms[v-1][0] = v
				moves++
			}

			ns.energy += moves * energy[v-1]
			ns.cost = ns.energy + cost(ns)
			r = append(r, ns)
		}

		// Check if any pods in any rooms can be moved into the hallway
		for rnum, room := range s.rooms {
			if room[1] == 0 || (room[1] == rnum+1 && (room[0] == 0 || room[0] == rnum+1)) {
				continue
			}

			sx := 2 + 2*rnum
			validSpots := make([]int, 0, 7)
			for x := sx - 1; x >= 0 && s.hallway[x] == 0; x-- {
				if x != 2 && x != 4 && x != 6 && x != 8 {
					validSpots = append(validSpots, x)
				}
			}
			for x := sx + 1; x < 11 && s.hallway[x] == 0; x++ {
				if x != 2 && x != 4 && x != 6 && x != 8 {
					validSpots = append(validSpots, x)
				}
			}

			for _, tx := range validSpots {
				ns := s.clone()

				moves := int(math.Abs(float64(tx - sx)))

				if room[0] == 0 {
					ns.hallway[tx] = room[1]
					ns.rooms[rnum][1] = 0
					moves += 2
				} else {
					ns.hallway[tx] = room[0]
					ns.rooms[rnum][0] = 0
					moves += 1
				}

				if actuallyTrulyValid(ns) {
					ns.energy += moves * energy[ns.hallway[tx]-1]
					ns.cost = ns.energy + cost(ns)
					r = append(r, ns)
				}
			}
		}

		return r
	}

	complete := func(s *state) bool {
		return (s.rooms[0][0] == 1 && s.rooms[0][1] == 1 &&
			s.rooms[1][0] == 2 && s.rooms[1][1] == 2 &&
			s.rooms[2][0] == 3 && s.rooms[2][1] == 3 &&
			s.rooms[3][0] == 4 && s.rooms[3][1] == 4)
	}

	_, _ = debug, recursiveDebug
	initial := &state{}
	for y := 0; y < 2; y++ {
		for x := 0; x < 4; x++ {
			initial.rooms[x][y] = int(input[28+14*y+3+2*x] - 'A' + 1)
		}
	}

	states := make(PriorityQueue, 0, 1000)
	states = append(states, initial)
	heap.Init(&states)

	for step := 1; len(states) > 0; step++ {
		s := heap.Pop(&states).(*state)
		if step == 10000000 {
			recursiveDebug(s)
			return 0, nil
		}
		for _, v := range validMoves(s) {
			if complete(v) {
				recursiveDebug(v)
				return v.energy, nil
			}
			heap.Push(&states, v)
		}
		if step%10000 == 0 {
			fmt.Println(step, len(states))
		}
	}
	return 0, nil
}

func AnswerV1(input string) (int, error) {
	energy := [4]int{1, 10, 100, 1000}

	debug := func(s *state) {
		m := map[int]string{
			0: " ",
			1: "A",
			2: "B",
			3: "C",
			4: "D",
		}
		fmt.Println("#############")
		fmt.Print("#")
		for _, c := range s.hallway {
			fmt.Print(m[c])
		}
		fmt.Println("#")
		fmt.Printf("###%s#%s#%s#%s### energy: %d\n", m[s.rooms[0][0]], m[s.rooms[1][0]], m[s.rooms[2][0]], m[s.rooms[3][0]], s.energy)
		fmt.Printf("  #%s#%s#%s#%s#  \n", m[s.rooms[0][1]], m[s.rooms[1][1]], m[s.rooms[2][1]], m[s.rooms[3][1]])
		fmt.Printf("  #%s#%s#%s#%s#  \n", m[s.rooms[0][2]], m[s.rooms[1][2]], m[s.rooms[2][2]], m[s.rooms[3][2]])
		fmt.Printf("  #%s#%s#%s#%s#  \n", m[s.rooms[0][3]], m[s.rooms[1][3]], m[s.rooms[2][3]], m[s.rooms[3][3]])
		fmt.Println("  #########  ")
	}
	var recursiveDebug func(*state)
	recursiveDebug = func(s *state) {
		if s.parent != nil {
			recursiveDebug(s.parent)
		}
		debug(s)
	}

	cost := func(s *state) int {
		// Determine the minimum energy required to teleport each amphipod directly to where it needs to go, ignoring collisions
		c := 0
		for x, v := range s.hallway {
			if v == 0 {
				continue
			}
			tx := 2 * v
			c += (int(math.Abs(float64(tx-x))) + 1) * energy[v-1]
		}
		for rnum, room := range s.rooms {
			for _, v := range room {
				if rnum+1 == v || v == 0 {
					continue
				}
				x, tx := 2+2*rnum, 2*v
				c += (int(math.Abs(float64(tx-x))) + 2) * energy[v-1]
			}
		}
		return c
	}

	hrtMemo := map[[11]int]bool{}
	hrtMemo[[11]int{}] = true
	var hallwayReachableTest func([11]int) bool
	hallwayReachableTest = func(hallway [11]int) bool {
		if v, ok := hrtMemo[hallway]; ok {
			return v
		}

		for ix, v := range hallway {
			if v == 0 {
				continue
			}

			tx := 2 * v
			reachable := true
			if tx < ix {
				for x := ix - 1; x >= tx; x-- {
					if hallway[x] != 0 {
						reachable = false
					}
				}
			} else {
				for x := ix + 1; x <= tx; x++ {
					if hallway[x] != 0 {
						reachable = false
					}
				}
			}

			if reachable {
				nh := hallway
				nh[ix] = 0
				if hallwayReachableTest(nh) {
					hrtMemo[hallway] = true
					return true
				}
			}
		}

		hrtMemo[hallway] = false
		return false
	}

	enterable := func(room []int, v int) bool {
		for _, c := range room {
			if c != 0 && c != v {
				return false
			}
		}
		return true
	}
	leavable := func(room []int, v int) bool {
		for _, c := range room {
			if c != 0 && c != v {
				return true
			}
		}
		return false
	}

	validMoves := func(s *state) []*state {
		r := make([]*state, 0, 100)

		// Check if we can move any hallway pods into their room
		for ix, v := range s.hallway {
			// Empty cells can't move
			if v == 0 {
				continue
			}

			// Room needs to be empty or partially filled with the correct amphipod
			if !enterable(s.rooms[v-1][:], v) {
				continue
			}

			clear := true
			sx, tx := ix, 2*v
			if tx < sx {
				sx, tx = tx, sx
			}
			for x := sx + 1; x < tx; x++ {
				clear = clear && s.hallway[x] == 0
			}
			if !clear {
				continue
			}

			ns := s.clone()
			ns.hallway[ix] = 0

			moves := tx - sx
			for d := 3; d >= 0; d-- {
				if ns.rooms[v-1][d] != 0 {
					continue
				}
				ns.rooms[v-1][d] = v
				moves += d + 1
				break
			}

			ns.energy += moves * energy[v-1]
			ns.cost = ns.energy + cost(ns)
			r = append(r, ns)
		}

		// Check if any pods in any rooms can be moved into the hallway
		for rnum, room := range s.rooms {
			if !leavable(room[:], rnum+1) {
				continue
			}

			sx := 2 + 2*rnum
			validSpots := make([]int, 0, 7)
			for x := sx - 1; x >= 0 && s.hallway[x] == 0; x-- {
				if x != 2 && x != 4 && x != 6 && x != 8 {
					validSpots = append(validSpots, x)
				}
			}
			for x := sx + 1; x < 11 && s.hallway[x] == 0; x++ {
				if x != 2 && x != 4 && x != 6 && x != 8 {
					validSpots = append(validSpots, x)
				}
			}

			for _, tx := range validSpots {
				ns := s.clone()

				moves := int(math.Abs(float64(tx - sx)))

				for d := 0; d < 4; d++ {
					if room[d] == 0 {
						continue
					}
					ns.hallway[tx] = room[d]
					ns.rooms[rnum][d] = 0
					moves += d + 1
					break
				}

				if hallwayReachableTest(ns.hallway) {
					ns.energy += moves * energy[ns.hallway[tx]-1]
					ns.cost = ns.energy + cost(ns)
					r = append(r, ns)
				}
			}
		}

		return r
	}

	complete := func(s *state) bool {
		return (s.rooms[0][0] == 1 && s.rooms[0][1] == 1 && s.rooms[0][2] == 1 && s.rooms[0][3] == 1 &&
			s.rooms[1][0] == 2 && s.rooms[1][1] == 2 && s.rooms[1][2] == 2 && s.rooms[1][3] == 2 &&
			s.rooms[2][0] == 3 && s.rooms[2][1] == 3 && s.rooms[2][2] == 3 && s.rooms[2][3] == 3 &&
			s.rooms[3][0] == 4 && s.rooms[3][1] == 4 && s.rooms[3][2] == 4 && s.rooms[3][3] == 4)
	}

	initial := &state{}
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			initial.rooms[x][y] = int(input[28+14*y+3+2*x] - 'A' + 1)
		}
	}

	states := make(PriorityQueue, 0, 1000)
	states = append(states, initial)
	heap.Init(&states)

	for step := 1; len(states) > 0; step++ {
		s := heap.Pop(&states).(*state)
		for _, v := range validMoves(s) {
			if complete(v) {
				recursiveDebug(v)
				return v.energy, nil
			}
			heap.Push(&states, v)
		}
		if step%100000 == 0 {
			fmt.Println(step, len(states), s.energy, s.cost)
		}
	}
	return 0, nil
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*state

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// min-priority queue actually uses less properly
	return pq[i].cost < pq[j].cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*state)
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
