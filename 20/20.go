package main

import (
	"context"
	"fmt"
	"slices"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/queue"
	"github.com/magejiCoder/magejiAoc/set"
)

type direction int

const (
	up direction = iota
	down
	left
	right
)

type dirVec struct {
	v grid.Vec
	d direction
}

type raceTrack struct {
	m grid.VecMatrix[byte]

	start, end grid.Vec
	wall       grid.Vec
}

type state struct {
	v    grid.Vec
	d    direction
	path *set.Set[grid.Vec]
}

func copyMaze(m grid.VecMatrix[byte]) grid.VecMatrix[byte] {
	n := make(grid.VecMatrix[byte], len(m))
	for k, v := range m {
		n[k] = v
	}
	return n
}

func (r raceTrack) clone() raceTrack {
	return raceTrack{
		m:     copyMaze(r.m),
		start: r.start,
		end:   r.end,
	}
}

func (r *raceTrack) removeInnerWall(ctx context.Context, base int) map[int]int {
	// find wall
	innerWalls := set.New[grid.Vec]()
	r.m.ForEach(func(v grid.Vec, b byte) {
		if v.X == 0 || v.Y == 0 || v.X == r.m.Cols()-1 || v.Y == r.m.Rows()-1 || v == r.wall {
			return
		}
		if b == '#' {
			innerWalls.Add(v)
		}
	})
	paths := make(map[int]int)
	innerWalls.Each(func(item grid.Vec) bool {
		r2 := r.clone()
		r2.m.Add(item, '.')
		_, p := r2.walk(ctx)
		fmt.Printf("rm wall: %v,cost: %d\n", item, p)
		if p != -1 && base-p > 0 {
			paths[base-p]++
		}
		return true
	})
	return paths
}

func (r raceTrack) walk2(ctx context.Context) int {
	visited := make(map[dirVec]int)
	dirs := []direction{up, down, left, right}
	var dfs func(v grid.Vec, d direction, path *set.Set[grid.Vec]) int
	dfs = func(v grid.Vec, d direction, path *set.Set[grid.Vec]) int {
		if v == r.end {
			return path.Size()
		}
		var nextV grid.Vec
		switch d {
		case up:
			nextV = v.Add(grid.Vec{X: 0, Y: -1})
		case down:
			nextV = v.Add(grid.Vec{X: 0, Y: 1})
		case left:
			nextV = v.Add(grid.Vec{X: -1, Y: 0})
		case right:
			nextV = v.Add(grid.Vec{X: 1, Y: 0})
		}
		c, ok := r.m.Get(nextV)
		if !ok {
			// wall
			return -1
		}
		if c == '#' {
			// wall
			return -1
		}
		// fmt.Printf("v: %v,d: %d,path: %v,trace: %d\n", nextV, d, path.List(), path.Size())
		if v, ok := visited[dirVec{
			v: nextV,
			d: d,
		}]; ok && v <= path.Size()+1 {
			return -1
		}
		visited[dirVec{
			v: nextV,
			d: d,
		}] = path.Size() + 1
		// fmt.Printf("visited: %v\n", visited)
		p := path.Copy()
		p.Add(nextV)
		var totals []int
		fmt.Printf("visit: %v, trace: %v\n", nextV, p)
		for _, d := range dirs {
			t := dfs(nextV, d, p)
			if t != -1 {
				totals = append(totals, t)
			}
		}
		fmt.Printf("totals: %v\n", totals)
		if len(totals) == 0 {
			return -1
		}
		slices.Sort(totals)
		return totals[0]
	}
	var t int
	for _, d := range dirs {
		visited[dirVec{
			v: r.start,
			d: d,
		}] = 0
		t += dfs(r.start, d, set.New[grid.Vec]())
	}
	return t
}

func (r raceTrack) walk(_ context.Context) (grid.Vec, int) {
	// visited := make(map[grid.Vec]int)
	q := queue.NewQueue[state]()
	dirs := []direction{up, down, left, right}
	var lastV grid.Vec
	for _, d := range dirs {
		q.Push(state{
			v:    r.start,
			d:    d,
			path: set.New[grid.Vec](),
		})
		// visited[r.start] = 0
	}
	var totals []int
	for q.Len() > 0 {
		st := q.Pop()
		if st.v == r.end {
			totals = append(totals, st.path.Size())
			lastV = st.path.List()[st.path.Size()-1]
			break
		}
		var nextV grid.Vec
		switch st.d {
		case up:
			nextV = st.v.Add(grid.Vec{X: 0, Y: -1})
		case down:
			nextV = st.v.Add(grid.Vec{X: 0, Y: 1})
		case left:
			nextV = st.v.Add(grid.Vec{X: -1, Y: 0})
		case right:
			nextV = st.v.Add(grid.Vec{X: 1, Y: 0})
		}
		c, ok := r.m.Get(nextV)
		if !ok {
			// wall
			continue
		}
		if c == '#' {
			// wall
			continue
		}
		if st.path.Has(nextV) {
			continue
		}
		// fmt.Printf("next: %v\n", nextV)
		for _, d := range dirs {
			// vc, ok := visited[nextV]
			// if ok && vc <= st.path.Size()+1 {
			// 	// visited
			// 	continue
			// }
			// visited[nextV] = st.path.Size() + 1
			p := st.path.Copy()
			p.Add(st.v)
			// fmt.Printf("push: %v,path: %v\n", nextV, p.List())
			q.Push(state{
				v:    nextV,
				d:    d,
				path: p,
			})
		}

	}
	if len(totals) == 0 {
		return grid.Vec{}, -1
	}
	return lastV, totals[0]
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("20.txt")
	rs := &raceTrack{
		m: grid.NewVecMatrix[byte](),
	}
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		for j, c := range line {
			if c == 'S' {
				rs.start = grid.Vec{
					X: j,
					Y: i,
				}
			}
			if c == 'E' {
				rs.end = grid.Vec{
					X: j,
					Y: i,
				}
			}
			rs.m.Add(grid.Vec{
				X: j,
				Y: i,
			}, byte(c))
		}
		return nil
	})
	last, base := rs.walk(ctx)
	fmt.Printf("base: %d\n", base)
	rs.wall = last
	deltas := rs.removeInnerWall(ctx, base)
	least := 100
	var total int
	for k, v := range deltas {
		if k >= least {
			total += v
		}
	}
	fmt.Printf("p1: %d\n", total)
}

func main() {
	ctx := context.Background()
	p1(ctx)
}
