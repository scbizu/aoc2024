package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/queue"
	"github.com/magejiCoder/magejiAoc/set"
)

type memorySpace struct {
	m             grid.VecMatrix[byte]
	start         grid.Vec
	end           grid.Vec
	width, height int
}

type byteCord struct {
	v grid.Vec
}

func parseByteCord(s string) byteCord {
	parts := strings.Split(s, ",")
	return byteCord{
		v: grid.Vec{
			X: input.Atoi(parts[0]),
			Y: input.Atoi(parts[1]),
		},
	}
}

func (m *memorySpace) dropBytes(bc []byteCord) {
	for _, b := range bc {
		m.m.Add(b.v, '#')
	}
}

type state struct {
	v    grid.Vec
	d    direction
	path *set.Set[grid.Vec]
}

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

func (m *memorySpace) escape(_ context.Context) int {
	visited := make(map[dirVec]int)
	q := queue.NewQueue[state]()
	q.Push(state{
		v:    m.start,
		d:    down,
		path: set.New[grid.Vec](),
	})
	q.Push(state{
		v:    m.start,
		d:    right,
		path: set.New[grid.Vec](),
	})
	visited[dirVec{
		v: m.start,
		d: down,
	}] = 0
	visited[dirVec{
		v: m.start,
		d: right,
	}] = 0
	var totals []int
	for q.Len() > 0 {
		st := q.Pop()
		// fmt.Printf("v: %v,d: %v\n", st.v, st.d)
		if st.v == m.end {
			// escaped
			totals = append(totals, st.path.Size())
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
		c, ok := m.m.Get(nextV)
		if !ok {
			// wall
			continue
		}
		if c == '#' {
			// dropped byte
			continue
		}
		dirs := []direction{up, down, left, right}
		for _, d := range dirs {
			vc, ok := visited[dirVec{
				v: nextV,
				d: d,
			}]
			if ok && vc <= st.path.Size()+1 {
				// visited
				continue
			}
			visited[dirVec{
				v: nextV,
				d: d,
			}] = st.path.Size() + 1
			p := st.path.Copy()
			p.Add(st.v)
			q.Push(state{
				v:    nextV,
				d:    d,
				path: p,
			})
		}

	}
	slices.Sort(totals)
	return totals[0]
}

func main() {
	ctx := context.Background()
	p1(ctx)
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("18.txt")
	var bs []byteCord
	txt.ReadByLineEx(ctx, func(_ int, line string) error {
		bc := parseByteCord(line)
		bs = append(bs, bc)
		return nil
	})
	width, height := 71, 71
	ms := &memorySpace{
		m: grid.NewVecMatrix[byte](),
		start: grid.Vec{
			X: 0,
			Y: 0,
		},
		end: grid.Vec{
			X: width - 1,
			Y: height - 1,
		},
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			ms.m.Add(grid.Vec{
				X: x,
				Y: y,
			}, '.')
		}
	}
	dropped := 1024
	// fmt.Printf("dropped: %d\n", bs[:dropped])
	ms.dropBytes(bs[:dropped])
	fmt.Printf("p1: %d\n", ms.escape(ctx))
}
