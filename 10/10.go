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
		if u != '0' {
			return
		}
		v9 := make(map[grid.Vec]struct{})
		q := queue.NewQueue[grid.Vec]()
		q.Push(v)
		for q.Len() > 0 {
			n := q.Pop()
			if m.g[n] == '9' {
				if _, ok := v9[n]; !ok {
					v9[n] = struct{}{}
					total++
				}
			}
			for _, next := range m.g.GetNeighbor(n) {
				if m.g[next] == m.g[n]+1 {
					q.Push(next)
				}
			}
		}
	})
	return total
}

type traceItem struct {
	head   grid.Vec
	trails []grid.Vec
}

func (m maze) walk2() int {
	var total int
	m.g.ForEach(func(v grid.Vec, u uint8) {
		if u != '0' {
			return
		}
		q := queue.NewQueue[*traceItem]()
		q.Push(&traceItem{
			head:   v,
			trails: []grid.Vec{v},
		})
		for q.Len() > 0 {
			n := q.Pop()
			if m.g[n.head] == '9' {
				if len(n.trails) == 10 {
					total++
				}
			}
			for _, next := range m.g.GetNeighbor(n.head) {
				if m.g[next] == m.g[n.head]+1 {
					q.Push(&traceItem{
						head:   next,
						trails: append(n.trails, next),
					})
				}
			}
		}
	})
	return total
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("10.txt")
	g := grid.NewVecMatrix[byte]()
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

func p2(ctx context.Context) {
	txt := input.NewTXTFile("10.txt")
	g := grid.NewVecMatrix[byte]()
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
	fmt.Printf("p2: %d\n", m.walk2())
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
