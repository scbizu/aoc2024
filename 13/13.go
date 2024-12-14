package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/queue"
)

type Claw struct {
	A, B  Axis
	Prize Axis
	cache map[try]struct{}
}

func (c Claw) String() string {
	return fmt.Sprintf("A: %v, B: %v, Prize: %v", c.A, c.B, c.Prize)
}

func (ax Axis) String() string {
	return fmt.Sprintf("X+%d, Y+%d", ax.x, ax.y)
}

type Axis struct {
	x, y int
}

type try struct {
	a, b   int
	rx, ry int
}

func (c *Claw) Guess2(ctx context.Context) int {
	c.Prize.x += 10000000000000
	c.Prize.y += 10000000000000
	// 使用克莱姆法则获得数学解
	// | prize.x b.x |
	// | prize.y b.y |
	// ---------------
	// | a.x b.x |     <= det
	// | a.y b.y |
	n1 := (c.Prize.x*c.B.y - c.Prize.y*c.B.x) / (c.A.x*c.B.y - c.A.y*c.B.x)
	n2 := (c.Prize.y*c.A.x - c.Prize.x*c.A.y) / (c.A.x*c.B.y - c.A.y*c.B.x)
	if n1*c.A.x+n2*c.B.x == c.Prize.x && n1*c.A.y+n2*c.B.y == c.Prize.y {
		// fmt.Printf("guess2: %d\n", 3*n1+n2)
		return 3*n1 + n2
	}
	return 0
}

func (c *Claw) Guess(ctx context.Context) int {
	t0 := try{
		a:  0,
		b:  0,
		rx: 0,
		ry: 0,
	}
	min := 3*100 + 100
	q := queue.NewQueue[try]()
	q.Push(t0)
	for q.Len() > 0 {
		t := q.Pop()
		if _, ok := c.cache[t]; ok {
			continue
		}
		// fmt.Printf("%v\n", t)
		if t.rx > c.Prize.x || t.ry > c.Prize.y {
			continue
		}
		if t.rx == c.Prize.x && t.ry == c.Prize.y {
			if 3*t.a+t.b < min {
				min = 3*t.a + t.b
			}
		}
		// 三种情况
		// 1. push A
		// 2. push B
		// 3. push A + push B
		if t.a < 100 && t.b < 100 {
			// push A + push B
			q.Push(try{
				a:  t.a + 1,
				b:  t.b + 1,
				rx: t.rx + c.A.x + c.B.x,
				ry: t.ry + c.A.y + c.B.y,
			})
			c.cache[t] = struct{}{}
		}
		if t.a < 100 {
			// push A
			q.Push(try{
				a:  t.a + 1,
				b:  t.b,
				rx: t.rx + c.A.x,
				ry: t.ry + c.A.y,
			})
			c.cache[t] = struct{}{}
		}
		if t.b < 100 {
			// push B
			q.Push(try{
				a:  t.a,
				b:  t.b + 1,
				rx: t.rx + c.B.x,
				ry: t.ry + c.B.y,
			})
			c.cache[t] = struct{}{}
		}
	}
	if min == 400 {
		return 0
	}
	return min
}

func parseToClaw(cw *Claw, s string) {
	switch {
	case strings.Contains(s, "Button A:"):
		sub := strings.TrimPrefix(s, "Button A: ")
		xy := strings.Split(sub, ", ")
		var xint, yint int
		_, err := fmt.Sscanf(xy[0], "X+%d", &xint)
		if err != nil {
			panic(err)
		}
		_, err = fmt.Sscanf(xy[1], "Y+%d", &yint)
		if err != nil {
			panic(err)
		}
		cw.A = Axis{xint, yint}
	case strings.Contains(s, "Button B:"):
		sub := strings.TrimPrefix(s, "Button B: ")
		xy := strings.Split(sub, ", ")
		var xint, yint int
		_, err := fmt.Sscanf(xy[0], "X+%d", &xint)
		if err != nil {
			panic(err)
		}
		_, err = fmt.Sscanf(xy[1], "Y+%d", &yint)
		if err != nil {
			panic(err)
		}
		cw.B = Axis{xint, yint}
	case strings.Contains(s, "Prize: "):
		sub := strings.TrimPrefix(s, "Prize: ")
		xy := strings.Split(sub, ", ")
		var xint, yint int
		// fmt.Printf("xy: %v\n", xy)
		_, err := fmt.Sscanf(xy[0], "X=%d", &xint)
		if err != nil {
			panic(err)
		}
		_, err = fmt.Sscanf(xy[1], "Y=%d", &yint)
		if err != nil {
			panic(err)
		}
		cw.Prize = Axis{xint, yint}
	}
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("13.txt")
	var cws []*Claw
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		// fmt.Printf("block: %v\n", len(block))
		for _, s := range block {
			cw := Claw{
				cache: make(map[try]struct{}),
			}
			parts := strings.Split(s, "\n")
			for _, s := range parts {
				parseToClaw(&cw, s)
			}
			cws = append(cws, &cw)
		}
		return nil
	})

	// fmt.Printf("claw: %v\n", cws)

	var total int
	for _, cw := range cws {
		total += cw.Guess(ctx)
	}
	fmt.Printf("p1: %d\n", total)
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("13.txt")
	var cws []*Claw
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		// fmt.Printf("block: %v\n", len(block))
		for _, s := range block {
			cw := Claw{
				cache: make(map[try]struct{}),
			}
			parts := strings.Split(s, "\n")
			for _, s := range parts {
				parseToClaw(&cw, s)
			}
			cws = append(cws, &cw)
		}
		return nil
	})

	// fmt.Printf("claw: %v\n", cws)

	var total int
	for _, cw := range cws {
		total += cw.Guess2(ctx)
	}
	fmt.Printf("p2: %d\n", total)
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
