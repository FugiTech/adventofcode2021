package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

var testInput = strings.TrimSpace(`
NNCB

CH -> B
HH -> N
CB -> H
NH -> C
HB -> C
HC -> B
HN -> C
NN -> C
BH -> H
NC -> B
NB -> B
BN -> B
BB -> N
BC -> B
CC -> N
CN -> C
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
	is.Equal(answer, 2188189693529)
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
