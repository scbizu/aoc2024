package main

import (
	"context"
	"fmt"
	"math"
	"slices"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/queue"
	"github.com/magejiCoder/magejiAoc/set"
)

// ctrl1 is the digital controller .
type ctrl1 struct {
	current byte
}

type direction uint8

const (
	up direction = iota
	down
	left
	right
)

func (d direction) String() string {
	switch d {
	case up:
		return "^"
	case down:
		return "v"
	case left:
		return "<"
	case right:
		return ">"
	}
	return "unknown"
}

type pathTrack []string

type state struct {
	current byte
	path    []direction
	visited *set.Set[byte]
}

func dirsToString(dirs []direction) string {
	var ret string
	for _, d := range dirs {
		ret += d.String()
	}
	ret += "A"
	return ret
}

type numDir struct {
	num byte
	dir direction
}

var numPad = map[byte][]numDir{
	'A': {{num: '0', dir: left}, {num: '3', dir: up}},
	'0': {{num: 'A', dir: right}, {num: '2', dir: up}},
	'1': {{num: '2', dir: right}, {num: '4', dir: up}},
	'2': {{num: '1', dir: left}, {num: '5', dir: up}, {num: '0', dir: down}, {num: '3', dir: right}},
	'3': {{num: 'A', dir: down}, {num: '6', dir: up}, {num: '2', dir: left}},
	'4': {{num: '1', dir: down}, {num: '7', dir: up}, {num: '5', dir: right}},
	'5': {{num: '2', dir: down}, {num: '8', dir: up}, {num: '4', dir: left}, {num: '6', dir: right}},
	'6': {{num: '3', dir: down}, {num: '9', dir: up}, {num: '5', dir: left}},
	'7': {{num: '4', dir: down}, {num: '8', dir: right}},
	'8': {{num: '5', dir: down}, {num: '9', dir: right}, {num: '7', dir: left}},
	'9': {{num: '6', dir: down}, {num: '8', dir: left}},
}

var dirPad = map[byte][]numDir{
	'^': {{num: 'v', dir: down}, {num: 'A', dir: right}},
	'v': {{num: '^', dir: up}, {num: '<', dir: left}, {num: '>', dir: right}},
	'<': {{num: 'v', dir: right}},
	'>': {{num: 'v', dir: left}, {num: 'A', dir: up}},
	'A': {{num: '^', dir: left}, {num: '>', dir: down}},
}

func (c *ctrl1) apply(b byte) pathTrack {
	q := queue.NewQueue[state]()
	q.Push(state{current: c.current, path: []direction{}, visited: set.New(c.current)})
	var ret pathTrack

	for q.Len() > 0 {
		p := q.Pop()
		// fmt.Printf("current: %c, path: %v\n", p.current, p.path)
		if p.current == b {
			// fmt.Printf("dir: %v\n", dirsToString(p.path))
			c.current = b
			ret = append(ret, dirsToString(p.path))
		}
		dirs := p.path
		for _, next := range numPad[p.current] {
			if p.visited.Has(next.num) {
				continue
			}
			// fmt.Printf("next: %c,dir: %v\n", next.num, next.dir)
			pc := p.visited.Copy()
			pc.Add(next.num)
			np := slices.Clone(dirs)
			q.Push(state{current: next.num, path: append(np, next.dir), visited: pc})
		}
	}
	return ret
}

// ctrl2 is the robot controller .
type ctrl2 struct {
	current   byte
	ctrlCache map[target]pathTrack
}

type target struct {
	from byte
	to   byte
}

func (c2 *ctrl2) apply(fb byte, last bool) pathTrack {
	var ret pathTrack
	q := queue.NewQueue[state]()
	q.Push(state{current: c2.current, path: []direction{}, visited: set.New(c2.current)})
	for q.Len() > 0 {
		p := q.Pop()
		if p.current == fb {
			ret = append(ret, dirsToString(p.path))
			c2.current = fb
			if last {
				break
			}
		}
		dirs := p.path
		for _, next := range dirPad[p.current] {
			if p.visited.Has(next.num) {
				continue
			}
			pc := p.visited.Copy()
			pc.Add(next.num)
			pn := slices.Clone(dirs)
			q.Push(state{current: next.num, path: append(pn, next.dir), visited: pc})
		}
	}
	return ret
}

func toInt(seq string) int {
	var intPart string
	for _, b := range []byte(seq) {
		if '0' <= b && b <= '9' {
			intPart += string(b)
		}
	}
	return input.Atoi(intPart)
}

type item struct {
	index int
	path  []string
}

type key struct {
	from   byte
	to     byte
	nRobot int
}

// cache for the digital pad mapping
var cache = map[key]int{}

// buildTogether builds the shortest path length for the robot and the digital pad.
// For every digitalPad x->y path , per to robots numbers , the result is fixed.
// For Example: (029A)
// 1. A->0 (Robots: 3 (plus 1 for my pad)): (A -> < , A) (Robots: 2): (A -> <, , < -> A, A) (Robots: 1)
// 2. 0->2 (Robots: 3) : xxx (Robots: 2) : yyy (Robots: 1)
// 3. 2->9
func buildTogether(seq []byte, robots int) int {
	var length int
	var from byte = 'A'
	for i := 0; i < len(seq); i++ {
		to := seq[i]
		r := robots
		ds := buildDigitalPad([]byte{to})
		for r > 0 {
			if l, ok := cache[key{from: from, to: to, nRobot: r}]; ok {
				length += l
				break
			}
			for _, d := range ds {
				r2 := buildRobotPad([]byte(d), false)
				if len(filterShortest(r2)[0]) <= length {
					length = len(filterShortest(r2)[0])
				}
			}
		}
		from = to
	}
	return length
}

func buildDigitalPad(seq []byte) []string {
	var ret []string
	q := queue.NewQueue[item]()
	q.Push(item{index: 0, path: []string{}})
	c1 := &ctrl1{current: 'A'}
	for q.Len() > 0 {
		it := q.Pop()
		if it.index >= len(seq) {
			ret = it.path
			break
		}
		paths := c1.apply(seq[it.index])
		var nps []string
		if len(it.path) == 0 {
			nps = paths
		} else {
			for _, p := range it.path {
				for _, path := range paths {
					nps = append(nps, p+path)
				}
			}
		}
		// fmt.Printf("nps: %v\n", paths)
		q.Push(item{index: it.index + 1, path: nps})
	}
	// filter out the shortest ret
	minLen := len(ret[0])
	for _, r := range ret {
		if len(r) < minLen {
			minLen = len(r)
		}
	}
	var nret []string
	for _, r := range ret {
		if len(r) == minLen {
			nret = append(nret, r)
		}
	}
	return nret
}

func buildRobotPad(seq []byte, last bool) []string {
	var ret []string
	q := queue.NewQueue[item]()
	q.Push(item{index: 0, path: []string{}})
	c2 := &ctrl2{current: 'A'}
	for q.Len() > 0 {
		it := q.Pop()
		if it.index >= len(seq) {
			ret = it.path
			break
		}
		paths := c2.apply(seq[it.index], last)
		var nps []string
		if len(it.path) == 0 {
			nps = paths
		} else {
			for _, p := range it.path {
				for _, path := range paths {
					nps = append(nps, p+path)
				}
			}
		}
		q.Push(item{index: it.index + 1, path: nps})
	}
	// filter out the shortest ret
	minLen := len(ret[0])
	for _, r := range ret {
		if len(r) < minLen {
			minLen = len(r)
		}
	}
	var nret []string
	for _, r := range ret {
		if len(r) == minLen {
			nret = append(nret, r)
		}
	}
	return nret
}

func filterShortest(ret []string) []string {
	minLen := len(ret[0])
	for _, r := range ret {
		if len(r) < minLen {
			minLen = len(r)
		}
	}
	var nret []string
	for _, r := range ret {
		if len(r) == minLen {
			nret = append(nret, r)
		}
	}
	return nret
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("21.txt")
	var seqs []string
	txt.ReadByLineEx(ctx, func(_ int, line string) error {
		if len(line) == 0 {
			return nil
		}
		seqs = append(seqs, line)
		return nil
	})
	fmt.Printf("seqs: %v\n", seqs)
	var sum int
	for _, seq := range seqs {
		next := buildDigitalPad([]byte(seq))
		// fmt.Printf("next: %v\n", next)
		for i := 0; i < 2; i++ {
			shortest := math.MaxInt32
			var nx []string
			for _, n := range next {
				// if i == 1 { fmt.Printf("n: %v\n", n)
				// }
				r2 := buildRobotPad([]byte(n), i == 1)
				if len(filterShortest(r2)[0]) <= shortest {
					shortest = len(filterShortest(r2)[0])
					nx = append(nx, filterShortest(r2)...)
				}
			}
			next = nx
			// fmt.Printf("next r: %v\n", next)
		}
		minLen := len(next[0])
		for _, r := range next {
			if len(r) < minLen {
				minLen = len(r)
			}
		}
		var nret []string
		for _, r := range next {
			if len(r) == minLen {
				nret = append(nret, r)
			}
		}
		// fmt.Printf("n: %v\n", len(nret[0]))
		sum += toInt(seq) * len(nret[0])
	}
	fmt.Printf("p1: %d\n", sum)
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("21.txt")
	var seqs []string
	txt.ReadByLineEx(ctx, func(_ int, line string) error {
		if len(line) == 0 {
			return nil
		}
		seqs = append(seqs, line)
		return nil
	})
	// fmt.Printf("seqs: %v\n", seqs)
	var sum int
	for _, seq := range seqs {
		next := buildDigitalPad([]byte(seq))
		// fmt.Printf("next: %v\n", next)
		for i := 0; i < 25; i++ {
			shortest := math.MaxInt32
			var nx []string
			for _, n := range next {
				// if i == 1 { fmt.Printf("n: %v\n", n)
				// }
				r2 := buildRobotPad([]byte(n), i == 24)
				if len(filterShortest(r2)[0]) <= shortest {
					shortest = len(filterShortest(r2)[0])
					nx = append(nx, filterShortest(r2)...)
				}
			}
			next = nx
			// fmt.Printf("next r: %v\n", next)
		}
		minLen := len(next[0])
		for _, r := range next {
			if len(r) < minLen {
				minLen = len(r)
			}
		}
		var nret []string
		for _, r := range next {
			if len(r) == minLen {
				nret = append(nret, r)
			}
		}
		// fmt.Printf("n: %v\n", len(nret[0]))
		sum += toInt(seq) * len(nret[0])
	}
	fmt.Printf("p2: %d\n", sum)
}

func main() {
	ctx := context.Background()
	p1(ctx)
	// p2(ctx)
}
