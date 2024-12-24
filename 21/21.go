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

func toInt(seq string) int64 {
	var intPart string
	for _, b := range []byte(seq) {
		if '0' <= b && b <= '9' {
			intPart += string(b)
		}
	}
	return int64(input.Atoi(intPart))
}

type cacheItem struct {
	count int
	path  []string
}
type key2 struct {
	current    byte
	next       byte
	robotIndex int
}

func (k key2) String() string {
	return fmt.Sprintf("%c->%c:%d", k.current, k.next, k.robotIndex)
}

// cache has 5*5*n(robot) items
var (
	// cache2 存的是下一个robot的路径
	// 比如原始路径是 A -> 0，2个 robot
	// 那么第一个 robot 是 {current: A, next: <, robotIndex: 1}: {count: len(['<','v','<','A'])}
	// > 如果 robotIndex == 1 的再次接到 {'A'-> '<'} 的请求，那么会直接返回 4
	//   > 对于同一个 robotIndex 的路径，数据量是一定的, 所以总共的数量是 5 * 5 * n(robot的数量)
	// 我的 pad (第二个robot) 是 {current: A,next: <, robotIndex: 0}: {count: 4}
	// > 因为是最后一个 robot，所以直接返回 对应的最小路径，为 4
	cache2 = map[key2]cacheItem{}
)

// build builds the shortest path length for the robot and the digital pad.
// For every digitalPad x->y path , per to robots numbers , the result is fixed.
// For Example: (029A)
// 1. A->0 (Robots: 3 (plus 1 for my pad)): (A -> < , A) (Robots: 2): (A -> <, , < -> A, A) (Robots: 1)
// 2. 0->2 (Robots: 3) : xxx (Robots: 2) : yyy (Robots: 1)
// 3. 2->9
func build(seq []byte, totalRobots int) int64 {
	var from byte = 'A'
	var total int
	for i := 0; i < len(seq); i++ {
		to := seq[i]
		c := &ctrl1{current: from}
		paths := c.apply(to)
		_, next := filterShortest(paths)
		// fmt.Printf("c1: next: %v\n", next)
		var dfs func(current, target byte, robots int) int
		dfs = func(current, target byte, robots int) int {
			// fmt.Printf("current: %c, target: %c, robots: %d\n", current, target, robots)
			if min, ok := cache2[key2{current: current, next: target, robotIndex: robots}]; ok {
				// fmt.Printf("[%d][%c->%c]shortest: %d\n", robots, current, target, min.count)
				return min.count
			}
			c2 := &ctrl2{current: current}
			paths := c2.apply(target, false)
			min, next := filterShortest(paths)
			if robots == 0 {
				// fmt.Printf("c2: next: %v\n", next)
				return min
			}
			// fmt.Printf("c2: next: %v\n", next)
			var cr byte = 'A'
			shortest := math.MaxInt64
			for _, path := range next {
				var l int
				for _, c := range path {
					l += dfs(cr, byte(c), robots-1)
					cr = byte(c)
				}
				// fmt.Printf("l: %d\n", l)
				if l < shortest {
					shortest = l
				}
			}
			// fmt.Printf("[%d][%c->%c]shortest: %d\n", robots, current, target, shortest)
			cache2[key2{current: current, next: target, robotIndex: robots}] = cacheItem{count: shortest, path: next}
			return shortest
		}
		s := math.MaxInt64
		for _, nxt := range next {
			var l int
			var fromN byte = 'A'
			for _, n := range nxt {
				p := dfs(fromN, byte(n), totalRobots-1)
				// fmt.Printf("c1: [%c->%c]shortest: %d\n", fromN, n, p)
				l += p
				fromN = byte(n)
			}
			if l < s {
				s = l
			}
			// fmt.Printf("c1: [%s]shortest: %d\n", nxt, s)
		}
		total += s
		// fmt.Printf("[%c->%c]shortest: %d\n", from, to, shortest)
		from = to
	}
	return int64(total)
}

func mergePathTrack(pt1, pt2 pathTrack) pathTrack {
	var ret pathTrack
	for _, p := range pt1 {
		for _, p2 := range pt2 {
			ret = append(ret, p+p2)
		}
	}
	return ret
}

func filterShortest(ret []string) (int, []string) {
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
	return minLen, nret
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
	var sum int64
	for _, seq := range seqs {
		p := build([]byte(seq), 2)
		fmt.Printf("path: %d\n", p)
		sum += toInt(seq) * p
	}
	// fmt.Printf("cache2: %v\n", cache2)
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
	var sum int64
	for _, seq := range seqs {
		p := build([]byte(seq), 25)
		fmt.Printf("path: %d\n", p)
		sum += toInt(seq) * p
	}
	fmt.Printf("p2: %d\n", sum)
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
