package main

import (
	"encoding/json"
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
	var N *SnailfishNumber

	lines := strings.Split(input, "\n")
	for _, l := range lines {
		v := ParseSnailfishNumber(l)
		N = AddSnailfishNumber(N, v)
	}

	return N.Magnitude(), nil
}

func AnswerV1(input string) (int, error) {

	return 0, nil
}

func AddSnailfishNumber(l, r *SnailfishNumber) *SnailfishNumber {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}

	v := &SnailfishNumber{L: l, R: r}
	v.Reduce()
	return v
}

func ParseSnailfishNumber(input string) *SnailfishNumber {
	var data interface{}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		panic(err)
	}

	var interpretInterface func(d interface{}) *SnailfishNumber
	interpretInterface = func(d interface{}) *SnailfishNumber {
		switch v := d.(type) {
		case float64:
			return &SnailfishNumber{V: int(v)}
		case []interface{}:
			return &SnailfishNumber{
				L: interpretInterface(v[0]),
				R: interpretInterface(v[1]),
			}
		default:
			panic(fmt.Errorf("unknown type: %t %v", d, d))
		}
	}

	return interpretInterface(data)
}

type SnailfishNumber struct {
	L *SnailfishNumber
	R *SnailfishNumber
	V int
}

func (n *SnailfishNumber) Reduce() {
	for run := true; run; {
		// fmt.Println(n)
		_, _, run = n.explode(nil, 0)
		run = run || n.split()
	}
}

func (n *SnailfishNumber) explode(p *SnailfishNumber, depth int) (L int, R int, actionTaken bool) {
	if n == nil {
		return 0, 0, false
	}

	// Explode
	if depth >= 4 && n.L != nil && n.R != nil {
		L = n.L.V
		R = n.R.V
		n.L = nil
		n.R = nil
		n.V = 0
		return L, R, true
	}

	// Handle childs exploding
	if l, r, at := n.L.explode(n, depth+1); at {
		if n.R.addLeft(r) {
			r = 0
		}
		return l, r, true
	}
	if l, r, at := n.R.explode(n, depth+1); at {
		if n.L.addRight(l) {
			l = 0
		}
		return l, r, true
	}
	return 0, 0, false
}

func (n *SnailfishNumber) addLeft(v int) (success bool) {
	if n == nil {
		return false
	}
	if n.L == nil && n.R == nil {
		n.V += v
		return true
	}
	return n.L.addLeft(v) || n.R.addLeft(v)
}

func (n *SnailfishNumber) addRight(v int) (success bool) {
	if n == nil {
		return false
	}
	if n.L == nil && n.R == nil {
		n.V += v
		return true
	}
	return n.R.addRight(v) || n.L.addRight(v)
}

func (n *SnailfishNumber) split() (actionTaken bool) {
	if n == nil {
		return false
	}

	if n.V > 9 {
		n.L = &SnailfishNumber{V: n.V / 2}
		n.R = &SnailfishNumber{V: n.V - n.L.V}
		n.V = 0
		return true
	}

	return n.L.split() || n.R.split()
}

func (n *SnailfishNumber) Magnitude() int {
	if n == nil {
		return 0
	}
	return 3*n.L.Magnitude() + 2*n.R.Magnitude() + n.V
}

func (n *SnailfishNumber) String() string {
	if n == nil {
		return ""
	}
	if n.L != nil || n.R != nil {
		return fmt.Sprintf("[%s,%s]", n.L, n.R)
	}
	return strconv.Itoa(n.V)
}
