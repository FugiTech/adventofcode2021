package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
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

var Answer = AnswerNelson

func AnswerPart1(input string) (int, error) {
	count := 0

	in := strings.Split(input, "\n")
	for _, l := range in {
		p := strings.Split(l, " | ")
		out := strings.Split(p[1], " ")
		for _, o := range out {
			if len(o) < 5 || len(o) > 6 {
				count++
			}
		}
	}

	return count, nil
}

func AnswerV1(input string) (int, error) {
	count := 0

	lines := strings.Split(input, "\n")
	for _, l := range lines {
		p := strings.Split(l, " | ")
		in := strings.Split(p[0], " ")
		out := strings.Split(p[1], " ")

		sort.Slice(in, func(a, b int) bool {
			return len(in[a]) < len(in[b])
		})

		keyToNum := map[uint8]int{}
		numToKey := map[int]uint8{}
		add := func(k uint8, n int) {
			keyToNum[k] = n
			numToKey[n] = k
		}

		add(segmentKey(in[0]), 1)
		add(segmentKey(in[1]), 7)
		add(segmentKey(in[2]), 4)
		add(segmentKey(in[9]), 8)

		for i := 6; i <= 8; i++ { // 6 segment nums
			v := segmentKey(in[i])
			four := numToKey[4]
			seven := numToKey[7]
			if v&four == four {
				add(v, 9)
			} else if v&seven == seven {
				add(v, 0)
			} else {
				add(v, 6)
			}
		}

		for i := 3; i <= 5; i++ { // 5 segment nums
			v := segmentKey(in[i])
			seven := numToKey[7]
			six := numToKey[6]
			if v&seven == seven {
				add(v, 3)
			} else if v&six == v {
				add(v, 5)
			} else {
				add(v, 2)
			}
		}

		count += (1000*keyToNum[segmentKey(out[0])] +
			100*keyToNum[segmentKey(out[1])] +
			10*keyToNum[segmentKey(out[2])] +
			1*keyToNum[segmentKey(out[3])])
	}

	return count, nil
}

func AnswerV2(input string) (int, error) {
	count := 0

	var (
		keyToNum  [256]int
		numToKey  [10]uint8
		in        uint8
		inLen     int
		fivers    [3]uint8
		fiversIdx int
		sixers    [3]uint8
		sixersIdx int
	)

	for i, l := 0, len(input); i < l; i++ {
		if input[i] == '|' {
			for idx := 0; idx < sixersIdx; idx++ {
				v := sixers[idx]
				if v&numToKey[4] == numToKey[4] {
					keyToNum[v] = 9
					numToKey[9] = v
				} else if v&numToKey[7] == numToKey[7] {
					keyToNum[v] = 0
					numToKey[0] = v
				} else {
					keyToNum[v] = 6
					numToKey[6] = v
				}
			}

			for idx := 0; idx < sixersIdx; idx++ {
				v := fivers[idx]
				if v&numToKey[7] == numToKey[7] {
					keyToNum[v] = 3
					numToKey[3] = v
				} else if v&numToKey[6] == v {
					keyToNum[v] = 5
					numToKey[5] = v
				} else {
					keyToNum[v] = 2
					numToKey[2] = v
				}
			}

			fiversIdx = 0
			sixersIdx = 0

			out := 0
			for i += 2; i < l && input[i] != '\n'; i++ {
				if input[i] != ' ' {
					in |= 1 << (input[i] - 'a')
				} else {
					out = 10*out + keyToNum[in]
					in = 0
				}
			}
			out = 10*out + keyToNum[in]
			in = 0

			count += out
		} else if input[i] != ' ' {
			in |= 1 << (input[i] - 'a')
			inLen++
		} else {
			switch inLen {
			case 2:
				keyToNum[in] = 1
				numToKey[1] = in
			case 3:
				keyToNum[in] = 7
				numToKey[7] = in
			case 4:
				keyToNum[in] = 4
				numToKey[4] = in
			case 5:
				fivers[fiversIdx] = in
				fiversIdx++
			case 6:
				sixers[sixersIdx] = in
				sixersIdx++
			case 7:
				keyToNum[in] = 8
				numToKey[8] = in
			}

			in = 0
			inLen = 0
		}
	}

	return count, nil
}

func segmentKey(s string) uint8 {
	m := map[byte]uint8{
		'a': 1 << 7,
		'b': 1 << 6,
		'c': 1 << 5,
		'd': 1 << 4,
		'e': 1 << 3,
		'f': 1 << 2,
		'g': 1 << 1,
	}
	var r uint8
	for _, c := range []byte(s) {
		r |= m[c]
	}
	return r
}

func AnswerNelson(input string) (int, error) {
	inputs := [][]string{}
	outputs := [][]string{}

	lines := strings.Split(input, "\n")
	for _, l := range lines {
		inout := strings.Split(l, "|")
		in := strings.Split(inout[0], " ")
		inputs = append(inputs, in)
		out := strings.Split(inout[1], " ")
		outputs = append(outputs, out)
	}

	outputNums := []int{}

	for i, input := range inputs {
		mapping := map[int]number{}
		filteredList := []string{}
		filteredList = append(filteredList, input...)
		// Get 1, 4, 7, and 8
		for _, str := range input {
			if len(str) == 2 {
				num := createNumber(1, str)
				mapping[1] = num
				filteredList = removeStrFromList(filteredList, str)
				continue
			}
			if len(str) == 4 {
				num := createNumber(4, str)
				mapping[4] = num
				filteredList = removeStrFromList(filteredList, str)
				continue
			}
			if len(str) == 3 {
				num := createNumber(7, str)
				mapping[7] = num
				filteredList = removeStrFromList(filteredList, str)
				continue
			}

			if len(str) == 7 {
				num := createNumber(8, str)
				mapping[8] = num
				filteredList = removeStrFromList(filteredList, str)
				continue
			}
		}
		// Figure out 5s
		for _, str := range filteredList {
			if len(str) != 5 {
				continue
			}
			strMap := map[string]bool{}
			for _, letter := range str {
				strMap[string(letter)] = true
			}
			// If it shares 2 letters with 1, and has 5 letters, it must be 3
			if checkSimilarLetters(mapping[1].letters, strMap) == 2 {
				mapping[3] = createNumber(3, str)
				filteredList = removeStrFromList(filteredList, str)
			}
		}

		for _, str := range filteredList {
			if len(str) != 5 {
				continue
			}
			strMap := map[string]bool{}
			for _, letter := range str {
				strMap[string(letter)] = true
			}
			// Shares 3 letters with 4, must be 5
			if checkSimilarLetters(mapping[4].letters, strMap) == 3 {
				mapping[5] = createNumber(5, str)
				filteredList = removeStrFromList(filteredList, str)
			}
		}

		// Last number must be 2
		for _, str := range filteredList {
			if len(str) != 5 {
				continue
			}
			mapping[2] = createNumber(2, str)
			filteredList = removeStrFromList(filteredList, str)
		}

		// Need to find 0, 6, 9
		for _, str := range filteredList {
			if len(str) != 6 {
				continue
			}
			strMap := map[string]bool{}
			for _, letter := range str {
				strMap[string(letter)] = true
			}
			// 9 shares 4 with 4
			if checkSimilarLetters(mapping[4].letters, strMap) == 4 {
				mapping[9] = createNumber(9, str)
				filteredList = removeStrFromList(filteredList, str)
			}
		}

		// 6 vs 0
		for _, str := range filteredList {
			if len(str) != 6 {
				continue
			}
			strMap := map[string]bool{}
			for _, letter := range str {
				strMap[string(letter)] = true
			}
			//
			if checkSimilarLetters(mapping[7].letters, strMap) == 3 {
				mapping[0] = createNumber(0, str)
				filteredList = removeStrFromList(filteredList, str)
			}
		}

		// Last number must be 6
		for _, str := range filteredList {
			if len(str) != 6 {
				continue
			}
			mapping[6] = createNumber(6, str)
			filteredList = removeStrFromList(filteredList, str)
		}

		if len(mapping) != 10 {
			log.Fatal("wtf", mapping)
		}

		outputStr := ""
		for _, out := range outputs[i] {
			if len(out) == 0 {
				continue
			}
			num := strMap2Number(mapping, out)
			if num == 10 {
				log.Fatal("oh no 10")
			} // shitty error handle
			outputStr += fmt.Sprintf("%d", num)
		}
		num, _ := strconv.ParseInt(outputStr, 10, 64)
		outputNums = append(outputNums, int(num))
	}

	total := 0
	for _, num := range outputNums {
		total += num
	}

	return total, nil
}

type number struct {
	number  int
	letters map[string]bool
}

func strMap2Number(mapping map[int]number, letters string) int {
	strMap := map[string]bool{}
	for _, letter := range letters {
		strMap[string(letter)] = true
	}

	for key, value := range mapping {
		if reflect.DeepEqual(value.letters, strMap) {
			return key
		}
	}
	fmt.Println(mapping)
	fmt.Println(letters)
	return 10
}

func checkSimilarLetters(a, b map[string]bool) int {
	total := 0
	for key, _ := range a {
		_, ok := b[key]
		if ok {
			total++
		}
	}
	return total
}

func createNumber(num int, letters string) number {
	number := number{
		number:  num,
		letters: map[string]bool{},
	}

	for _, letter := range letters {
		number.letters[string(letter)] = true
	}

	return number
}

func removeStrFromList(list []string, str string) []string {
	output := []string{}
	for _, entry := range list {
		if entry != str {
			output = append(output, entry)
		}
	}

	return output
}
