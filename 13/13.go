package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/queue"
)

type Claw struct {
	A, B  Axia
	Prize Axia
}

func (c Claw) String() string {
	return fmt.Sprintf("A: %v, B: %v, Prize: %v", c.A, c.B, c.Prize)
}

func (ax Axia) String() string {
	return fmt.Sprintf("X+%d, Y+%d", ax.x, ax.y)
}

type Axia struct {
	x, y int
}

type try struct {
	a, b   int
	rx, ry int
}

func (c Claw) Guess(ctx context.Context) int {
	t0 := try{
		a:  100,
		b:  100,
		rx: 0,
		ry: 0,
	}
	q := queue.NewQueue[try]()
	q.Push(t0)
	for q.Len() > 0 {
		t := q.Pop()
		if t.rx == c.Prize.x && t.ry == c.Prize.y {
			return t.a + t.b
		}
		if t.a > 0 {
			// push A
			q.Push(try{
				a:  t.a - 1,
				b:  t.b,
				rx: t.rx + c.A.x,
				ry: t.ry + c.A.y,
			})
		}
		if t.b > 0 {
			// push B
			q.Push(try{
				a:  t.a,
				b:  t.b - 1,
				rx: t.rx + c.B.x,
				ry: t.ry + c.B.y,
			})
		}
	}
	return 0
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
		cw.A = Axia{xint, yint}
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
		cw.B = Axia{xint, yint}
	case strings.Contains(s, "Prize: "):
		sub := strings.TrimPrefix(s, "Prize: ")
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
		cw.Prize = Axia{xint, yint}
	}
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("13.txt")
	var cws []Claw
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		cw := Claw{}
		fmt.Printf("block: %v\n", len(block))
		for _, s := range block {
			parseToClaw(&cw, s)
		}
		cws = append(cws, cw)
		return nil
	})

	fmt.Printf("claw: %v\n", cws)

	var total int
	for _, cw := range cws {
		total += cw.Guess(ctx)
	}
	fmt.Printf("p1: %d\n", total)
}

func main() {
	ctx := context.Background()
	p1(ctx)
}
