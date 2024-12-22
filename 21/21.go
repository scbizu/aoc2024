package main

import (
	"bytes"
	"context"
	"fmt"

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
	'2': {{num: '1', dir: left}, {num: '5', dir: up}, {num: '0', dir: down}},
	'3': {{num: 'A', dir: down}, {num: '6', dir: up}, {num: '1', dir: left}},
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
	// fmt.Printf("current: %c, target: %c\n", c.current, b)

	q := queue.NewQueue[state]()
	q.Push(state{current: c.current, path: []direction{}, visited: set.New(c.current)})
	var ret pathTrack

	for q.Len() > 0 {
		p := q.Pop()
		if p.current == b {
			ret = append(ret, dirsToString(p.path))
		}
		dirs := p.path
		for _, next := range numPad[p.current] {
			if p.visited.Has(next.num) {
				continue
			}
			pc := p.visited.Copy()
			pc.Add(next.num)
			q.Push(state{current: next.num, path: append(dirs, next.dir), visited: pc})
		}
	}
	return ret
}

// ctrl2 is the robot controller .
type ctrl2 struct {
	current byte
}

func (c2 *ctrl2) apply(fb byte) pathTrack {
	var ret pathTrack
	q := queue.NewQueue[state]()
	q.Push(state{current: c2.current, path: []direction{}, visited: set.New(c2.current)})
	for q.Len() > 0 {
		p := q.Pop()
		if p.current == fb {
			ret = append(ret, dirsToString(p.path))
		}
		dirs := p.path
		for _, next := range dirPad[p.current] {
			if p.visited.Has(next.num) {
				continue
			}
			pc := p.visited.Copy()
			pc.Add(next.num)
			q.Push(state{current: next.num, path: append(dirs, next.dir), visited: pc})
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

func p1(ctx context.Context) {
	txt := input.NewTXTFile("21.txt")
	var seqs []string
	txt.ReadByLineEx(ctx, func(_ int, line string) error {
		seqs = append(seqs, line)
		return nil
	})
	var outs []string
	for _, seq := range seqs {
		c1 := &ctrl1{current: 'A'}
		bf := bytes.NewBuffer(nil)
		for _, b := range []byte(seq) {
			dirs := c1.apply(b)
			for _, d := range dirs {
				bf.WriteString(d.String())
			}
			bf.WriteString("A")
		}
		fmt.Printf("c1: %s\n", bf.String())
		in := bf.Bytes()
		for i := 0; i < 2; i++ {
			c2 := &ctrl2{current: 'A'}
			c2b := bytes.NewBuffer(nil)
			for _, b := range in {
				dirs := c2.apply(b)
				for _, d := range dirs {
					c2b.WriteString(d.String())
				}
				c2b.WriteString("A")
			}
			in = c2b.Bytes()
			fmt.Printf("c2(%d): %s\n", i, c2b.String())
		}
		outs = append(outs, string(in))
	}
	var sum int
	for i, out := range outs {
		fmt.Printf("seq: %d, out: %d\n", toInt(seqs[i]), len(out))
		sum += len(out) * toInt(seqs[i])
	}
	fmt.Printf("p1: %d\n", sum)
}

func main() {
	ctx := context.Background()
	p1(ctx)
}
