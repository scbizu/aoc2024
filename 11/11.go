package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
)

type game struct {
	cache map[string][]string
}

func (g *game) load2(pattern map[string]int) map[string]int {
	next := make(map[string]int)
	for k, v := range pattern {
		if v == 0 {
			continue
		}
		for _, num := range g.load(k) {
			next[num] += v
		}
	}
	return next
}

func (g *game) load(num string) []string {
	if v, ok := g.cache[num]; ok {
		return v
	}
	var res []string
	switch {
	case num == "0":
		res = []string{"1"}
	case len(num)%2 == 0:
		half := len(num) / 2
		left := num[:half]
		right := strings.TrimLeft(num[half:], "0")
		if right == "" {
			right = "0"
		}
		res = []string{left, right}
	default:
		res = []string{fmt.Sprintf("%d", input.Atoi(num)*2024)}
	}
	g.cache[num] = res
	return res
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("11.txt")
	var nums []string
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		nums = strings.Split(line, " ")
		return nil
	})
	g := game{
		cache: make(map[string][]string),
	}
	var blink int
	for {
		if blink > 24 {
			break
		}
		var r []string
		for _, num := range nums {
			r = append(r, g.load(num)...)
		}
		nums = r
		blink++
	}
	fmt.Printf("p1: %v\n", len(nums))
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("11.txt")
	var nums []string
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		nums = strings.Split(line, " ")
		return nil
	})
	g := game{
		cache: make(map[string][]string),
	}
	pattern := make(map[string]int)
	for _, num := range nums {
		pattern[num]++
	}
	var blink int
	var total int
	for {
		if blink > 74 {
			break
		}
		pattern = g.load2(pattern)
		// fmt.Printf("pattern: %v\n", pattern)
		blink++
	}
	for _, v := range pattern {
		total += v
	}
	fmt.Printf("p2: %d\n", total)
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
