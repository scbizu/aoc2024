package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/set"
)

type Direction int

const (
	up Direction = iota
	down
	left
	right
)

func (d Direction) vec() grid.Vec {
	switch d {
	case up:
		return grid.Vec{X: 0, Y: -1}
	case down:
		return grid.Vec{X: 0, Y: 1}
	case left:
		return grid.Vec{X: -1, Y: 0}
	case right:
		return grid.Vec{X: 1, Y: 0}
	}
	return grid.Vec{}
}

func parseToD(s string) Direction {
	switch s {
	case "^":
		return up
	case "v":
		return down
	case "<":
		return left
	case ">":
		return right
	}
	panic("invalid direction: " + string(s))
}

type maze struct {
	m     grid.VecMatrix[byte]
	peers *set.Set[peer]
	robot grid.Vec
}

type peer struct {
	l grid.Vec
	r grid.Vec
}

func (m maze) extend() grid.VecMatrix[byte] {
	nm := grid.VecMatrix[byte]{}
	m.m.ForEach(func(v grid.Vec, b byte) {
		if b == '#' {
			nm.Add(grid.Vec{
				X: 2 * v.X,
				Y: v.Y,
			}, '#')
			nm.Add(grid.Vec{
				X: 2*v.X + 1,
				Y: v.Y,
			}, '#')
		}
		if b == '.' {
			nm.Add(grid.Vec{
				X: 2 * v.X,
				Y: v.Y,
			}, '.')
			nm.Add(grid.Vec{
				X: 2*v.X + 1,
				Y: v.Y,
			}, '.')
		}
		if b == 'O' {
			nm.Add(grid.Vec{
				X: 2 * v.X,
				Y: v.Y,
			}, '[')
			nm.Add(grid.Vec{
				X: 2*v.X + 1,
				Y: v.Y,
			}, ']')
		}
		if b == '@' {
			nm.Add(grid.Vec{
				X: 2 * v.X,
				Y: v.Y,
			}, '@')
			nm.Add(grid.Vec{
				X: 2*v.X + 1,
				Y: v.Y,
			}, '.')
		}
	})
	return nm
}

type state struct {
	kind  string
	at    grid.Vec
	d     Direction
	peers *set.Set[peer]
}

func (m *maze) move2(ctx context.Context, s state) (grid.Vec, bool) {
	switch s.kind {
	case "robot":
		dst := s.at.Add(s.d.vec())
		c, ok := m.m.Get(dst)
		if !ok || c == '#' {
			m.robot = s.at
			return s.at, false
		}
		if c == '.' {
			m.robot = dst
			return dst, true
		}
		if c == '[' || c == ']' {
			_, ok := m.move2(ctx, state{kind: "box", at: dst, d: s.d})
			if !ok {
				m.robot = s.at
				m.peers = s.peers
				return s.at, false
			}
			m.robot = dst
			return dst, true
		}
	case "box":
		at := s.at
		atc, _ := m.m.Get(at)
		pr := [2]grid.Vec{}
		if atc == '[' {
			pr[0] = at
			pr[1] = at.Add(grid.Vec{X: 1, Y: 0})
		} else {
			pr[0] = at.Add(grid.Vec{X: -1, Y: 0})
			pr[1] = at
		}
		switch s.d {
		case up, down:
			dst1, dst2 := pr[0].Add(s.d.vec()), pr[1].Add(s.d.vec())
			c1, ok1 := m.m.Get(dst1)
			c2, ok2 := m.m.Get(dst2)
			if !ok1 || !ok2 || c1 == '#' || c2 == '#' {
				return s.at, false
			}
			// case 1:
			// []
			// ..
			if c1 == '.' && c2 == '.' {
				m.peers.Remove(peer{l: pr[0], r: pr[1]})
				m.peers.Add(peer{l: dst1, r: dst2})
				return dst1, true
			}
			// case 2:
			// []
			// []
			if c1 == '[' && c2 == ']' {
				_, ok := m.move2(ctx, state{kind: "box", at: dst1, d: s.d})
				if ok {
					m.peers.Remove(peer{l: pr[0], r: pr[1]})
					m.peers.Add(peer{l: dst1, r: dst2})
					return dst1, true
				} else {
					return s.at, false
				}
			}
			// case 3:
			// []      O . []
			// .[]     R .[].
			if c1 == '.' && c2 == '[' {
				_, ok := m.move2(ctx, state{kind: "box", at: dst2, d: s.d})
				if ok {
					m.peers.Remove(peer{l: pr[0], r: pr[1]})
					m.peers.Add(peer{l: dst1, r: dst2})
					return dst1, true
				} else {
					return s.at, false
				}
			}
			if c1 == ']' && c2 == '.' {
				_, ok := m.move2(ctx, state{kind: "box", at: dst1, d: s.d})
				if ok {
					m.peers.Remove(peer{l: pr[0], r: pr[1]})
					m.peers.Add(peer{l: dst1, r: dst2})
					return dst1, true
				} else {
					return s.at, false
				}
			}
			// case 4:
			// []     []     []
			//[][]   [][]   [][]
			//....   #...  [][][]
			//              #
			if c1 == ']' && c2 == '[' {
				_, ok1 := m.move2(ctx, state{kind: "box", at: dst1, d: s.d})
				if !ok1 {
					return s.at, false
				}
				_, ok2 := m.move2(ctx, state{kind: "box", at: dst2, d: s.d})
				if !ok2 {
					return s.at, false
				}
				m.peers.Remove(peer{l: pr[0], r: pr[1]})
				m.peers.Add(peer{l: dst1, r: dst2})
				return dst1, true
			}
		case left, right:
			dst1, dst2 := pr[0].Add(s.d.vec()), pr[1].Add(s.d.vec())
			c1, ok1 := m.m.Get(dst1)
			c2, ok2 := m.m.Get(dst2)
			if !ok1 || !ok2 || c1 == '#' || c2 == '#' {
				return s.at, false
			}
			if s.d == left && c1 == '.' {
				m.peers.Remove(peer{l: pr[0], r: pr[1]})
				m.peers.Add(peer{l: dst1, r: dst2})
				return dst1, true
			}
			if s.d == right && c2 == '.' {
				m.peers.Remove(peer{l: pr[0], r: pr[1]})
				m.peers.Add(peer{l: dst1, r: dst2})
				return dst1, true
			}
			if s.d == left && c1 == ']' {
				_, ok := m.move2(ctx, state{kind: "box", at: dst1, d: s.d})
				if ok {
					m.peers.Remove(peer{l: pr[0], r: pr[1]})
					m.peers.Add(peer{l: dst1, r: dst2})
					return dst1, true
				} else {
					return s.at, false
				}
			}
			if s.d == right && c2 == '[' {
				_, ok := m.move2(ctx, state{kind: "box", at: dst2, d: s.d})
				if ok {
					m.peers.Remove(peer{l: pr[0], r: pr[1]})
					m.peers.Add(peer{l: dst1, r: dst2})
					return dst1, true
				} else {
					return s.at, false
				}
			}
		}
	}
	panic("invalid state")
}

func (m *maze) move(ctx context.Context, s state) (grid.Vec, bool) {
	switch s.kind {
	case "robot":
		dst := s.at.Add(s.d.vec())
		c, ok := m.m.Get(dst)
		if !ok || c == '#' {
			return s.at, false
		}
		if c == '.' {
			m.m.Add(s.at, '.')
			m.m.Add(dst, '@')
			return dst, true
		}
		if c == 'O' {
			_, ok := m.move(ctx, state{kind: "box", at: dst, d: s.d})
			if !ok {
				return s.at, false
			}
			m.m.Add(s.at, '.')
			m.m.Add(dst, '@')
			return dst, true
		}
	case "box":
		dst := s.at.Add(s.d.vec())
		c, ok := m.m.Get(dst)
		if !ok || c == '#' {
			return s.at, false
		}
		if c == '.' {
			m.m.Add(s.at, '.')
			m.m.Add(dst, 'O')
			return dst, true
		}
		if c == 'O' {
			_, ok := m.move(ctx, state{kind: "box", at: dst, d: s.d})
			if ok {
				m.m.Add(s.at, '.')
				m.m.Add(dst, 'O')
				return dst, true
			} else {
				return s.at, false
			}
		}
	}
	panic("invalid state")
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("15.txt")
	mz := maze{
		m: grid.VecMatrix[byte]{},
	}
	var mvs []Direction
	var at grid.Vec
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		// block 1 is the maze
		mzPart := block[0]
		lines := strings.Split(mzPart, "\n")
		for i, line := range lines {
			for j, c := range line {
				if c == '@' {
					at = grid.Vec{X: j, Y: i}
				}
				mz.m.Add(grid.Vec{X: j, Y: i}, byte(c))
			}
		}
		// block 2 is the moves
		rbPart := block[1]
		for _, mv := range rbPart {
			if string(mv) == "\n" {
				continue
			}
			mvs = append(mvs, parseToD(string(mv)))
		}
		return nil
	})

	for _, mv := range mvs {
		at, _ = mz.move(ctx, state{
			kind: "robot",
			at:   at,
			d:    mv,
		})
		// fmt.Printf("move: %d,d: %v\n", idx, mv)
		// mz.m.Print(os.Stdout, "%c")
		// fmt.Println()
	}
	var total int
	mz.m.ForEach(func(v grid.Vec, b byte) {
		if b == 'O' {
			// fmt.Printf("O: %v\n", v)
			total += 100*v.Y + v.X
		}
	})
	fmt.Printf("p1: %d\n", total)
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("15.txt")
	mz := maze{
		m:     grid.VecMatrix[byte]{},
		peers: set.New[peer](),
	}
	var mvs []Direction
	var at grid.Vec
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		// block 1 is the maze
		mzPart := block[0]
		lines := strings.Split(mzPart, "\n")
		for i, line := range lines {
			for j, c := range line {
				mz.m.Add(grid.Vec{X: j, Y: i}, byte(c))
			}
		}
		// block 2 is the moves
		rbPart := block[1]
		for _, mv := range rbPart {
			if string(mv) == "\n" {
				continue
			}
			mvs = append(mvs, parseToD(string(mv)))
		}
		return nil
	})

	m := mz.extend()
	mz.m = m

	mz.m.Print(os.Stdout, "%c")

	mz.m.ForEach(func(v grid.Vec, b byte) {
		if b == '@' {
			at = v
			mz.robot = at
		}
		if b == '[' {
			mz.peers.Add(peer{l: v, r: v.Add(grid.Vec{X: 1, Y: 0})})
		}
	})

	// fmt.Printf("peers: %v\n", mz.peers)
	// fmt.Printf("@: %v\n", mz.robot)

	for _, mv := range mvs {
		// move
		at, _ = mz.move2(ctx, state{
			kind:  "robot",
			at:    at,
			d:     mv,
			peers: mz.peers.Copy(),
		})
		// fmt.Printf("peers: %v\n", mz.peers)
		// fmt.Printf("@: %v\n", mz.robot)
		// reset
		mz.m.ForEach(func(v grid.Vec, b byte) {
			if b == '@' {
				mz.m.Add(v, '.')
			}
			if b == '[' || b == ']' {
				mz.m.Add(v, '.')
			}
		})
		// load peers to maze
		for _, pr := range mz.peers.List() {
			mz.m.Add(pr.l, '[')
			mz.m.Add(pr.r, ']')
		}
		mz.m.Add(mz.robot, '@')
		// fmt.Printf("move: %d,d: %v\n", idx, mv)
		// mz.m.Print(os.Stdout, "%c")
		// fmt.Println()
	}
	// mz.m.Print(os.Stdout, "%c")
	// fmt.Println()
	var total int
	mz.m.ForEach(func(v grid.Vec, b byte) {
		if b == '[' {
			// fmt.Printf("O: %v\n", v)
			total += 100*v.Y + v.X
		}
	})
	fmt.Printf("p2: %d\n", total)
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
