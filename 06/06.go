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
	currentAt  grid.Vec
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

func p1(ctx context.Context) {
	txt := input.NewTXTFile("06.txt")
	pn := panel{
		m: grid.NewVecMatrix[byte](),
	}
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		for j, c := range line {
			if c == '^' {
				pn.currentAt = grid.Vec{X: j, Y: i}
				pn.currentDir = up
			}
			if c == 'v' {
				pn.currentAt = grid.Vec{X: j, Y: i}
				pn.currentDir = down
			}
			if c == '<' {
				pn.currentAt = grid.Vec{X: j, Y: i}
				pn.currentDir = left
			}
			if c == '>' {
				pn.currentAt = grid.Vec{X: j, Y: i}
				pn.currentDir = right
			}
			pn.m.Add(grid.Vec{X: j, Y: i}, byte(c))
		}
		return nil
	})
	v := pn.currentAt
	next := pn.currentAt
	// fmt.Printf("start at %v: direction: %v\n", v, pn.currentDir)
	paths := set.New[grid.Vec](
		pn.currentAt,
	)
	for {
		if _, ok := pn.m.Get(next); !ok {
			break
		}
		switch pn.currentDir {
		case up:
			next = v.Add(grid.Vec{X: 0, Y: -1})
			nextV, ok := pn.m.Get(next)
			if !ok {
				break
			}
			if ok && nextV == '#' {
				// turn right
				next = v.Add(grid.Vec{X: 1, Y: 0})
				pn.currentDir = right
			}
			v = next
		case down:
			next = v.Add(grid.Vec{X: 0, Y: 1})
			nextV, ok := pn.m.Get(next)
			if !ok {
				break
			}
			if ok && nextV == '#' {
				// turn left
				next = v.Add(grid.Vec{X: -1, Y: 0})
				pn.currentDir = left
			}
			v = next
		case left:
			next = v.Add(grid.Vec{X: -1, Y: 0})
			nextV, ok := pn.m.Get(next)
			if !ok {
				break
			}
			if ok && nextV == '#' {
				// turn up
				next = v.Add(grid.Vec{X: 0, Y: -1})
				pn.currentDir = up
			}
			v = next
		case right:
			next = v.Add(grid.Vec{X: 1, Y: 0})
			nextV, ok := pn.m.Get(next)
			if !ok {
				break
			}
			if ok && nextV == '#' {
				// turn down
				next = v.Add(grid.Vec{X: 0, Y: 1})
				pn.currentDir = down
			}
			v = next
		}
		// fmt.Printf("next at %v: direction: %v\n", v, pn.currentDir)
		pn.currentAt = v
		if !paths.Has(v) {
			paths.Add(v)
		}
	}
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

func main() {
	ctx := context.TODO()
	p1(ctx)
}
