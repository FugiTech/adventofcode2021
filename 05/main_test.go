package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

var testInput = strings.TrimSpace(`
0,9 -> 5,9
8,0 -> 0,8
9,4 -> 3,4
2,2 -> 2,1
7,0 -> 7,4
6,4 -> 2,0
0,9 -> 2,9
3,4 -> 1,4
0,0 -> 8,8
5,5 -> 8,2
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
	is.Equal(answer, 12)
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

func BenchmarkNormalize(b *testing.B) {
	b.Skip()

	b.Run("slow", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			normalizeSlow(-10)
			normalizeSlow(-5)
			normalizeSlow(-1)
			normalizeSlow(0)
			normalizeSlow(1)
			normalizeSlow(10)
			normalizeSlow(5)
		}
	})

	b.Run("fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			normalize(-10)
			normalize(-5)
			normalize(-1)
			normalize(0)
			normalize(1)
			normalize(10)
			normalize(5)
		}
	})
}
