package main

import (
	"context"
	"fmt"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
)

type path []grid.Vec

type panel struct {
	g     grid.VecMatrix[string]
	paths []path
}

var findXMAS = func(c, v string) bool {
	switch c {
	case "X":
		return v == "M"
	case "M":
		return v == "A"
	case "A":
		return v == "S"
	}
	return false
}

var findMAS = func(c, v string) bool {
	switch c {
	case "M":
		return v == "A"
	case "A":
		return v == "S"
	}
	return false
}

type direction int

const (
	directionUpRight direction = iota
	directionUpLeft
	directionDownRight
	directionDownLeft
	directionUp
	directionDown
	directionRight
	directionLeft
)

func getNext(v grid.Vec, d direction) []grid.Vec {
	switch d {
	case directionUpRight:
		return []grid.Vec{
			{X: v.X + 1, Y: v.Y - 1},
		}
	case directionUpLeft:
		return []grid.Vec{
			{X: v.X - 1, Y: v.Y - 1},
		}
	case directionDownRight:
		return []grid.Vec{
			{X: v.X + 1, Y: v.Y + 1},
		}
	case directionDownLeft:
		return []grid.Vec{
			{X: v.X - 1, Y: v.Y + 1},
		}
	case directionUp:
		return []grid.Vec{
			{X: v.X, Y: v.Y - 1},
		}
	case directionDown:
		return []grid.Vec{
			{X: v.X, Y: v.Y + 1},
		}
	case directionRight:
		return []grid.Vec{
			{X: v.X + 1, Y: v.Y},
		}
	case directionLeft:
		return []grid.Vec{
			{X: v.X - 1, Y: v.Y},
		}
	}
	return nil
}

func (p *panel) traverse2(
	ctx context.Context,
	start grid.Vec,
	d direction,
	find func(c, v string) bool,
	finalPath []grid.Vec,
) {
	v, ok := p.g.Get(start)
	if !ok {
		return
	}
	if len(finalPath) == 3 {
		fmt.Printf("finalPath: %v\n", finalPath)
		p.paths = append(p.paths, finalPath)
		return
	}
	next := getNext(start, d)
	for _, n := range next {
		vv, ok := p.g.Get(n)
		if ok && find(v, vv) {
			finalPath = append(finalPath, n)
			p.traverse(ctx, n, d, find, finalPath)
		}
	}
}

func (p *panel) traverse(
	ctx context.Context,
	start grid.Vec,
	d direction,
	find func(c, v string) bool,
	finalPath []grid.Vec,
) {
	v, ok := p.g.Get(start)
	if !ok {
		return
	}
	if len(finalPath) == 4 {
		fmt.Printf("finalPath: %v\n", finalPath)
		p.paths = append(p.paths, finalPath)
		return
	}
	next := getNext(start, d)
	for _, n := range next {
		vv, ok := p.g.Get(n)
		if ok && find(v, vv) {
			finalPath = append(finalPath, n)
			p.traverse(ctx, n, d, find, finalPath)
		}
	}
}

func p1() {
	txt := input.NewTXTFile("04.txt")
	g := grid.NewVecMatrix[string]()
	ctx := context.Background()
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		for index, c := range line {
			g.Add(grid.Vec{
				X: index,
				Y: i,
			}, string(c))
		}
		return nil
	})
	p := panel{
		g: g,
	}
	p.g.ForEach(func(v grid.Vec, s string) {
		if s == "X" {
			directions := []direction{
				directionUpRight,
				directionUpLeft,
				directionDownRight,
				directionDownLeft,
				directionUp,
				directionDown,
				directionRight,
				directionLeft,
			}
			for _, d := range directions {
				p.traverse(ctx, v, d, findXMAS, []grid.Vec{v})
			}
		}
	})
	fmt.Printf("p1: %d\n", len(p.paths))
}

func p2() {
	txt := input.NewTXTFile("04.txt")
	g := grid.NewVecMatrix[string]()
	ctx := context.Background()
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		for index, c := range line {
			g.Add(grid.Vec{
				X: index,
				Y: i,
			}, string(c))
		}
		return nil
	})
	p := panel{
		g: g,
	}
	p.g.ForEach(func(v grid.Vec, s string) {
		if s == "M" {
			directions := []direction{
				directionUpRight,
				directionUpLeft,
				directionDownRight,
				directionDownLeft,
			}
			for _, d := range directions {
				p.traverse2(ctx, v, d, findMAS, []grid.Vec{v})
			}
		}
	})
	fmt.Printf("p2: %d\n", len(p.paths))
}

func main() {
	p1()
	p2()
}
