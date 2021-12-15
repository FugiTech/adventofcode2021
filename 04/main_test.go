package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

var testInput = strings.TrimSpace(`
7,4,9,5,11,17,23,2,0,14,21,24,10,16,13,6,15,25,12,22,18,20,8,19,3,26,1

22 13 17 11  0
 8  2 23  4 24
21  9 14 16  7
 6 10  3 18  5
 1 12 20 15 19

 3 15  0  2 22
 9 18 13 17  5
19  8  7 25 23
20 11 10 24  4
14 21 16 12  6

14 21 17 24  4
10 16 15  9 19
18  8 23 26 20
22 11 13  6  5
 2  0 12  3  7
`)

func getLargeInput() string {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	input, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(input))
}

func TestAnswer(t *testing.T) {
	is := is.New(t)
	answer, err := Answer(testInput)
	is.NoErr(err)
	is.Equal(answer, 1924)
}

func BenchmarkAnswerSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Answer(testInput)
	}
}

func BenchmarkAnswerLarge(b *testing.B) {
	input := getLargeInput()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Answer(input)
	}
}
