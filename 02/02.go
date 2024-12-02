package main

import (
	"context"
	"slices"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/math"
)

type rule struct {
	nums []int
}

func isSafe2(nums []int) bool {
	if !isSafe(nums) {
		for i := 0; i < len(nums); i++ {
			if isSafe(slices.Delete(slices.Clone(nums), i, i+1)) {
				return true
			}
		}
	} else {
		return true
	}
	return false
}

func isSafe(nums []int) bool {
	// fmt.Printf("check: %v\n", nums)
	var desc bool
	first, second := nums[0], nums[1]
	if first == second {
		return false
	}
	if first < second {
		desc = false
	} else {
		desc = true
	}
	for i := 1; i < len(nums); i++ {
		v := nums[i] - nums[i-1]
		if desc && v >= 0 {
			return false
		}
		if !desc && v <= 0 {
			return false
		}
		if math.Abs(v) < 1 || math.Abs(v) > 3 {
			return false
		}
	}
	return true
}

func p1() {
	txt := input.NewTXTFile("02.txt")
	var rules []rule
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		parts := strings.Split(line, " ")
		var nums []int
		for _, part := range parts {
			nums = append(nums, input.Atoi(part))
		}
		rules = append(rules, rule{nums: nums})
		return nil
	})
	var count int
	for _, r := range rules {
		if isSafe(r.nums) {
			count++
		}
	}
	println("p1: ", count)
}

func p2() {
	txt := input.NewTXTFile("02.txt")
	var rules []rule
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		parts := strings.Split(line, " ")
		var nums []int
		for _, part := range parts {
			nums = append(nums, input.Atoi(part))
		}
		rules = append(rules, rule{nums: nums})
		return nil
	})
	var count int
	for _, r := range rules {
		if isSafe2(r.nums) {
			count++
		}
	}
	println("p2: ", count)
}

func main() {
	p1()
	p2()
}
