package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/queue"
	"github.com/magejiCoder/magejiAoc/set"
)

type towel struct {
	stripes *set.Set[string]
}

type state struct {
	read []string
	left string
}

func (t towel) spilitPattern(pattern []byte) [][]string {
	// fmt.Printf("pattern: %s\n", string(pattern))
	q := queue.NewQueue[state]()
	q.Push(state{
		read: []string{string(pattern[:1])},
		left: string(pattern[1:]),
	})
	var patterns [][]string
	for q.Len() > 0 {
		s := q.Pop()
		fmt.Printf("pop: %s\n", s)
		latest := s.read[len(s.read)-1]
		if len(s.left) == 0 {
			// handle the last byte
			patterns = append(patterns, s.read)
		}
		if len(s.left) > 0 {
			fmt.Printf("next1: %s\n", s.left[:1])
			if t.stripes.Has(string(s.left[:1])) {
				// append 1 byte
				q.Push(state{
					read: append(slices.Clone(s.read), string(s.left[:1])),
					left: string(s.left[1:]),
				})
			}
			fmt.Printf("next2: %s\n", latest+s.left[:1])
			if t.stripes.Has(latest + s.left[:1]) {
				// append to latest byte
				q.Push(state{
					read: append(slices.Clone(s.read[:len(s.read)-1]), latest+s.left[:1]),
					left: string(s.left[1:]),
				})
			}
		}
	}
	return patterns
}

func (t towel) isPatternValid(
	pattern string,
) bool {
	for _, ps := range t.spilitPattern([]byte(pattern)) {
		valid := true
		// fmt.Printf("ps: %v\n", ps)
		// every ps should present in stripes
		for _, p := range ps {
			if !t.stripes.Has(p) {
				fmt.Printf("don't present in strips: %v\n", p)
				valid = false
				break
			}
		}
		if valid {
			return true
		}
	}
	return false
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("19.txt")
	t := &towel{
		stripes: set.New[string](),
	}
	var patterns []string
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		// block0 is the stripes
		stripes := block[0]
		for _, stripe := range strings.Split(stripes, ", ") {
			t.stripes.Add(stripe)
		}
		// block1 is the patterns
		patterns = append(patterns, strings.Split(block[1], "\n")...)
		return nil
	})
	var total int
	for _, pattern := range patterns {
		if t.isPatternValid(pattern) {
			fmt.Printf("valid: %s\n", pattern)
			total++
		}
	}
	fmt.Printf("p1: %d\n", total)
}

func main() {
	ctx := context.Background()
	p1(ctx)
}
