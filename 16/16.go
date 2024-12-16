package main

import (
	"context"
	"fmt"
	"math"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/queue"
	"github.com/magejiCoder/magejiAoc/set"
)

type direction int

const (
	north direction = iota
	east
	south
	west
)

type maze struct {
	m          grid.VecMatrix[byte]
	start, end grid.Vec
}

type dv struct {
	v grid.Vec
	d direction
}

type state struct {
	loc  grid.Vec
	path *set.Set[dv]
	d    direction
}

func score(d direction, p *set.Set[dv]) int {
	var score int
	p.Each(func(item dv) bool {
		// 转向
		if d%2 != item.d%2 {
			score += 1000
			d = item.d
		} else {
			score++
		}
		return true
	})
	return score
}

func (m maze) walk(_ context.Context) int {
	q := queue.NewQueue[state]()
	q.Push(state{
		loc:  m.start,
		path: set.New[dv](),
		d:    east,
	})
	min := math.MaxInt32
	for q.Len() > 0 {
		n := q.Pop()
		if n.loc == m.end {
			fmt.Printf("path: %v\n", n.path)
			if score(east, n.path) < min {
				min = score(east, n.path)
			}
		}
		c, ok := m.m.Get(n.loc)
		if !ok {
			continue
		}
		if c == '#' {
			continue
		}
		p := n.path.Copy()
		// go north
		{
			next := n.loc.Add(grid.Vec{X: 0, Y: -1})
			if n.path.Has(dv{
				v: next,
				d: north,
			}) {
				continue
			}
			nc, ok := m.m.Get(next)
			if ok && (nc == '.' || nc == 'E') {
				p.Add(dv{
					v: next,
					d: north,
				})
				q.Push(state{
					loc:  next,
					path: p,
					d:    north,
				})
			}
		}
		// go east
		{
			next := n.loc.Add(grid.Vec{X: 1, Y: 0})
			if n.path.Has(dv{
				v: next,
				d: east,
			}) {
				continue
			}
			nc, ok := m.m.Get(next)
			if ok && (nc == '.' || nc == 'E') {
				p.Add(dv{
					v: next,
					d: east,
				})
				q.Push(state{
					loc:  next,
					path: p,
					d:    east,
				})
			}
		}
		// go south
		{
			next := n.loc.Add(grid.Vec{X: 0, Y: 1})
			if n.path.Has(dv{
				v: next,
				d: south,
			}) {
				continue
			}
			nc, ok := m.m.Get(next)
			if ok && (nc == '.' || nc == 'E') {
				p.Add(dv{
					v: next,
					d: south,
				})
				q.Push(state{
					loc:  next,
					path: p,
					d:    south,
				})
			}
		}
		// go west
		{
			next := n.loc.Add(grid.Vec{X: -1, Y: 0})
			if n.path.Has(dv{
				v: next,
				d: west,
			}) {
				continue
			}
			nc, ok := m.m.Get(next)
			if ok && (nc == '.' || nc == 'E') {
				p.Add(dv{
					v: next,
					d: west,
				})
				q.Push(state{
					loc:  next,
					path: p,
					d:    west,
				})
			}
		}
	}
	return min
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("16.txt")
	m := maze{
		m: grid.NewVecMatrix[byte](),
	}
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		for j, c := range line {
			if c == 'S' {
				m.start = grid.Vec{X: j, Y: i}
			}
			if c == 'E' {
				m.end = grid.Vec{X: j, Y: i}
			}
			m.m.Add(grid.Vec{X: j, Y: i}, byte(c))
		}
		return nil
	})
	min := m.walk(ctx)
	fmt.Printf("p1: %d\n", min)
}

func main() {
	ctx := context.Background()
	p1(ctx)
}
