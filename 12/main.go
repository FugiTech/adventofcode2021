package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
	paths := map[string][]string{}

	lines := strings.Split(input, "\n")
	for _, l := range lines {
		p := strings.Split(l, "-")
		paths[p[0]] = append(paths[p[0]], p[1])
		paths[p[1]] = append(paths[p[1]], p[0])
	}

	var visit func(from string, visited []string) int
	visit = func(from string, visited []string) int {
		if from == "end" {
			return 1
		}

		visited = append(visited, from)
		count := 0
		for _, to := range paths[from] {
			if to[0] <= 'Z' || !contains(visited, to) {
				count += visit(to, visited)
			}
		}
		return count
	}

	return visit("start", nil), nil
}

func AnswerV1(input string) (int, error) {
	// Contains a map of "location" to "all locations reachable from that location"
	paths := map[string][]string{}

	lines := strings.Split(input, "\n")
	for _, l := range lines {
		p := strings.Split(l, "-")
		// Paths are bidirectional, so add to both keys
		paths[p[0]] = append(paths[p[0]], p[1])
		paths[p[1]] = append(paths[p[1]], p[0])
	}

	// visit returns the number of valid paths from "from" to the 'end'
	// "doubled" stores whether we've used our one-time double visit to a small cave
	var visit func(from string, visited []string, doubled bool) int
	visit = func(from string, visited []string, doubled bool) int {
		// If we found the end, immediately end
		if from == "end" {
			return 1
		}

		// Otherwise remember we've been here so we don't double back
		visited = append(visited, from)
		count := 0
		for _, to := range paths[from] {
			if to == "start" {
				// Can't go to start twice
			} else if to[0] <= 'Z' || !contains(visited, to) {
				// Either a big-cave or unvisited small-cave, either way it's free to go there
				count += visit(to, visited, doubled)
			} else if !doubled {
				// It's a visited small-cave but we have our double-visit ticket still. Use it.
				count += visit(to, visited, true)
			}
		}
		return count
	}

	return visit("start", nil, false), nil
}

func contains(a []string, s string) bool {
	for _, v := range a {
		if s == v {
			return true
		}
	}
	return false
}
