package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/set"
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
	swaps        map[string]string
}

func (w wire) getByOP(op op) []ins {
	var res []ins
	for _, i := range w.instructions {
		if i.op == op {
			res = append(res, i)
		}
	}
	return res
}

func (w *wire) run() {
	for _, i := range w.instructions {
		o := i.out
		if o2, ok := w.swaps[o]; ok {
			o = o2
		}
		w.output[o] = w.apply(i)
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
		swaps:        make(map[string]string),
	}
	w.run()
	var zsize int
	for s := range w.output {
		if strings.HasPrefix(s, "z") {
			zsize++
		}
	}
	zs := make([]int, zsize)
	for s, v := range w.output {
		if strings.HasPrefix(s, "z") {
			zs[toInt(strings.TrimPrefix(s, "z"))] = v
		}
	}
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
	wrongs := w.check()
	fmt.Printf("p2: %v\n", strings.Join(wrongs, ","))
}

// since it is a Ripple Carry Adder(RCA 行波加法器), it should follow the rules:
// 1. All XOR gates that inputs with x and y cannot every output to z (unless x0 and y0)
// 2. All other gates should output to z
// 3. All gates output to z should be XOR gates(except z45)
// 4. All gates in (1) must has input from (3)
// 5. at last , fix the sequence by removing the final OR gate
// Credit to /u/piman51277 , it is hard to read though
func (w *wire) check() []string {
	res := set.New[string]()
	fa0 := set.New[string]()
	fa3 := set.New[ins]()
	xors := w.getByOP(XOR)
	for _, i := range xors {
		if !strings.HasPrefix(i.in[0], "x") && !strings.HasPrefix(i.in[1], "x") {
			continue
		}
		fa0.Add(i.out)
		if i.in[0] == "x00" && i.in[1] == "y00" {
			if i.out != "z00" {
				// fmt.Printf("r1\n")
				res.Add(i.out)
			}
			continue
		} else {
			if i.out == "z00" {
				// fmt.Printf("r2\n")
				res.Add(i.out)
			}
		}
		if strings.HasPrefix(i.out, "z") {
			fmt.Printf("r3\n")
			res.Add(i.out)
		}
	}
	for _, i := range xors {
		if !strings.HasPrefix(i.in[0], "x") && !strings.HasPrefix(i.in[1], "x") {
			fa3.Add(i)
			if !strings.HasPrefix(i.out, "z") {
				// fmt.Printf("r4\n")
				res.Add(i.out)
			}
		}
	}

	for _, i := range w.instructions {
		if strings.HasPrefix(i.out, "z") {
			if i.out == "z45" {
				if i.op != OR {
					// fmt.Printf("r5\n")
					res.Add(i.out)
				}
				continue
			}
			if i.op != XOR {
				// fmt.Printf("r6\n")
				res.Add(i.out)
			}
		}
	}
	var rest []string
	fa0.Each(func(item string) bool {
		if res.Has(item) {
			return true
		}
		if item == "z00" {
			return true
		}
		var has bool
		fa3.Each(func(i ins) bool {
			if i.in[0] == item || i.in[1] == item {
				has = true
				return false
			}
			return true
		})
		if !has {
			// fmt.Printf("r7\n")
			res.Add(item)
			rest = append(rest, item)
		}
		return true
	})

	for _, o := range rest {
		index := strings.TrimPrefix(getInsByOut(w.instructions, o).in[0], "x")
		final2 := [2]string{}
		fa3.Each(func(item ins) bool {
			if item.out == "z"+index {
				final2 = item.in
				return false
			}
			return true
		})
		for i, f := range final2 {
			if getInsByOut(w.instructions, f).op == OR {
				if i == 0 {
					res.Add(final2[1])
				} else {
					res.Add(final2[0])
				}
				break
			}
		}
	}

	r := res.List()
	slices.Sort(r)
	return r
}

func getInsByOut(instructions []ins, out string) ins {
	for _, i := range instructions {
		if i.out == out {
			return i
		}
	}
	panic("not found")
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
	p1(ctx)
	p2(ctx)
}
