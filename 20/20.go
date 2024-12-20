package main

import (
	"context"
	"fmt"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/math"
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
	lastTrace  *set.Set[grid.Vec]
	cache      map[grid.Vec]int
}

type state struct {
	v          grid.Vec
	path       *set.Set[grid.Vec]
	foundOther bool
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
		m:         copyMaze(r.m),
		start:     r.start,
		end:       r.end,
		lastTrace: r.lastTrace,
	}
}

func isNear(v1, v2 grid.Vec) bool {
	v := v1.Sub(v2)
	return math.Abs(v.X)+math.Abs(v.Y) == 1
}

func (r *raceTrack) removeInnerWall(ctx context.Context) map[int]int {
	base := r.lastTrace.Size()
	// find wall
	innerWalls := set.New[grid.Vec]()
	r.m.ForEach(func(v grid.Vec, b byte) {
		if b != '#' {
			return
		}
		r.lastTrace.Each(func(item grid.Vec) bool {
			if isNear(v, item) {
				fmt.Printf("item,v: %v,%v\n", item, v)
				innerWalls.Add(v)
			}
			return true
		})
	})
	fmt.Printf("#: %d\n", innerWalls.Size())
	paths := make(map[int]int)
	innerWalls.Each(func(item grid.Vec) bool {
		r2 := r.clone()
		r2.m.Add(item, '.')
		p := r2.walk(ctx)
		// fmt.Printf("rm wall: %v,cost: %d\n", item, p.Size())
		if p.Size() > 0 && base-p.Size() > 0 {
			paths[base-p.Size()]++
		}
		return true
	})
	return paths
}

//	func (r raceTrack) walk2(ctx context.Context) *set.Set[grid.Vec] {
//		dirs := []direction{up, down, left, right}
//		var dfs func(v grid.Vec, d direction, path *set.Set[grid.Vec])
//		dfs = func(v grid.Vec, d direction, path *set.Set[grid.Vec])  {
//			if v == r.end {
//				return
//			}
//			var nextV grid.Vec
//			switch d {
//			case up:
//				nextV = v.Add(grid.Vec{X: 0, Y: -1})
//			case down:
//				nextV = v.Add(grid.Vec{X: 0, Y: 1})
//			case left:
//				nextV = v.Add(grid.Vec{X: -1, Y: 0})
//			case right:
//				nextV = v.Add(grid.Vec{X: 1, Y: 0})
//			}
//			if path.Has(nextV) {
//				return
//			}
//			c, ok := r.m.Get(nextV)
//			if !ok {
//				// wall
//				return
//			}
//			if c == '#' {
//				// wall
//				return
//			}
//			// fmt.Printf("visited: %v\n", visited)
//			p := path.Copy()
//			p.Add(nextV)
//			var totals []int
//			for _, d := range dirs {
//				t := dfs(nextV, d, p)
//				if t != -1 {
//					totals = append(totals, t)
//					break
//				}
//			}
//			fmt.Printf("totals: %v\n", totals)
//			if len(totals) == 0 {
//				return -1
//			}
//			slices.Sort(totals)
//			return totals[0]
//		}
//		for _, d := range dirs {
//			next := dfs(r.start, d, set.New[grid.Vec]())
//			if next != -1 {
//				return next
//			}
//		}
//		return -1
//	}
func addVec(v grid.Vec, d direction) grid.Vec {
	switch d {
	case up:
		return v.Add(grid.Vec{X: 0, Y: -1})
	case down:
		return v.Add(grid.Vec{X: 0, Y: 1})
	case left:
		return v.Add(grid.Vec{X: -1, Y: 0})
	case right:
		return v.Add(grid.Vec{X: 1, Y: 0})
	}
	panic("invalid direction")
}

func (r raceTrack) walk(_ context.Context) *set.Set[grid.Vec] {
	q := queue.NewQueue[state]()
	dirs := []direction{up, down, left, right}
	q.Push(state{
		v: r.start,
		path: set.New[grid.Vec](
			r.start,
		),
	})
	for q.Len() > 0 {
		st := q.Pop()
		if st.v == r.end {
			st.path.Remove(r.end)
			return st.path
		}
		var validVecs []grid.Vec
		for _, d := range dirs {
			vc := addVec(st.v, d)
			c, ok := r.m.Get(vc)
			if !ok {
				// wall
				continue
			}
			if c == '#' {
				// wall
				continue
			}
			if st.path.Has(vc) {
				continue
			}
			validVecs = append(validVecs, vc)
		}
		switch len(validVecs) {
		case 1:
			path := st.path.Copy()
			path.Add(validVecs[0])
			q.Push(state{
				v:          validVecs[0],
				path:       path,
				foundOther: st.foundOther,
			})
		case 2:
			for _, vv := range validVecs {
				if r.lastTrace.Has(vv) && !st.foundOther {
					continue
				}
				// choose the other way
				path := st.path.Copy()
				path.Add(vv)
				q.Push(state{
					v:          vv,
					path:       path,
					foundOther: true,
				})
			}
		case 3, 4:
			panic("bad route")
		}
	}
	return set.New[grid.Vec]()
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("20.txt")
	rs := &raceTrack{
		m:         grid.NewVecMatrix[byte](),
		lastTrace: set.New[grid.Vec](),
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
	base := rs.walk(ctx)
	fmt.Printf("base: %d\n", base.Size())
	rs.lastTrace = base
	deltas := rs.removeInnerWall(ctx)
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
