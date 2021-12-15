package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

var testInputA = strings.TrimSpace(`
start-A
start-b
A-c
A-b
b-d
A-end
b-end
`)

var testInputB = strings.TrimSpace(`
dc-end
HN-start
start-kj
dc-start
dc-HN
LN-dc
HN-end
kj-sa
kj-HN
kj-dc
`)

var testInputC = strings.TrimSpace(`
fs-end
he-DX
fs-he
start-DX
pj-DX
end-zg
zg-sl
zg-pj
pj-he
RW-he
fs-DX
pj-RW
zg-RW
start-pj
he-WI
zg-he
pj-fs
start-RW
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

func TestAnswerA(t *testing.T) {
	is := is.New(t)
	answer, err := Answer(testInputA)
	is.NoErr(err)
	is.Equal(answer, 36)
}
func TestAnswerB(t *testing.T) {
	is := is.New(t)
	answer, err := Answer(testInputB)
	is.NoErr(err)
	is.Equal(answer, 103)
}
func TestAnswerC(t *testing.T) {
	is := is.New(t)
	answer, err := Answer(testInputC)
	is.NoErr(err)
	is.Equal(answer, 3509)
}

func BenchmarkAnswerSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Answer(testInputA)
	}
}

func BenchmarkAnswerLarge(b *testing.B) {
	input := getLargeInput()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Answer(input)
	}
}
