package main

import (
	"context"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/queue"
)

type garden struct {
	plants  grid.VecMatrix[byte]
	visited map[grid.Vec]struct{}
}

type block struct {
	plants []grid.Vec
	v      byte
	// 内接
	insideLineaer int
}

func (g garden) Group(ctx context.Context) []block {
	g.plants.ForEach(func(v grid.Vec, b byte) {
		if _, ok := g.visited[v]; ok {
			return
		}
		var linear int
		q := queue.NewQueue[grid.Vec]()
		q.Push(v)
		for q.Len() > 0 {
			p := q.Pop()
			for _, next := range g.plants.GetNeighbor(p) {
				if _, ok := g.visited[next]; ok {
					continue
				}
				if g.plants[next] == b {
					q.Push(next)
					linear++
				}
				g.visited[next] = struct{}{}
			}
		}
	})
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("12.txt")
	g := garden{
		plants: grid.NewVecMatrix[byte](),
	}
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		for j, c := range line {
			g.plants.Add(grid.Vec{X: j, Y: i}, byte(c))
		}
		return nil
	})
}

func main() {
	ctx := context.Background()
	p1(ctx)
}
