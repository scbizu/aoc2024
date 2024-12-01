package main

import (
	"context"
	"slices"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/math"
)

type left struct {
	numbers []int
}

type right struct {
	numbers []int
}

func p1() {
	txt := input.NewTXTFile("01.txt")
	ctx := context.Background()
	l, r := left{}, right{}
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		parts := strings.Split(line, "   ")
		ln, rn := parts[0], parts[1]
		l.numbers = append(l.numbers, input.Atoi(ln))
		r.numbers = append(r.numbers, input.Atoi(rn))
		return nil
	})
	slices.Sort(l.numbers)
	slices.Sort(r.numbers)

	var sum int
	for i := 0; i < len(l.numbers); i++ {
		sum += math.Abs(l.numbers[i] - r.numbers[i])
	}
	println("p1: ", sum)
}

func p2() {
	txt := input.NewTXTFile("01.txt")
	ctx := context.Background()
	l, r := left{}, right{}
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		parts := strings.Split(line, "   ")
		ln, rn := parts[0], parts[1]
		l.numbers = append(l.numbers, input.Atoi(ln))
		r.numbers = append(r.numbers, input.Atoi(rn))
		return nil
	})
	rmap := make(map[int]int)
	for i := 0; i < len(r.numbers); i++ {
		if _, ok := rmap[r.numbers[i]]; !ok {
			rmap[r.numbers[i]] = 1
		} else {
			rmap[r.numbers[i]] += 1
		}
	}
	var sum int
	for i := 0; i < len(l.numbers); i++ {
		if _, ok := rmap[l.numbers[i]]; ok {
			sum += l.numbers[i] * rmap[l.numbers[i]]
		}
	}
	println("p2: ", sum)
}

func main() {
	p1()
	p2()
}
