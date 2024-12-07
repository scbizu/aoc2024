package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/queue"
)

type calibrator struct {
	result int
	pairs  []int
}

func (c calibrator) calibration() bool {
	for _, r := range cal(c.pairs) {
		if r == c.result {
			return true
		}
	}
	return false
}

func (c calibrator) calibration2() bool {
	for _, r := range cal2(c.pairs) {
		if r == c.result {
			return true
		}
	}
	return false
}

func opStrAdd(v1, v2 int) int {
	v1Str := fmt.Sprintf("%d", v1)
	v2Str := fmt.Sprintf("%d", v2)
	return input.Atoi(v1Str + v2Str)
}

func cal2(pairs []int) []int {
	// fmt.Printf("pairs: %v\n", pairs)
	q := queue.NewQueue[int]()
	q.Push(0)
	for _, pair := range pairs {
		newQ := queue.NewQueue[int]()
		for q.Len() > 0 {
			v := q.Pop()
			newQ.Push(v + pair)
			newQ.Push(v * pair)
			newQ.Push(opStrAdd(v, pair))
		}
		q = newQ
	}
	var res []int
	for q.Len() > 0 {
		res = append(res, q.Pop())
	}
	// fmt.Printf("res: %v\n", res)
	return res
}

func cal(pairs []int) []int {
	// fmt.Printf("pairs: %v\n", pairs)
	q := queue.NewQueue[int]()
	q.Push(0)
	for _, pair := range pairs {
		newQ := queue.NewQueue[int]()
		for q.Len() > 0 {
			v := q.Pop()
			newQ.Push(v + pair)
			newQ.Push(v * pair)
		}
		q = newQ
	}
	var res []int
	for q.Len() > 0 {
		res = append(res, q.Pop())
	}
	// fmt.Printf("res: %v\n", res)
	return res
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("07.txt")
	var cals []calibrator
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		parts := strings.Split(line, ": ")
		nums := strings.Split(parts[1], " ")
		var pairs []int
		for _, num := range nums {
			pair := input.Atoi(num)
			pairs = append(pairs, pair)
		}
		cals = append(cals, calibrator{
			result: input.Atoi(parts[0]),
			pairs:  pairs,
		})
		return nil
	})
	var total int
	for _, c := range cals {
		if c.calibration() {
			total += c.result
		}
	}
	fmt.Printf("p1: %d\n", total)
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("07.txt")
	var cals []calibrator
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		parts := strings.Split(line, ": ")
		nums := strings.Split(parts[1], " ")
		var pairs []int
		for _, num := range nums {
			pair := input.Atoi(num)
			pairs = append(pairs, pair)
		}
		cals = append(cals, calibrator{
			result: input.Atoi(parts[0]),
			pairs:  pairs,
		})
		return nil
	})
	var total int
	for _, c := range cals {
		if c.calibration2() {
			total += c.result
		}
	}
	fmt.Printf("p2: %d\n", total)
}

func main() {
	ctx := context.TODO()
	p1(ctx)
	p2(ctx)
}
