package main

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/math"
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
}
