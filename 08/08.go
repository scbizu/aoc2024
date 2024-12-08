package main

import (
	"context"
	"fmt"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/set"
)

// line: y = kx + b
type line struct {
	k       int
	b       int
	origins [2]grid.Vec
}

func (l line) String() string {
	return fmt.Sprintf("y = %dx + %d", l.k, l.b)
}

func (l line) IsPerfect(v grid.Vec) bool {
	return l.k*v.X+l.b == v.Y && (grid.Distance(l.origins[0], v) == 2*grid.Distance(l.origins[1], v) ||
		grid.Distance(l.origins[1], v) == 2*grid.Distance(l.origins[0], v))
}

// 两个向量确定一条直线
func NewLine(v1, v2 grid.Vec) line {
	k := (v2.Y - v1.Y) / (v2.X - v1.X)
	b := v1.Y - k*v1.X
	return line{
		k:       k,
		b:       b,
		origins: [2]grid.Vec{v1, v2},
	}
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("08.txt")
	g := grid.NewVecMatrix[byte]()
	cmap := make(map[byte][]grid.Vec)
	originSet := set.New[grid.Vec]()
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
				originSet.Add(grid.Vec{
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
				lines[c] = append(lines[c], NewLine(vecs[i], vecs[j]))
			}
		}
	}
	antiNodes := set.New[grid.Vec]()
	for c, ls := range lines {
		for _, l := range ls {
			fmt.Printf("current on antennas: %s,line: %s\n", string(c), l)
			g.ForEach(func(v grid.Vec, b byte) {
				if l.IsPerfect(v) {
					if !originSet.Has(v) {
						fmt.Printf("found: antiNode: %v\n", v)
						antiNodes.Add(v)
					}
				}
			})
		}
	}

	fmt.Printf("p1: %d\n", antiNodes.Size())
}

func main() {
	ctx := context.TODO()
	p1(ctx)
}
