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
	var counts [16]int
	lines := 1

	for i, j, l := 0, 0, len(input); i < l; i, j = i+1, j+1 {
		if input[i] == '1' {
			counts[j]++
		} else if input[i] == '\n' {
			j = -1
			lines++
		}
	}

	gamma, epsilon := 0, 0
	for _, v := range counts[:] {
		if v != 0 {
			if v > lines/2 { // mostly 1s
				gamma = gamma<<1 | 1
				epsilon = epsilon<<1 | 0
			} else { // mostly 0s
				gamma = gamma<<1 | 0
				epsilon = epsilon<<1 | 1
			}
		}
	}

	return gamma * epsilon, nil
}

func AnswerV1(input string) (int, error) {
	lines := strings.Split(input, "\n")
	bits := len(lines[0])

	numbers := make([]int, len(lines))
	for i, l := range lines {
		v, _ := strconv.ParseInt(l, 2, 64)
		numbers[i] = int(v)
	}

	filter := func(arr []int, bit int, invert bool) []int {
		count := 0
		for _, v := range arr {
			count += (v >> bit) & 1
		}

		search := 0
		if count >= (len(arr)+1)/2 {
			search = 1
		}
		if invert {
			search ^= 1
		}

		ret := make([]int, 0, len(arr))
		for _, v := range arr {
			if (v>>bit)&1 == search {
				ret = append(ret, v)
			}
		}

		// fmt.Printf("filter(bit=%d,invert=%t,search=%d) => len(%d)\n", bit, invert, search, len(ret))
		return ret
	}

	find := func(invert bool) int {
		arr := numbers
		for i := bits - 1; i >= 0; i-- {
			arr = filter(arr, i, invert)
			if len(arr) == 1 {
				return arr[0]
			}
		}
		return arr[0]
	}

	o2GenRating, co2ScrubRating := find(false), find(true)
	return o2GenRating * co2ScrubRating, nil
}

func AnswerV2(input string) (int, error) {
	bits := strings.IndexByte(input, '\n')
	numbers := make([]int, (len(input)+bits)/(bits+1))

	/*
		for i, l := 0, len(input); i < l; i += bits + 1 {
			v := 0
			for j := 0; j < bits; j++ {
				v |= int(input[i+j]-'0') << (bits - (j + 1))
			}
			numbers = append(numbers, v)
		}
	*/

	for i, j, l := 0, 0, len(input); i < l; i, j = i+bits+1, j+1 {
		v, _ := strconv.ParseUint(input[i:i+bits], 2, 64)
		numbers[j] = int(v)
	}

	filter := func(arr []int, bit int, invert bool) []int {
		count := 0
		for _, v := range arr {
			count += (v >> bit) & 1
		}

		search := 0
		if count >= (len(arr)+1)/2 {
			search = 1
		}
		if invert {
			search ^= 1
		}

		ret := make([]int, 0, len(arr))
		for _, v := range arr {
			if (v>>bit)&1 == search {
				ret = append(ret, v)
			}
		}

		// fmt.Printf("filter(bit=%d,invert=%t,search=%d) => len(%d)\n", bit, invert, search, len(ret))
		return ret
	}

	find := func(invert bool) int {
		arr := numbers
		for i := bits - 1; i >= 0; i-- {
			arr = filter(arr, i, invert)
			if len(arr) == 1 {
				return arr[0]
			}
		}
		return arr[0]
	}

	o2GenRating, co2ScrubRating := find(false), find(true)
	return o2GenRating * co2ScrubRating, nil
}

func AnswerV3(input string) (int, error) {
	bits := strings.IndexByte(input, '\n')
	numbers := make([]int, (len(input)+bits)/(bits+1))

	for i, j, l := 0, 0, len(input); i < l; i, j = i+bits+1, j+1 {
		v := 0
		for k := 0; k < bits; k++ {
			v |= int(input[i+k]-'0') << (bits - (k + 1))
		}
		numbers[j] = int(v)
	}

	filter := func(arr []int, bit int, invert bool) []int {
		count := 0
		for _, v := range arr {
			count += (v >> bit) & 1
		}

		search := 0
		if count >= (len(arr)+1)/2 {
			search = 1
		}
		if invert {
			search ^= 1
		}

		ret := make([]int, 0, len(arr))
		for _, v := range arr {
			if (v>>bit)&1 == search {
				ret = append(ret, v)
			}
		}

		// fmt.Printf("filter(bit=%d,invert=%t,search=%d) => len(%d)\n", bit, invert, search, len(ret))
		return ret
	}

	find := func(invert bool) int {
		arr := numbers
		for i := bits - 1; i >= 0; i-- {
			arr = filter(arr, i, invert)
			if len(arr) == 1 {
				return arr[0]
			}
		}
		return arr[0]
	}

	o2GenRating, co2ScrubRating := find(false), find(true)
	return o2GenRating * co2ScrubRating, nil
}

func AnswerV4(input string) (int, error) {
	bits := strings.IndexByte(input, '\n')
	numbers1 := make([]uint16, (len(input)+bits)/(bits+1))
	numbers2 := make([]uint16, len(numbers1))

	for i, j, l := 0, 0, len(input); i < l; i, j = i+bits+1, j+1 {
		v := 0
		for k := 0; k < bits; k++ {
			v |= int(input[i+k]-'0') << (bits - (k + 1))
		}
		numbers1[j] = uint16(v)
		numbers2[j] = uint16(v)
	}

	filter := func(arr []uint16, bit int, invert bool) []uint16 {
		count := uint16(0)
		for _, v := range arr {
			count += (v >> bit) & 1
		}

		search := uint16(0)
		if count >= uint16(len(arr)+1)/2 {
			search = 1
		}
		if invert {
			search ^= 1
		}

		s := 0
		for i, l := 0, len(arr); i < l; i++ {
			if (arr[i]>>bit)&1 != search {
				s++
			} else {
				arr[i-s] = arr[i]
			}
		}

		// fmt.Printf("filter(bit=%d,invert=%t,search=%d) => len(%d)\n", bit, invert, search, len(ret))
		return arr[:len(arr)-s]
	}

	find := func(arr []uint16, invert bool) int {
		for i := bits - 1; i >= 0; i-- {
			arr = filter(arr, i, invert)
			if len(arr) == 1 {
				return int(arr[0])
			}
		}
		return int(arr[0])
	}

	o2GenRating, co2ScrubRating := find(numbers1, false), find(numbers2, true)
	return o2GenRating * co2ScrubRating, nil
}

type NodeV5 struct {
	Children [2]*NodeV5
	Counts   [2]uint16
	Value    uint16
}

func AnswerV5(input string) (int, error) {
	bits := strings.IndexByte(input, '\n')
	root := &NodeV5{}

	// Build tree of input
	for i, l := 0, len(input); i < l; i += bits + 1 {
		v := uint16(0)
		selected := root
		for k := 0; k < bits; k++ {
			b := uint16(input[i+k] - '0')
			v |= b << (bits - (k + 1))
			if selected.Children[b] == nil {
				selected.Children[b] = &NodeV5{Value: v}
			}
			selected.Counts[b]++
			selected = selected.Children[b]
		}
	}

	// Find o2GenRating with most popular input
	selected := root
	for i := 0; i < bits; i++ {
		if selected.Counts[1] >= selected.Counts[0] {
			selected = selected.Children[1]
		} else {
			selected = selected.Children[0]
		}
	}
	o2GenRating := int(selected.Value)

	// Find co2ScrubRating with least popular input
	selected = root
	for i := 0; i < bits; i++ {
		if selected.Counts[0] == 0 {
			selected = selected.Children[1]
		} else if selected.Counts[1] == 0 {
			selected = selected.Children[0]
		} else if selected.Counts[0] <= selected.Counts[1] {
			selected = selected.Children[0]
		} else {
			selected = selected.Children[1]
		}
	}
	co2ScrubRating := int(selected.Value)

	return o2GenRating * co2ScrubRating, nil
}
