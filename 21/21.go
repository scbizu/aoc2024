package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/magejiCoder/magejiAoc/input"
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

func (c *ctrl1) apply(b byte) []direction {
	fmt.Printf("current: %c, target: %c\n", c.current, b)
	var ret []direction
	if c.current == b {
		return ret
	}
	switch c.current {
	case 'A':
		if b < '3' {
			ret = append(ret, left)
			c.current = '0'
			ret = append(ret, c.apply(b)...)
		} else {
			ret = append(ret, up)
			c.current = '3'
			ret = append(ret, c.apply(b)...)
		}
	case '0':
		if b == 'A' {
			ret = append(ret, right)
		} else {
			ret = append(ret, up)
			c.current = '2'
			ret = append(ret, c.apply(b)...)
		}
	default:
		if b == 'A' {
			ret = append(ret, c.apply('3')...)
			ret = append(ret, down)
			return ret
		}
		if c.current > b {
			sub := c.current - b
			switch sub {
			case 2:
				if c.current-'0' != 0 {
					ret = append(ret, c.apply(c.current-3)...)
					ret = append(ret, right)
				} else {
					ret = append(ret, c.apply(c.current-1)...)
				}
			case 1:
				ret = append(ret, left)
			case 3:
				ret = append(ret, down)
			default:
				ret = append(ret, c.apply(c.current-3)...)
				ret = append(ret, c.apply(b)...)
			}
		} else {
			sub := b - c.current
			switch sub {
			case 2:
				if c.current-'0' != 0 {
					ret = append(ret, c.apply(c.current+3)...)
					ret = append(ret, left)
				} else {
					ret = append(ret, c.apply(c.current+1)...)
				}
			case 1:
				ret = append(ret, right)
			case 3:
				ret = append(ret, up)
			default:
				ret = append(ret, c.apply(c.current+3)...)
				ret = append(ret, c.apply(b)...)
			}
		}
	}
	c.current = b
	return ret
}

// ctrl2 is the robot controller .
type ctrl2 struct {
	current byte

	node *node
}

func (c *ctrl2) init() {
	nodeA := &node{v: 'A'}
	nodeU := &node{v: '^'}
	nodeD := &node{v: 'v'}
	nodeL := &node{v: '<'}
	nodeR := &node{v: '>'}
	nodeA.neighbors = []*node{nodeU, nodeR}
	nodeU.neighbors = []*node{nodeA, nodeD}
	nodeD.neighbors = []*node{nodeU, nodeR, nodeL}
	nodeL.neighbors = []*node{nodeD}
	c.node = nodeA
}

type node struct {
	v         byte
	neighbors []*node
}

func (c2 *ctrl2) apply(fb byte) []direction {
	var ret []direction
	switch {
	case c2.current == 'A' && fb == '>':
		ret = append(ret, down)
		c2.current = '>'
	case c2.current == 'A' && fb == '^':
		ret = append(ret, up)
		c2.current = '^'
	case c2.current == '^' && fb == 'A':
		ret = append(ret, right)
		c2.current = 'A'
	case c2.current == '^' && fb == 'v':
		ret = append(ret, down)
		c2.current = 'v'
	case c2.current == 'v' && fb == '>':
		ret = append(ret, right)
		c2.current = '>'
	case c2.current == 'v' && fb == '<':
		ret = append(ret, left)
		c2.current = '<'
	case c2.current == 'v' && fb == '^':
		ret = append(ret, up)
		c2.current = '^'
	case c2.current == '<' && fb == 'v':
		ret = append(ret, right)
		c2.current = 'v'
	}
	return ret
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("21.txt")
	var seqs []string
	txt.ReadByLineEx(ctx, func(_ int, line string) error {
		seqs = append(seqs, line)
		return nil
	})
	buf := bytes.NewBuffer(nil)
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
		buf.WriteString(bf.String())
	}
	fmt.Println(buf.String())
}

func main() {
	ctx := context.Background()
	p1(ctx)
}
