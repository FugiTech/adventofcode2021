package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

type test struct {
	Input  string
	Answer int
}

var tests = []test{
	/*
		{"EE00D40C823060", 14},
		{"8A004A801A8002F478", 16},
		{"620080001611562C8802118E34", 12},
		{"C0015000016115A2E0802F182340", 23},
		{"A0016C880162017C3686B18A3D4780", 31},
	*/
	{"C200B40A82", 3},
	{"04005AC33890", 54},
	{"880086C3E88112", 7},
	{"CE00C43D881120", 9},
	{"D8005AC2A8F0", 1},
	{"F600BC2D8F", 0},
	{"9C005AC2F8F0", 0},
	{"9C0141080250320F1802104A08", 1},
}

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
	for _, t := range tests {
		answer, err := Answer(t.Input)
		is.NoErr(err)
		is.Equal(answer, t.Answer)
	}
}

func BenchmarkAnswerSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Answer(tests[0].Input)
	}
}

func BenchmarkAnswerLarge(b *testing.B) {
	input := getLargeInput()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Answer(input)
	}
}
