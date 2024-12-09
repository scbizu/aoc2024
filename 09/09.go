package main

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/stack"
)

type diskMap struct {
	rawStr     string
	parsed     []rune
	blockStack *stack.Stack[*block]
}

type block struct {
	start, end int
	value      rune
}

func (dm *diskMap) parse(_ context.Context) {
	var index int
	parsed := bytes.NewBuffer(nil)
	var size int
	for index < len(dm.rawStr) {
		// index % 2 == 0 indicates the file blocks
		// index % 2 == 1 indicates the free space blocks
		n := dm.rawStr[index]
		if index%2 == 0 {
			// fill n bytes of repeated index
			indexStr := fmt.Sprintf("%v", string((index/2 + '0')))
			fmt.Fprintf(parsed, "%s", strings.Repeat(indexStr, input.Atoi(string(n))))
			dm.blockStack.Push(&block{start: size, end: size + input.Atoi(string(n)), value: rune(index/2 + '0')})
		}
		if index%2 == 1 {
			// fill n bytes of repeated dot
			fmt.Fprintf(parsed, "%s", strings.Repeat(".", input.Atoi(string(n))))
		}
		size += input.Atoi(string(n))
		index++
	}
	dm.parsed = []rune(parsed.String())
}

func (dm *diskMap) refmt2(_ context.Context) []rune {
	newBytes := slices.Clone(dm.parsed)
	rd := slices.Clone(dm.parsed)
	for {
		b := dm.blockStack.Pop()
		if b == nil {
			break
		}
		// fmt.Printf("b: %v\n", b)
		bsize := b.end - b.start
		var ack bool
		var start, end int
		for i := 0; i < len(rd); i++ {
			if rd[i] == '.' {
				// fmt.Printf("i: %d\n", i)
				if i+bsize > len(rd) {
					break
				}
				// fmt.Printf("i: %d, bsize: %d\n", i, bsize)
				ack = true
				for j := 0; j < bsize; j++ {
					if rd[i+j] != '.' {
						ack = false
						break
					}
				}
				if ack {
					start = i
					end = i + bsize
					break
				}
			}
		}
		// acked new location and new location is better :)
		if ack && start < b.start {
			for i := start; i < end; i++ {
				newBytes[i] = b.value
				rd[i] = b.value
			}
			for i := b.start; i < b.end; i++ {
				newBytes[i] = '.'
			}
		}
		// fmt.Printf("newBytes: %s\n", string(newBytes))
	}
	return newBytes
}

func (dm diskMap) refmt(_ context.Context) []rune {
	var firstDotIndex int
	var lastNonDotIndex int
	newBytes := slices.Clone(dm.parsed)
	for {
		// fmt.Printf("first: %d,last: %d\n", firstDotIndex, lastNonDotIndex)
		for firstDotIndex = 0; firstDotIndex < len(newBytes); firstDotIndex++ {
			if newBytes[firstDotIndex] == '.' {
				break
			}
		}
		for lastNonDotIndex = len(newBytes) - 1; lastNonDotIndex >= 0; lastNonDotIndex-- {
			if newBytes[lastNonDotIndex] != '.' {
				break
			}
		}
		if firstDotIndex > lastNonDotIndex {
			break
		}
		newBytes[firstDotIndex] = dm.parsed[lastNonDotIndex]
		newBytes[lastNonDotIndex] = '.'
		// fmt.Printf("newBytes: %s\n", string(newBytes))
	}
	return newBytes
}

func checksum(bs []rune) int {
	var sum int
	for i := 0; i < len(bs); i++ {
		if bs[i] == '.' {
			continue
		}
		sum += int(bs[i]-'0') * i
	}
	return sum
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("09.txt")
	dm := &diskMap{
		blockStack: stack.NewStack[*block](),
	}
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		dm.rawStr = line
		return nil
	})
	dm.parse(ctx)
	// fmt.Printf("parsed: %s\n", string(dm.parsed))
	rf := dm.refmt(ctx)
	// fmt.Printf("refmt: %s\n", string(rf))
	total := checksum(rf)
	fmt.Printf("p1: %d\n", total)
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("09.txt")
	dm := &diskMap{
		blockStack: stack.NewStack[*block](),
	}
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		dm.rawStr = line
		return nil
	})
	dm.parse(ctx)
	// fmt.Printf("parsed: %s\n", string(dm.parsed))
	rf := dm.refmt2(ctx)
	// fmt.Printf("refmt: %s\n", string(rf))
	total := checksum(rf)
	fmt.Printf("p2: %d\n", total)
}

func main() {
	ctx := context.TODO()
	p1(ctx)
	p2(ctx)
}
