package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/set"
)

type pageOrderPair struct {
	l int
	r int
}

type update struct {
	pages []int
}

type checker struct {
	rules map[int]*set.Set[int]
}

func p1() {
	txt := input.NewTXTFile("05.txt")
	var pairs []pageOrderPair
	var updates []update
	txt.ReadByBlock(
		context.Background(),
		"\n\n",
		func(block []string) error {
			rules, datasets := block[0], block[1]
			rparts := strings.Split(rules, "\n")
			for _, rpart := range rparts {
				rules := strings.Split(rpart, "|")
				pairs = append(pairs, pageOrderPair{
					l: input.Atoi(rules[0]),
					r: input.Atoi(rules[1]),
				})
			}
			for _, dataset := range strings.Split(datasets, "\n") {
				nums := strings.Split(dataset, ",")
				var up update
				for _, num := range nums {
					n := input.Atoi(num)
					up.pages = append(up.pages, n)
				}
				updates = append(updates, up)
			}
			return nil
		},
	)
	ck := checker{rules: make(map[int]*set.Set[int])}
	for _, p := range pairs {
		if _, ok := ck.rules[p.r]; !ok {
			ck.rules[p.r] = set.New(p.l)
		} else {
			ck.rules[p.r].Add(p.l)
		}
	}
	var total int
	for _, u := range updates {
		valid := true
		for i := 0; i < len(u.pages)-1; i++ {
			r, ok := ck.rules[u.pages[i]]
			if ok && r.HasAny([]int{u.pages[i+1]}...) {
				valid = false
				break
			}
		}
		if valid {
			midIndex := len(u.pages) / 2
			total += u.pages[midIndex]
		}
	}
	fmt.Printf("p1: %d\n", total)
}

func p2() {
	txt := input.NewTXTFile("05.txt")
	var pairs []pageOrderPair
	var updates []update
	txt.ReadByBlock(
		context.Background(),
		"\n\n",
		func(block []string) error {
			rules, datasets := block[0], block[1]
			rparts := strings.Split(rules, "\n")
			for _, rpart := range rparts {
				rules := strings.Split(rpart, "|")
				pairs = append(pairs, pageOrderPair{
					l: input.Atoi(rules[0]),
					r: input.Atoi(rules[1]),
				})
			}
			for _, dataset := range strings.Split(datasets, "\n") {
				nums := strings.Split(dataset, ",")
				var up update
				for _, num := range nums {
					n := input.Atoi(num)
					up.pages = append(up.pages, n)
				}
				updates = append(updates, up)
			}
			return nil
		},
	)
	ck := checker{rules: make(map[int]*set.Set[int])}
	for _, p := range pairs {
		if _, ok := ck.rules[p.r]; !ok {
			ck.rules[p.r] = set.New(p.l)
		} else {
			ck.rules[p.r].Add(p.l)
		}
	}
	var total int
	for _, u := range updates {
		for i := 0; i < len(u.pages)-1; i++ {
			r, ok := ck.rules[u.pages[i]]
			if ok && r.HasAny([]int{u.pages[i+1]}...) {
				nu := fixUpdate(u, ck)
				total += nu.pages[len(nu.pages)/2]
				break
			}
		}
	}
	fmt.Printf("p2: %d\n", total)
}

func fixUpdate(u update, ck checker) update {
	for i := 0; i < len(u.pages)-1; i++ {
		r, ok := ck.rules[u.pages[i]]
		if ok && r.HasAny([]int{u.pages[i+1]}...) {
			// swap i and i+1
			u.pages[i], u.pages[i+1] = u.pages[i+1], u.pages[i]
			u = fixUpdate(u, ck)
		}
	}
	return u
}

func main() {
	p1()
	p2()
}
