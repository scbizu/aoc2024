package main

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/math"
	"github.com/magejiCoder/magejiAoc/queue"
)

type computer struct {
	ra int
	rb int
	rc int

	pointerIndex *int
	program      []int
}

func (c *computer) combo(v int) int {
	if v >= 0 && v <= 3 {
		return v
	}
	if v == 4 {
		return c.ra
	}
	if v == 5 {
		return c.rb
	}
	if v == 6 {
		return c.rc
	}
	// 7 is reserved
	panic("invalid combo")
}

type state struct {
	rIndex int
	ra     int
}

// 因为 RC 的存在使得 RB/RC 的过程不可推导
// RB 取反之后再让 RA 右移赋值给RC 真的不知道想要干嘛，所以只能用逆向工程来硬解
// 先找到最后一个program代表的RA,再往前推，因为每次RA结果之后都会 / 8 (题目中的 0 3 序列)，所以找到第一个之后就开始左移3位找到下一个RA
func (c *computer) guess() int {
	// reverse the program
	pmc := slices.Clone(c.program)
	slices.Reverse(pmc)
	raq := queue.NewQueue[state]()
	raq.Push(state{
		rIndex: 0,
		ra:     0,
	})
	var results []int
	for raq.Len() > 0 {
		s := raq.Pop()
		if s.rIndex == len(pmc) {
			break
		}
		for i := 0 + s.ra; i < 8+s.ra; i++ {
			cm := &computer{
				rb:      c.rb,
				rc:      c.rc,
				program: c.program,
			}
			cm.ra = i
			buf := bytes.NewBuffer(nil)
			cm.load(buf)
			r := slices.Clone(pmc[:s.rIndex+1])
			slices.Reverse(r)
			// fmt.Printf("buf: %s\n", buf.String())
			// fmt.Printf("ra: %v\n", i)
			if buf.String() == joinInts(r) {
				if len(r) == len(c.program) {
					results = append(results, i)
					continue
				}
				a := i << 3
				// fmt.Printf("[%d]ra: %d, p: %d\n", s.rIndex, s.ra, i)
				raq.Push(state{
					rIndex: s.rIndex + 1,
					ra:     a,
				})
			}
		}
	}
	slices.Sort(results)
	// panic("no solution found")
	return results[0]
}

func joinInts(ints []int) string {
	var ss []string
	for _, i := range ints {
		ss = append(ss, strconv.Itoa(i))
	}
	return strings.Join(ss, ",")
}

func (c *computer) load(buf *bytes.Buffer) {
	var readIndex int
	var outs []string
	for readIndex < len(c.program) {
		// fmt.Printf("read index: %d\n", readIndex)
		opcode := c.program[readIndex]
		operand := c.program[readIndex+1]
		// fmt.Printf("opcode: %d, operand: %d\n", opcode, operand)
		switch opcode {
		case 0:
			c.adv(power2(c.combo(operand)))
		case 1:
			c.bxl(operand)
		case 2:
			c.bst(c.combo(operand))
		case 3:
			c.jnz(operand)
		case 4:
			c.bxc(operand)
		case 5:
			outs = append(outs, c.out(c.combo(operand)%8))
		case 6:
			c.bdv(power2(c.combo(operand)))
		case 7:
			c.cdv(power2(c.combo(operand)))
		}
		if c.pointerIndex == nil {
			readIndex += 2
		} else {
			readIndex = *c.pointerIndex
			c.pointerIndex = nil
		}
	}
	// fmt.Printf("ra: %d, rb: %d, rc: %d\n", c.ra, c.rb, c.rc)
	buf.WriteString(strings.Join(outs, ","))
}

func power2(v int) int {
	return math.Power(2, v)
}

func (c *computer) adv(v2 int) {
	c.ra = c.ra / v2
}

func (c *computer) bxl(v int) {
	c.rb = c.rb ^ v
}

func (c *computer) bst(v int) {
	c.rb = v % 8
}

func (c *computer) jnz(v int) {
	if c.ra == 0 {
		return
	}
	c.pointerIndex = &v
}

func (c *computer) bxc(_ int) {
	c.rb = c.rb ^ c.rc
}

// combo
func (c *computer) out(v int) string {
	var cs []string
	for _, c := range strconv.Itoa(v) {
		cs = append(cs, string(c))
	}
	return strings.Join(cs, ",")
}

func (c *computer) bdv(v2 int) {
	c.rb = c.ra / v2
}

func (c *computer) cdv(v2 int) {
	c.rc = c.ra / v2
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("17.txt")
	c := &computer{}
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		registers := block[0]
		parts := strings.Split(registers, "\n")
		ras := strings.TrimPrefix(parts[0], "Register A: ")
		ra := input.Atoi(ras)
		bas := strings.TrimPrefix(parts[1], "Register B: ")
		rb := input.Atoi(bas)
		cas := strings.TrimPrefix(parts[2], "Register C: ")
		rc := input.Atoi(cas)
		c.ra = ra
		c.rb = rb
		c.rc = rc
		program := block[1]
		programs := strings.TrimPrefix(program, "Program: ")
		pparts := strings.Split(programs, ",")
		for _, p := range pparts {
			c.program = append(c.program, input.Atoi(p))
		}
		return nil
	})
	fmt.Printf("p2: %d\n", c.guess())
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("17.txt")
	c := &computer{}
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		registers := block[0]
		parts := strings.Split(registers, "\n")
		ras := strings.TrimPrefix(parts[0], "Register A: ")
		ra := input.Atoi(ras)
		bas := strings.TrimPrefix(parts[1], "Register B: ")
		rb := input.Atoi(bas)
		cas := strings.TrimPrefix(parts[2], "Register C: ")
		rc := input.Atoi(cas)
		c.ra = ra
		c.rb = rb
		c.rc = rc
		program := block[1]
		programs := strings.TrimPrefix(program, "Program: ")
		pparts := strings.Split(programs, ",")
		for _, p := range pparts {
			c.program = append(c.program, input.Atoi(p))
		}
		return nil
	})
	buf := bytes.NewBuffer(nil)
	c.load(buf)
	fmt.Printf("p1: %s\n", buf.String())
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
