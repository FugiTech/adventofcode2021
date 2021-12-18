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

type BitReader struct {
	Input       string
	index       int
	buffer      uint64
	bitsWaiting int
}

func (r *BitReader) Remaining() bool {
	return r.index < len(r.Input) || r.bitsWaiting > 0
}

func (r *BitReader) Read(bits int) uint64 {
	for r.bitsWaiting < bits && r.index <= len(r.Input)-2 {
		v, _ := strconv.ParseUint(r.Input[r.index:r.index+2], 16, 64)
		r.index += 2
		r.buffer = (r.buffer << 8) | v
		r.bitsWaiting += 8
	}
	r.bitsWaiting -= bits

	v := r.buffer >> r.bitsWaiting
	r.buffer = r.buffer & ((1 << r.bitsWaiting) - 1)
	return v
}

func (r *BitReader) ClearBuffer() {
	r.bitsWaiting = 0
}

func (r *BitReader) ReadPacket() (V, T, val, bits uint64) {
	V = r.Read(3)
	bits += 3

	T = r.Read(3)
	bits += 3

	if T == 4 { // Literal
		for {
			tmp := r.Read(5)
			bits += 5

			val = (val << 4) | (tmp & 0b1111)
			if tmp>>4 == 0 {
				break
			}
		}
	} else { // Operator
		LT := r.Read(1)
		bits += 1

		var process func(uint64)
		subvals := make([]uint64, 0, 64)
		switch T {
		case 0:
			process = func(v uint64) { val += v }
		case 1:
			val = 1
			process = func(v uint64) { val *= v }
		case 2:
			val = math.MaxUint64
			process = func(v uint64) {
				if v < val {
					val = v
				}
			}
		case 3:
			process = func(v uint64) {
				if v > val {
					val = v
				}
			}
		case 5:
			process = func(_ uint64) {
				if len(subvals) == 2 && subvals[0] > subvals[1] {
					val = 1
				} else {
					val = 0
				}
			}
		case 6:
			process = func(_ uint64) {
				if len(subvals) == 2 && subvals[0] < subvals[1] {
					val = 1
				} else {
					val = 0
				}
			}
		case 7:
			process = func(_ uint64) {
				if len(subvals) == 2 && subvals[0] == subvals[1] {
					val = 1
				} else {
					val = 0
				}
			}
		}

		if LT == 0 {
			L := r.Read(15)
			bits += 15

			for i := uint64(0); i < L; {
				v, _, subval, b := r.ReadPacket()
				V += v
				bits += b
				i += b
				subvals = append(subvals, subval)
				process(subval)
			}

		} else {
			L := r.Read(11)
			bits += 11

			for i := uint64(0); i < L; i++ {
				v, _, subval, b := r.ReadPacket()
				V += v
				bits += b
				subvals = append(subvals, subval)
				process(subval)
			}
		}
	}

	return V, T, val, bits
}

func AnswerPart1(input string) (int, error) {
	versionSum := 0

	r := &BitReader{Input: input}

	for r.Remaining() {
		V, _, _, _ := r.ReadPacket()
		versionSum += int(V)
		r.ClearBuffer()
	}

	return versionSum, nil
}

func AnswerV1(input string) (int, error) {
	r := &BitReader{Input: input}
	_, _, v, _ := r.ReadPacket()

	return int(v), nil
}
