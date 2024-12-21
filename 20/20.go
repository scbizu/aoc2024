package main

import (
	"context"
	"fmt"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/math"
)

type raceTrack struct {
	m grid.VecMatrix[byte]

	start, end grid.Vec
	cache      map[grid.Vec]int
}

func manhattan(v1, v2 grid.Vec) int {
	return math.Abs(v1.X-v2.X) + math.Abs(v1.Y-v2.Y)
}

func (r *raceTrack) wallBreak(_ context.Context, atMost int, leastPath int) int {
	var total int
	// for every vec in cache, we should search max to `atMost` area to find if we can use our cheat time
	for vec, d := range r.cache {
		for j := vec.X - atMost; j <= vec.X+atMost; j++ {
			for k := vec.Y - atMost; k <= vec.Y+atMost; k++ {
				if p, ok := r.cache[grid.Vec{X: j, Y: k}]; ok {
					// 「Cheats don't need to use all 20 picoseconds」
					if manhattan(vec, grid.Vec{X: j, Y: k}) > atMost {
						continue
					}
					if p-d-manhattan(vec, grid.Vec{X: j, Y: k}) >= leastPath {
						total++
					}
				}
			}
		}
	}
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
	delta := rs.wallBreak(ctx, 2, least)
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
	// fmt.Printf("cache: %v\n", rs.cache)
	least := 100
	maxBreak := 20
	delta := rs.wallBreak(ctx, maxBreak, least)
	fmt.Printf("p2: %d\n", delta)
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
