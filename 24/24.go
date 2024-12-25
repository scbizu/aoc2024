package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
)

type ins struct {
	in  [2]string
	op  op
	out string
}

func (i ins) String() string {
	return fmt.Sprintf("%v %v -> %s", i.in, i.op, i.out)
}

type op uint8

const (
	AND op = iota
	OR
	XOR
)

func parseOP(opStr string) op {
	switch opStr {
	case "AND":
		return AND
	case "OR":
		return OR
	case "XOR":
		return XOR
	}
	panic("not a valid op")
}

func (w wire) buildMap() map[string]ins {
	m := make(map[string]ins)
	for _, i := range w.instructions {
		m[i.out] = i
	}
	return m
}

func (w *wire) apply(in ins) int {
	// fmt.Printf("in: %v\n", in)
	op1, op2 := in.in[0], in.in[1]
	var l, r int
	m := w.buildMap()
	if _, ok := w.inputs[op1]; !ok {
		l = w.apply(m[op1])
		w.output[op1] = l
		w.inputs[op1] = l
	} else {
		l = w.inputs[op1]
	}
	if _, ok := w.inputs[op2]; !ok {
		r = w.apply(m[op2])
		w.output[op2] = r
		w.inputs[op2] = r
	} else {
		r = w.inputs[op2]
	}
	switch in.op {
	case AND:
		return l & r
	case OR:
		return l | r
	case XOR:
		return l ^ r
	}
	panic("invalid op")
}

type wire struct {
	inputs       map[string]int
	instructions []ins
	output       map[string]int
}

func (w *wire) run() {
	for _, i := range w.instructions {
		w.output[i.out] = w.apply(i)
	}
}

// getX get all x in order as binary string
func (w *wire) getX() []int {
	var xsize int
	for x := range w.inputs {
		if strings.HasPrefix(x, "x") {
			xsize++
		}
	}
	xs := make([]int, xsize)
	for x, v := range w.inputs {
		if strings.HasPrefix(x, "x") {
			xs[xsize-toInt(strings.TrimPrefix(x, "x"))-1] = v
		}
	}
	return xs
}

func (w *wire) getY() []int {
	var ysize int
	for y := range w.inputs {
		if strings.HasPrefix(y, "y") {
			ysize++
		}
	}
	ys := make([]int, ysize)
	for y, v := range w.inputs {
		if strings.HasPrefix(y, "y") {
			ys[ysize-toInt(strings.TrimPrefix(y, "y"))-1] = v
		}
	}
	return ys
}

func (w *wire) getZ() []int {
	var zsize int
	for z := range w.output {
		if strings.HasPrefix(z, "z") {
			zsize++
		}
	}
	zs := make([]int, zsize)
	for z, v := range w.output {
		if strings.HasPrefix(z, "z") {
			zs[zsize-toInt(strings.TrimPrefix(z, "z"))-1] = v
		}
	}
	return zs
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("24.txt")

	var instructions []ins
	inputMap := make(map[string]int)
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		// block 0 is the input
		inputs := strings.Split(block[0], "\n")
		for _, in := range inputs {
			parts := strings.Split(in, ": ")
			inputMap[parts[0]] = input.Atoi(parts[1])
		}
		// block 1 is the instructions
		b1Parts := block[1]
		for _, line := range strings.Split(b1Parts, "\n") {
			insParts := strings.Split(line, " -> ")
			inParts := strings.Split(insParts[0], " ")
			instructions = append(instructions, ins{
				in:  [2]string{inParts[0], inParts[2]},
				op:  parseOP(inParts[1]),
				out: insParts[1],
			})
		}
		return nil
	})
	// fmt.Printf("ins: %v\n", instructions)
	w := &wire{
		inputs:       inputMap,
		instructions: instructions,
		output:       make(map[string]int),
	}
	w.run()
	fmt.Printf("output: %v\n", w.output)
	var zsize int
	for s := range w.output {
		if strings.HasPrefix(s, "z") {
			zsize++
		}
	}
	fmt.Printf("zsize: %v\n", zsize)
	zs := make([]int, zsize)
	for s, v := range w.output {
		if strings.HasPrefix(s, "z") {
			zs[toInt(strings.TrimPrefix(s, "z"))] = v
		}
	}
	fmt.Printf("%v\n", zs)
	fmt.Printf("p1: %d\n", binaryToDecimal(zs, false))
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("24.txt")

	var instructions []ins
	inputMap := make(map[string]int)
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		// block 0 is the input
		inputs := strings.Split(block[0], "\n")
		for _, in := range inputs {
			parts := strings.Split(in, ": ")
			inputMap[parts[0]] = input.Atoi(parts[1])
		}
		// block 1 is the instructions
		b1Parts := block[1]
		for _, line := range strings.Split(b1Parts, "\n") {
			insParts := strings.Split(line, " -> ")
			inParts := strings.Split(insParts[0], " ")
			instructions = append(instructions, ins{
				in:  [2]string{inParts[0], inParts[2]},
				op:  parseOP(inParts[1]),
				out: insParts[1],
			})
		}
		return nil
	})
	// fmt.Printf("ins: %v\n", instructions)
	w := &wire{
		inputs:       inputMap,
		instructions: instructions,
		output:       make(map[string]int),
	}
	w.run()
	xs, ys := w.getX(), w.getY()
	fmt.Printf("xs: %v;ys: %v\n", xs, ys)
	add := binAnd(xs, ys, true)
	zs := w.getZ()
	fmt.Printf("got: %v,exp: %v\n", zs, add)
}

func binAnd(a, b []int, reverse bool) []int {
	da := binaryToDecimal(a, reverse)
	db := binaryToDecimal(b, reverse)
	fmt.Printf("da: %v,db: %v\n", da, db)
	return decToBinary(da&db, reverse)
}

func decToBinary(dec int, reverse bool) []int {
	var bin []int
	for dec > 0 {
		bin = append(bin, dec%2)
		dec /= 2
	}
	if reverse {
		slices.Reverse(bin)
	}
	return bin
}

func toInt(s string) int {
	var s2 string
	for _, c := range s {
		if '0' <= c && c <= '9' {
			s2 += string(c)
		}
	}
	return input.Atoi(s2)
}

func binaryToDecimal(bin []int, reverse bool) int {
	var dec int
	if reverse {
		slices.Reverse(bin)
	}
	for i, b := range bin {
		dec += b << i
	}
	return dec
}

func main() {
	ctx := context.Background()
	// p1(ctx)
	p2(ctx)
}
