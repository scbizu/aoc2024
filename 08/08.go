package main

import (
	"context"
	"fmt"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/set"
)

type line struct {
	origins [2]grid.Vec
}

func NewLine(v1, v2 grid.Vec) line {
	return line{
		origins: [2]grid.Vec{v1, v2},
	}
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("08.txt")
	g := grid.NewVecMatrix[byte]()
	cmap := make(map[byte][]grid.Vec)
	uniques := set.New[grid.Vec]()
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		for j, c := range line {
			g.Add(grid.Vec{
				X: j,
				Y: i,
			}, byte(c))
			if c != '.' {
				cmap[byte(c)] = append(cmap[byte(c)], grid.Vec{
					X: j,
					Y: i,
				})
				uniques.Add(grid.Vec{
					X: j,
					Y: i,
				})
			}
		}
		return nil
	})

	lines := make(map[byte][]line)

	for c, vecs := range cmap {
		for i := 0; i < len(vecs); i++ {
			for j := i + 1; j < len(vecs); j++ {
				// fmt.Printf("vecs: %v,%v\n", vecs[i], vecs[j])
				ln := NewLine(vecs[i], vecs[j])
				// fmt.Printf("current on antennas: %s,line: %s\n", string(c), ln)
				lines[c] = append(lines[c], ln)
			}
		}
	}
	var total int
	for _, ls := range lines {
		for _, l := range ls {
			// fmt.Printf("current on antennas: %s,line: %s\n", string(c), l)
			h, w := l.origins[1].Y-l.origins[0].Y, l.origins[1].X-l.origins[0].X
			node1 := grid.Vec{
				X: l.origins[0].X - w,
				Y: l.origins[0].Y - h,
			}
			if _, ok := g.Get(node1); ok && !uniques.Has(node1) {
				// fmt.Printf("node1: %v\n", node1)
				total++
			}
			node2 := grid.Vec{
				X: l.origins[1].X + w,
				Y: l.origins[1].Y + h,
			}
			if _, ok := g.Get(node2); ok && !uniques.Has(node2) {
				// fmt.Printf("node2: %v\n", node2)
				total++
			}
		}
	}

	fmt.Printf("p1: %d\n", total)
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("08.txt")
	g := grid.NewVecMatrix[byte]()
	cmap := make(map[byte][]grid.Vec)
	uniques := set.New[grid.Vec]()
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		for j, c := range line {
			g.Add(grid.Vec{
				X: j,
				Y: i,
			}, byte(c))
			if c != '.' {
				cmap[byte(c)] = append(cmap[byte(c)], grid.Vec{
					X: j,
					Y: i,
				})
				uniques.Add(grid.Vec{
					X: j,
					Y: i,
				})
			}
		}
		return nil
	})

	lines := make(map[byte][]line)

	for c, vecs := range cmap {
		for i := 0; i < len(vecs); i++ {
			for j := i + 1; j < len(vecs); j++ {
				// fmt.Printf("vecs: %v,%v\n", vecs[i], vecs[j])
				ln := NewLine(vecs[i], vecs[j])
				// fmt.Printf("current on antennas: %s,line: %s\n", string(c), ln)
				lines[c] = append(lines[c], ln)
			}
		}
	}
	for _, ls := range lines {
		for _, l := range ls {
			// fmt.Printf("current on antennas: %s,line: %s\n", string(c), l)
			h, w := l.origins[1].Y-l.origins[0].Y, l.origins[1].X-l.origins[0].X
			n1 := grid.Vec{
				X: l.origins[0].X - w,
				Y: l.origins[0].Y - h,
			}
			for {
				if _, ok := g.Get(n1); !ok {
					break
				}
				if !uniques.Has(n1) {
					// fmt.Printf("n1: %v\n", n1)
					uniques.Add(n1)
				}
				n1 = grid.Vec{
					X: n1.X - w,
					Y: n1.Y - h,
				}
			}
			n2 := grid.Vec{
				X: l.origins[1].X + w,
				Y: l.origins[1].Y + h,
			}
			for {
				if _, ok := g.Get(n2); !ok {
					break
				}
				if !uniques.Has(n2) {
					// fmt.Printf("n2: %v\n", n2)
					uniques.Add(n2)
				}
				n2 = grid.Vec{
					X: n2.X + w,
					Y: n2.Y + h,
				}
			}
		}
	}

	fmt.Printf("p2: %d\n", uniques.Size())
}

func main() {
	ctx := context.TODO()
	p1(ctx)
	p2(ctx)
}
