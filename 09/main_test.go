package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

var testInput = strings.TrimSpace(`
2199943210
3987894921
9856789892
8767896789
9899965678
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
	is.Equal(answer, 1134)
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
