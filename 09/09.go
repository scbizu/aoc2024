package main

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
)

type diskMap struct {
	rawStr string
	parsed []byte
}

func (dm *diskMap) parse(_ context.Context) {
	var index int
	parsed := bytes.NewBuffer(nil)
	for index < len(dm.rawStr) {
		// index % 2 == 0 indicates the file blocks
		// index % 2 == 1 indicates the free space blocks
		if index%2 == 0 {
			// fill n bytes of repeated index
			n := dm.rawStr[index]
			indexStr := fmt.Sprintf("%v", string((index/2 + '0')))
			fmt.Fprintf(parsed, "%s", strings.Repeat(indexStr, input.Atoi(string(n))))
		}
		if index%2 == 1 {
			// fill n bytes of repeated dot
			n := dm.rawStr[index]
			fmt.Fprintf(parsed, "%s", strings.Repeat(".", input.Atoi(string(n))))
		}
		index++
	}
	dm.parsed = parsed.Bytes()
}

func (dm diskMap) refmt(_ context.Context) []byte {
	var firstDotIndex int
	var lastNonDotIndex int
	newBytes := slices.Clone(dm.parsed)
	for {
		// fmt.Printf("first: %d,last: %d\n", firstDotIndex, lastNonDotIndex)
		firstDotIndex = strings.IndexByte(string(newBytes), '.')
		lastNonDotIndex = strings.LastIndexFunc(string(newBytes), func(r rune) bool {
			return r != '.'
		})
		if firstDotIndex > lastNonDotIndex {
			break
		}
		newBytes[firstDotIndex] = dm.parsed[lastNonDotIndex]
		newBytes[lastNonDotIndex] = '.'
		// fmt.Printf("newBytes: %s\n", string(newBytes))
	}
	return newBytes
}

func checksum(bs []byte) int {
	var sum int
	for i := 0; i < len(bs); i++ {
		if bs[i] == '.' {
			break
		}
		sum += int(bs[i]-'0') * i
	}
	return sum
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("09.txt")
	dm := &diskMap{}
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		dm.rawStr = line
		return nil
	})
	dm.parse(ctx)
	fmt.Printf("parsed: %s\n", string(dm.parsed))
	rf := dm.refmt(ctx)
	fmt.Printf("refmt: %s\n", string(rf))
	total := checksum(rf)
	fmt.Printf("p1: %d\n", total)
}

func main() {
	ctx := context.TODO()
	p1(ctx)
}
