package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/set"
)

type direction int

const (
	up direction = iota
	down
	left
	right
)

type panel struct {
	m          grid.VecMatrix[byte]
	startAt    grid.Vec
	startAtDir direction
	currentDir direction
}

func (p panel) String() string {
	var buf bytes.Buffer
	for i := 0; i < p.m.Rows(); i++ {
		for j := 0; j < p.m.Cols(); j++ {
			v := grid.Vec{X: j, Y: i}
			b, ok := p.m.Get(v)
			if !ok {
				buf.WriteRune(' ')
			} else {
				buf.WriteByte(b)
			}
		}
		buf.WriteRune('\n')
	}
	return buf.String()
}

func (p panel) clone() panel {
	m := grid.NewVecMatrix[byte]()
	p.m.ForEach(func(v grid.Vec, b byte) {
		m.Add(v, b)
	})
	return panel{
		m:          m,
		startAt:    p.startAt,
		startAtDir: p.startAtDir,
		currentDir: p.currentDir,
	}
}

type traceItem struct {
	v grid.Vec
	d direction
}

func findPaths(_ context.Context, pn panel) (*set.Set[grid.Vec], bool) {
	v := pn.startAt
	traces := set.New[traceItem](traceItem{v: v, d: pn.currentDir})
	paths := set.New[grid.Vec](v)
	for {
		if _, ok := pn.m.Get(v); !ok {
			return paths, false
		}
		var foundNext bool
		for {
			if foundNext {
				break
			}
			switch pn.currentDir {
			case up:
				next := v.Add(grid.Vec{X: 0, Y: -1})
				nextV, ok := pn.m.Get(next)
				if !ok {
					return paths, false
				}
				if nextV == '#' {
					// turn right
					pn.currentDir = right
				} else {
					v = next
					foundNext = true
				}
			case down:
				next := v.Add(grid.Vec{X: 0, Y: 1})
				nextV, ok := pn.m.Get(next)
				if !ok {
					return paths, false
				}
				if nextV == '#' {
					// turn left
					// next = v.Add(grid.Vec{X: -1, Y: 0})
					pn.currentDir = left
				} else {
					v = next
					foundNext = true
				}
			case left:
				next := v.Add(grid.Vec{X: -1, Y: 0})
				nextV, ok := pn.m.Get(next)
				if !ok {
					return paths, false
				}
				if nextV == '#' {
					// turn up
					pn.currentDir = up
				} else {
					v = next
					foundNext = true
				}
			case right:
				next := v.Add(grid.Vec{X: 1, Y: 0})
				nextV, ok := pn.m.Get(next)
				if !ok {
					return paths, false
				}
				if nextV == '#' {
					// turn down
					pn.currentDir = down
				} else {
					v = next
					foundNext = true
				}
			}
		}
		// fmt.Printf("next at %v: direction: %v\n", v, pn.currentDir)
		// detect loop
		if traces.Has(traceItem{v: v, d: pn.currentDir}) {
			// fmt.Printf("loop detected at %v,direction: %v\n", v, pn.currentDir)
			return paths, true
		}
		traces.Add(traceItem{v: v, d: pn.currentDir})
		// fmt.Printf("traces: %v\n", traces)
		// fmt.Printf("next at %v: direction: %v\n", v, pn.currentDir)
		paths.Add(v)
	}
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("06.txt")
	pn := panel{
		m: grid.NewVecMatrix[byte](),
	}
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		for j, c := range line {
			if c == '^' {
				pn.startAt = grid.Vec{X: j, Y: i}
				pn.startAtDir = up
				pn.currentDir = up
			}
			if c == 'v' {
				pn.startAt = grid.Vec{X: j, Y: i}
				pn.currentDir = down
				pn.startAtDir = down
			}
			if c == '<' {
				pn.startAt = grid.Vec{X: j, Y: i}
				pn.currentDir = left
				pn.startAtDir = left
			}
			if c == '>' {
				pn.startAt = grid.Vec{X: j, Y: i}
				pn.startAtDir = right
				pn.currentDir = right
			}
			pn.m.Add(grid.Vec{X: j, Y: i}, byte(c))
		}
		return nil
	})
	paths, _ := findPaths(ctx, pn)
	// pn.m.ForEach(func(v grid.Vec, b byte) {
	// 	if b == '#' {
	// 		return
	// 	}
	// 	if paths.Has(v) {
	// 		pn.m[v] = 'X'
	// 	} else {
	// 		pn.m[v] = '.'
	// 	}
	// })
	// fmt.Println(pn)
	fmt.Printf("p1: %d\n", paths.Size())
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("06.txt")
	pn := panel{
		m: grid.NewVecMatrix[byte](),
	}
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		for j, c := range line {
			if c == '^' {
				pn.startAt = grid.Vec{X: j, Y: i}
				pn.startAtDir = up
				pn.currentDir = up
			}
			if c == 'v' {
				pn.startAt = grid.Vec{X: j, Y: i}
				pn.startAtDir = down
				pn.currentDir = down
			}
			if c == '<' {
				pn.startAt = grid.Vec{X: j, Y: i}
				pn.startAtDir = left
				pn.currentDir = left
			}
			if c == '>' {
				pn.startAt = grid.Vec{X: j, Y: i}
				pn.startAtDir = right
				pn.currentDir = right
			}
			pn.m.Add(grid.Vec{X: j, Y: i}, byte(c))
		}
		return nil
	})
	cp := pn.clone()
	paths, _ := findPaths(ctx, cp)
	var options int
	paths.Each(func(item grid.Vec) bool {
		cp := pn.clone()
		cp.m.Add(item, '#')
		_, isInLoop := findPaths(ctx, cp)
		if isInLoop {
			// fmt.Printf("col: %d,row: %d\n", item.X, item.Y)
			options += 1
		}
		return true
	})
	fmt.Printf("p2: %d\n", options)
}

func main() {
	ctx := context.TODO()
	p1(ctx)
	p2(ctx)
}
