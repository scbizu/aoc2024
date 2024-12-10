package main

import (
	"context"
	"fmt"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/queue"
)

type maze struct {
	g grid.VecMatrix[uint8]
}

func (m maze) walk() int {
	var total int
	m.g.ForEach(func(v grid.Vec, u uint8) {
		if u != '9' {
			return
		}
		q := queue.NewQueue[grid.Vec]()
		q.Push(v)
		fmt.Printf("current: %v\n", v)
		for q.Len() > 0 {
			n := q.Pop()
			if m.g[n] == '0' {
				total++
			}
			for _, next := range m.g.GetNeighbor(n) {
				if m.g[next] == m.g[n]-1 {
					fmt.Printf("next: %v\n", next)
					q.Push(next)
				}
			}
		}
	})
	return total
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("10.txt")
	g := grid.NewVecMatrix[uint8]()
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		for j, c := range line {
			g.Add(grid.Vec{
				X: j,
				Y: i,
			}, byte(c))
		}
		return nil
	})
	m := maze{g}
	fmt.Printf("p1: %d\n", m.walk())
}

func main() {
	ctx := context.Background()
	p1(ctx)
}
