package main

import (
	"context"
	"fmt"
	"slices"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/queue"
)

type raceTrack struct {
	m grid.VecMatrix[byte]

	start, end grid.Vec
	cache      map[grid.Vec]int
}

type pair struct {
	min int
	max int
}

type gridPath struct {
	v    grid.Vec
	path int
}

func (r *raceTrack) wallBreak(_ context.Context, atMost int, leastPath int) int {
	// for every wall. we can find the break path until it reaches the path or atMost
	prs := make(map[pair]int)
	r.m.ForEach(func(v grid.Vec, b byte) {
		if b != '#' {
			return
		}
		var nums []gridPath
		visited := make(map[grid.Vec]struct{})
		q := queue.NewQueue[gridPath]()
		q.Push(gridPath{
			v:    v,
			path: 1,
		})
		for q.Len() > 0 {
			cur := q.Pop()
			if _, ok := visited[cur.v]; ok {
				continue
			}
			visited[cur.v] = struct{}{}
			if len(visited) > atMost {
				break
			}
			for _, n := range r.m.GetNeighbor(cur.v) {
				_, ok := r.cache[n]
				// find the entrance out
				if ok {
					nums = append(nums, gridPath{
						v:    n,
						path: cur.path,
					})
					continue
				}
				nc, ok := r.m.Get(n)
				if ok && nc == '#' {
					q.Push(gridPath{
						v:    n,
						path: cur.path + 1,
					})
				}
			}
		}
	})
	return total
}

func (r *raceTrack) leastPath(_ context.Context, least int) int {
	var total int
	// find wall
	r.m.ForEach(func(v grid.Vec, b byte) {
		if b != '#' {
			return
		}
		var nums []int
		for _, n := range r.m.GetNeighbor(v) {
			if c, ok := r.cache[n]; ok {
				nums = append(nums, c)
			}
		}
		if len(nums) > 1 {
			slices.Sort(nums)
			if nums[len(nums)-1]-nums[0]-2 >= least {
				// fmt.Printf("v: %v, nums: %v,saves: %d\n", v, nums, nums[len(nums)-1]-nums[0]-2)
				total++
			}
		}
	})
	return total
}

func (r *raceTrack) walk(_ context.Context) {
	var dfs func(v grid.Vec, d int, path map[grid.Vec]int)
	dfs = func(v grid.Vec, d int, path map[grid.Vec]int) {
		if v == r.end {
			r.cache[v] = d
			return
		}
		for _, n := range r.m.GetNeighbor(v) {
			if _, ok := path[n]; ok {
				continue
			}
			v, ok := r.m.Get(n)
			if !ok {
				continue
			}
			if v == '.' || v == 'E' {
				path[n] = d + 1
				r.cache[n] = d + 1
				dfs(n, d+1, path)
				break
			}
		}
	}
	p := map[grid.Vec]int{
		r.start: 0,
	}
	dfs(r.start, 0, p)
	// fmt.Printf("total: %d\n", len(p))
	for k, v := range r.cache {
		r.cache[k] = len(p) - v
	}
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("20.txt")
	rs := &raceTrack{
		m:     grid.NewVecMatrix[byte](),
		cache: make(map[grid.Vec]int),
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
	rs.cache[rs.start] = 0
	rs.walk(ctx)
	// fmt.Printf("cache: %d\n", len(rs.cache))
	least := 100
	delta := rs.leastPath(ctx, least)
	fmt.Printf("p1: %d\n", delta)
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("20.txt")
	rs := &raceTrack{
		m:     grid.NewVecMatrix[byte](),
		cache: make(map[grid.Vec]int),
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
	rs.cache[rs.start] = 0
	rs.walk(ctx)
	fmt.Printf("cache: %d\n", len(rs.cache))
	least := 76
	maxBreak := 20
	delta := rs.wallBreak(ctx, maxBreak, least)
	fmt.Printf("p2: %d\n", delta)
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
