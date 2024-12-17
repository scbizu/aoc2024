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
	loc   grid.Vec
	path  *set.Set[dv]
	d     direction
	score int
}

func copyMaze(m grid.VecMatrix[byte]) grid.VecMatrix[byte] {
	n := grid.NewVecMatrix[byte]()
	for k, v := range m {
		n.Add(k, v)
	}
	return n
}

func (m maze) walk(_ context.Context, min int) (int, *set.Set[grid.Vec]) {
	visited := make(map[dv]int)
	q := queue.NewQueue[state]()
	q.Push(state{
		loc:  m.start,
		path: set.New[dv](),
		d:    east,
	})
	fp := set.New[grid.Vec]()
	if min == 0 {
		min = math.MaxInt32
	}
	for q.Len() > 0 {
		n := q.Pop()
		// fmt.Printf("loc: %v,d: %v,visited: %v\n", n.loc, n.d, n.visited)
		if n.loc == m.end {
			if n.score <= min {
				// fmt.Printf("score: %d\n", n.score)
				n.path.Each(func(item dv) bool {
					fp.Add(item.v)
					return true
				})
				min = n.score
			}
			continue
		}
		c, ok := m.m.Get(n.loc)
		if !ok {
			continue
		}
		if c == '#' {
			continue
		}
		// go north
		{
			next := n.loc.Add(grid.Vec{X: 0, Y: -1})
			nc, ok := m.m.Get(next)
			if ok && (nc == '.' || nc == 'E') {
				isTurn := n.d == east || n.d == west
				score := n.score
				if isTurn {
					score += 1001
				} else {
					score++
				}
				if v, ok := visited[dv{
					v: n.loc,
					d: north,
				}]; (ok && v >= score) || !ok {
					p := n.path.Copy()
					p.Add(dv{
						v: n.loc,
						d: north,
					})
					visited[dv{
						v: n.loc,
						d: north,
					}] = score
					q.Push(state{
						loc:   next,
						path:  p,
						d:     north,
						score: score,
					})
				}
			}
		}
		// go east
		{
			next := n.loc.Add(grid.Vec{X: 1, Y: 0})
			nc, ok := m.m.Get(next)
			if ok && (nc == '.' || nc == 'E') {
				isTurn := n.d == north || n.d == south
				score := n.score
				if isTurn {
					score += 1001
				} else {
					score++
				}
				if v, ok := visited[dv{
					v: n.loc,
					d: east,
				}]; (ok && v >= score) || !ok {
					p := n.path.Copy()
					p.Add(dv{
						v: n.loc,
						d: east,
					})
					q.Push(state{
						loc:   next,
						path:  p,
						d:     east,
						score: score,
					})
					visited[dv{
						v: n.loc,
						d: east,
					}] = score
				}
			}
		}
		// go south
		{
			next := n.loc.Add(grid.Vec{X: 0, Y: 1})
			nc, ok := m.m.Get(next)
			if ok && (nc == '.' || nc == 'E') {
				isTurn := n.d == west || n.d == east
				score := n.score
				if isTurn {
					score += 1001
				} else {
					score++
				}
				if v, ok := visited[dv{
					v: n.loc,
					d: south,
				}]; (ok && v >= score) || !ok {
					p := n.path.Copy()
					p.Add(dv{
						v: n.loc,
						d: south,
					})
					q.Push(state{
						loc:   next,
						path:  p,
						d:     south,
						score: score,
					})
					visited[dv{
						v: n.loc,
						d: south,
					}] = score
				}
			}
		}
		// go west
		{
			next := n.loc.Add(grid.Vec{X: -1, Y: 0})
			nc, ok := m.m.Get(next)
			if ok && (nc == '.' || nc == 'E') {
				isTurn := n.d == north || n.d == south
				score := n.score
				if isTurn {
					score += 1001
				} else {
					score++
				}
				if v, ok := visited[dv{
					v: n.loc,
					d: west,
				}]; (ok && v >= score) || !ok {
					p := n.path.Copy()
					p.Add(dv{
						v: n.loc,
						d: west,
					})
					q.Push(state{
						loc:   next,
						path:  p,
						d:     west,
						score: score,
					})
					visited[dv{
						v: n.loc,
						d: west,
					}] = score
				}
			}
		}
	}
	return min, fp
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
	min, _ := m.walk(ctx, 0)
	fmt.Printf("p1: %d\n", min)
}

func p2(ctx context.Context) {
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
	min, _ := m.walk(ctx, 0)
	_, fp := m.walk(ctx, min)
	// fp.Each(func(item grid.Vec) bool {
	// 	m.m.Add(item, 'O')
	// 	return true
	// })
	// m.m.Print(os.Stdout, "%c")
	fmt.Printf("p2: %d\n", fp.Size()+1)
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
