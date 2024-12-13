package main

import (
	"context"
	"fmt"

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
	side          int
}

func (b block) String() string {
	return fmt.Sprintf("block: %d side: %d\n", len(b.plants), b.side)
}

func (blk *block) handleCorner(
	g garden,
	p grid.Vec,
	v1, v2 grid.Vec,
) {
	var x, y int
	if v1.X != 0 {
		x = v1.X
	} else {
		x = v2.X
	}
	if v1.Y != 0 {
		y = v1.Y
	} else {
		y = v2.Y
	}
	otherR := p.Add(grid.Vec{X: x, Y: y})
	v, ok := g.plants.Get(otherR)
	if !ok || v != blk.v {
		// println("+1")
		blk.side++
	}
}

// v1,v2 是否垂直
func isCross(v1, v2 grid.Vec) bool {
	return v1.X*v2.X+v1.Y*v2.Y == 0
}

func (g garden) Group2(ctx context.Context) []block {
	var blks []block
	g.plants.ForEach(func(v grid.Vec, b byte) {
		if _, ok := g.visited[v]; ok {
			return
		}
		q := queue.NewQueue[grid.Vec]()
		q.Push(v)
		g.visited[v] = struct{}{}
		blk := &block{
			plants: []grid.Vec{v},
			v:      b,
			side:   0,
		}

		for q.Len() > 0 {
			p := q.Pop()
			var nbs, others []grid.Vec
			for _, next := range g.plants.GetNeighbor(p) {
				n, ok := g.plants.Get(next)
				if !ok {
					others = append(others, next.Sub(p))
					continue
				}
				_, ok = g.visited[next]
				if n == b {
					if !ok {
						q.Push(next)
						blk.plants = append(blk.plants, next)
						g.visited[next] = struct{}{}
					}
					nbs = append(nbs, next.Sub(p))
				} else {
					others = append(others, next.Sub(p))
				}
			}
			// fmt.Printf("on: %v,nbs: %d,others: %v\n", p, len(nbs), others)
			for _, ov := range pick2(others) {
				if len(ov) <= 1 {
					continue
				}
				if isCross(ov[0], ov[1]) {
					// println("+1")
					blk.side++
				}
			}
			switch len(nbs) {
			case 1:
			// 存在一个相邻接点时，corner 不增加 , side 也不会增加
			case 2:
				// 分情况讨论
				// 1. 存在两个反向的相邻接点时，e.g.: O<-O->O , 此时 不存在 corner
				//                                O ?
				// 2. 存在两个垂直的相邻接点时，e.g.: O O ,此时:
				//    2.1 如果 ? 不是 O , 那么 corner + 1 (inner corner)
				//    2.2 如果 ? 是 O, 那么 corner 不增加
				// 左，上
				// 左，下
				// 上, 右
				// 下, 右
				blk.handleCorner(g, p, nbs[0], nbs[1])
			case 3:
				// 左，上，右
				// 左，上，下
				// 左，下，右
				// 上，下，右
				for _, vecs := range pick2(nbs) {
					blk.handleCorner(g, p, vecs[0], vecs[1])
				}
			case 4:
				// 左，上，右，下
				for _, vecs := range pick2(nbs) {
					blk.handleCorner(g, p, vecs[0], vecs[1])
				}
			}
		}
		blks = append(blks, *blk)
	})
	return blks
}

func pick2(vecs []grid.Vec) [][]grid.Vec {
	if len(vecs) <= 2 {
		return [][]grid.Vec{vecs}
	}
	var group [][]grid.Vec
	for i := 0; i < len(vecs); i++ {
		for j := i + 1; j < len(vecs); j++ {
			group = append(group, []grid.Vec{vecs[i], vecs[j]})
		}
	}
	return group
}

func p1(ctx context.Context) {
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("12.txt")
	g := garden{
		plants:  grid.NewVecMatrix[byte](),
		visited: make(map[grid.Vec]struct{}),
	}
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		for j, c := range line {
			g.plants.Add(grid.Vec{X: j, Y: i}, byte(c))
		}
		return nil
	})
	blks := g.Group2(ctx)
	var totalSide int
	for _, blk := range blks {
		// fmt.Println(blk)
		totalSide += blk.side * len(blk.plants)
	}
	fmt.Printf("p2: %d\n", totalSide)
}

func main() {
	ctx := context.Background()
	// p1(ctx)
	p2(ctx)
}
