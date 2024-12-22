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
}

func (c2 *ctrl2) apply(fb byte) []direction {
	var ret []direction
	if c2.current == fb {
		return ret
	}
	switch c2.current {
	case 'A':
		switch fb {
		case '>':
			ret = append(ret, down)
			c2.current = '>'
		case '^':
			ret = append(ret, left)
			c2.current = '^'
		default:
			ret = append(ret, left)
			c2.current = '^'
			ret = append(ret, c2.apply(fb)...)
		}
	case '^':
		switch fb {
		case 'A':
			ret = append(ret, right)
			c2.current = 'A'
		case 'v':
			ret = append(ret, down)
			c2.current = 'v'
		default:
			ret = append(ret, down)
			c2.current = 'v'
			ret = append(ret, c2.apply(fb)...)
		}
	case 'v':
		switch fb {
		case '^':
			ret = append(ret, up)
			c2.current = '^'
		case '<':
			ret = append(ret, left)
			c2.current = '<'
		case '>':
			ret = append(ret, right)
			c2.current = '>'
		default:
			ret = append(ret, right)
			c2.current = '>'
			ret = append(ret, c2.apply(fb)...)
		}
	case '>':
		switch fb {
		case 'A':
			ret = append(ret, up)
			c2.current = 'A'
		case 'v':
			ret = append(ret, left)
			c2.current = 'v'
		default:
			ret = append(ret, left)
			c2.current = 'v'
			ret = append(ret, c2.apply(fb)...)
		}
	case '<':
		switch fb {
		case 'v':
			ret = append(ret, right)
			c2.current = 'v'
		default:
			ret = append(ret, right)
			c2.current = 'v'
			ret = append(ret, c2.apply(fb)...)
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
