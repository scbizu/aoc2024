package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
)

type secret struct{}

func (s *secret) grow(num, n int) int {
	for i := 0; i < n; i++ {
		num = (num*64 ^ num) % 16777216
		num = (num/32 ^ num) % 16777216
		num = (num*2048 ^ num) % 16777216
	}
	return num
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("22.txt")
	var initials []int
	txt.ReadByLineEx(ctx, func(_ int, line string) error {
		initials = append(initials, input.Atoi(line))
		return nil
	})
	var sum int
	n := 2000
	for _, initial := range initials {
		s := &secret{}
		g := s.grow(initial, n)
		// fmt.Printf("g: %d\n", g)
		sum += g
	}
	fmt.Printf("p1: %d\n", sum)
}

type indexNumber struct {
	index int
	v     int
}

func ints2ToStrs(ints []int) []string {
	var out []string
	for _, i := range ints {
		out = append(out, fmt.Sprintf("%d", i))
	}
	return out
}

type delta struct {
	d int
	v int
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("22.txt")
	var initials []int
	txt.ReadByLineEx(ctx, func(_ int, line string) error {
		initials = append(initials, input.Atoi(line))
		return nil
	})
	n := 2000
	// seq per buyer
	seqmap := make(map[string][]int)
	for index, initial := range initials {
		s := &secret{}
		var prices []int
		for i := 0; i < n; i++ {
			g := s.grow(initial, i)
			prices = append(prices, g%10)
		}
		// fmt.Printf("prices: %v\n", prices)
		var deltas []delta
		for i := 0; i < len(prices)-1; i++ {
			j := i + 1
			deltas = append(deltas, delta{prices[j] - prices[i], prices[j]})
		}
		// fmt.Printf("deltas: %v\n", deltas)
		var every4 []int
		for i := 0; i < len(deltas); i++ {
			if len(every4) == 4 {
				strs := strings.Join(ints2ToStrs(every4), ",")
				// monkey will sell it at once
				if _, ok := seqmap[strs]; !ok {
					// make seq per buyer
					seqmap[strs] = make([]int, len(initials))
					bp := seqmap[strs]
					bp[index] = deltas[i-1].v
					seqmap[strs] = bp
				} else {
					bp := seqmap[strs]
					if bp[index] == 0 {
						bp[index] = deltas[i-1].v
						seqmap[strs] = bp
					}
				}
				// if bp[index] < deltas[i-1].v {
				// }
				// remove first
				every4 = slices.Delete(every4, 0, 1)
			}
			every4 = append(every4, deltas[i].d)
		}
	}
	// fmt.Printf("seqmap: %v\n", seqmap)
	var max int
	for _, v := range seqmap {
		var sum int
		for _, i := range v {
			sum += i
		}
		if sum > max {
			// fmt.Printf("seq: %s, v: %v\n", seq, v)
			max = sum
		}
	}
	fmt.Printf("p2: %d\n", max)
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
