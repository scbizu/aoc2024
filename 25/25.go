package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
)

type key struct {
	m grid.VecMatrix[byte]
}

func (k key) toPattern() []int {
	var rows []int
	for col := 0; col < k.m.Cols(); col++ {
		var c int
		for row := k.m.Rows() - 2; row >= 0; row-- {
			if k.m[grid.Vec{X: col, Y: row}] != '#' {
				continue
			}
			c++
		}
		rows = append(rows, c)
	}
	return rows
}

type lock struct {
	m grid.VecMatrix[byte]
}

func (l lock) toPattern() []int {
	var rows []int
	for col := 0; col < l.m.Cols(); col++ {
		var c int
		for row := 1; row < l.m.Rows(); row++ {
			if l.m[grid.Vec{X: col, Y: row}] != '#' {
				continue
			}
			c++
		}
		rows = append(rows, c)
	}
	return rows
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("25.txt")
	var keys []key
	var locks []lock
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		for _, b := range block {
			m := grid.NewVecMatrix[byte]()
			for y, line := range strings.Split(b, "\n") {
				for x, c := range line {
					m.Add(grid.Vec{X: x, Y: y}, byte(c))
				}
			}
			if m[grid.Vec{X: 0, Y: 0}] == '.' {
				keys = append(keys, key{m: m})
			}
			if m[grid.Vec{X: 0, Y: 0}] == '#' {
				locks = append(locks, lock{m: m})
			}
		}
		return nil
	})
	// fmt.Printf("keys: %d;locks: %d\n", len(keys), len(locks))
	col, row := keys[0].m.Cols(), keys[0].m.Rows()
	// fmt.Printf("cols: %d,rows: %d\n", col, row)
	var fit int
	for _, k := range keys {
		for _, l := range locks {
			var isFit bool
			// fmt.Printf("[key:%d] pattern: %v\n", ki, k.toPattern())
			// fmt.Printf("[lock:%d] pattern: %v\n", li, l.toPattern())
			for i := 0; i < col; i++ {
				// fmt.Printf("[%d:%d]num: %d to %d\n", ki, li, k.toPattern()[i], l.toPattern()[i])
				if k.toPattern()[i]+l.toPattern()[i] <= row-2 {
					isFit = true
				} else {
					isFit = false
					break
				}
			}
			if isFit {
				// fmt.Printf("keys: %v\n; locks: %v\n", k, l)
				fit++
			}
		}
	}
	fmt.Printf("p1: %d\n", fit)
}

func main() {
	ctx := context.Background()
	p1(ctx)
}
