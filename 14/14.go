package main

import (
	"context"
	"fmt"
	stdmath "math"
	"strings"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/math"
)

type maze struct {
	m          grid.VecMatrix[rune]
	wide, tall int
}

func (m maze) String() string {
	for i := 0; i < m.wide; i++ {
		for j := 0; j < m.tall; j++ {
			fmt.Printf("%c", m.m[grid.Vec{X: i, Y: j}])
		}
		fmt.Println()
	}
	return ""
}

type robot struct {
	loc grid.Vec
	mv  grid.Vec
}

func (m maze) AddRobots(rs []robot) {
	for _, r := range rs {
		m.m.Add(r.loc, '#')
	}
}

func (m maze) countQuadrant(rloc map[grid.Vec]int) [4]int {
	quadrants := [4]int{
		0, 0, 0, 0,
	}
	// q1
	for i := 0; i < m.wide/2; i++ {
		for j := 0; j < m.tall/2; j++ {
			if c, ok := rloc[grid.Vec{X: i, Y: j}]; ok {
				quadrants[0] += c
			}
		}
	}
	// q2
	for i := m.wide/2 + 1; i < m.wide; i++ {
		for j := 0; j < m.tall/2; j++ {
			if c, ok := rloc[grid.Vec{X: i, Y: j}]; ok {
				quadrants[1] += c
			}
		}
	}
	// q3
	for i := 0; i < m.wide/2; i++ {
		for j := m.tall/2 + 1; j < m.tall; j++ {
			if c, ok := rloc[grid.Vec{X: i, Y: j}]; ok {
				quadrants[2] += c
			}
		}
	}
	// q4
	for i := m.wide/2 + 1; i < m.wide; i++ {
		for j := m.tall/2 + 1; j < m.tall; j++ {
			if c, ok := rloc[grid.Vec{X: i, Y: j}]; ok {
				quadrants[3] += c
			}
		}
	}
	return quadrants
}

func (r robot) move(g maze, s int) robot {
	d := r.loc
	// fmt.Printf("d: %v\n", d)
	for i := 0; i < s; i++ {
		vec := d.Add(r.mv)
		if vec.X < 0 {
			// -1 => wide - 1
			vec.X = g.wide - math.Abs(vec.X)
		}
		if vec.Y < 0 {
			vec.Y = g.tall - math.Abs(vec.Y)
		}
		if vec.X >= g.wide {
			// wide => 0
			vec.X = vec.X - g.wide
		}
		if vec.Y >= g.tall {
			vec.Y = vec.Y - g.tall
		}
		// fmt.Printf("vec: %v\n", vec)
		d = vec
	}
	return robot{
		loc: d,
		mv:  r.mv,
	}
}

func parseRobot(s string) robot {
	parts := strings.Split(s, " ")
	p, v := parts[0], parts[1]
	var x, y int
	fmt.Sscanf(p, "p=%d,%d", &x, &y)
	var vx, vy int
	fmt.Sscanf(v, "v=%d,%d", &vx, &vy)
	return robot{
		loc: grid.Vec{X: x, Y: y},
		mv:  grid.Vec{X: vx, Y: vy},
	}
}

func buildMaze(line, col int) maze {
	mz := maze{
		m:    grid.NewVecMatrix[rune](),
		wide: line,
		tall: col,
	}
	for i := 0; i < line; i++ {
		for j := 0; j < col; j++ {
			mz.m.Add(grid.Vec{X: i, Y: j}, '.')
		}
	}
	return mz
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("14.txt")
	mz := buildMaze(101, 103)
	var robots []robot
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		r := parseRobot(line)
		robots = append(robots, r)
		return nil
	})
	finals := make(map[grid.Vec]int)
	for _, r := range robots {
		next := r.move(mz, 100)
		finals[next.loc] += 1
	}
	// fmt.Printf("finals: %v\n", finals)
	qt := mz.countQuadrant(finals)
	fmt.Printf("Q1: %d\nQ2: %d\nQ3: %d\nQ4: %d\n", qt[0], qt[1], qt[2], qt[3])
	total := 1
	for _, v := range qt {
		total *= v
	}
	fmt.Printf("p1: %d\n", total)
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("14.txt")
	mz := buildMaze(101, 103)
	var robots []robot
	txt.ReadByLineEx(ctx, func(i int, line string) error {
		r := parseRobot(line)
		robots = append(robots, r)
		return nil
	})
	i := 1
	for {
		var finals []robot
		for _, r := range robots {
			next := r.move(mz, 1)
			finals = append(finals, next)
		}
		if findTheTree(finals) {
			fmt.Printf("p2: %d\n", i)
			mz.AddRobots(finals)
			fmt.Printf("%s\n", mz)
			mz.m.Reset('.')
			break
		}
		robots = finals
		i++
	}
}

// 如果要组成圣诞树，大部分点都会聚合在一起，也就是说，这时候的标准差会很小
func findTheTree(robots []robot) bool {
	var xs, xy []int
	for _, r := range robots {
		xs = append(xs, r.loc.X)
		xy = append(xy, r.loc.Y)
	}
	stdx := stddev(xs)
	stdy := stddev(xy)
	// 19 是个 magic number，是根据实际情况调整出来的
	// 实际可以先从最差的情况(max(x, y))开始调整，然后逐渐减小
	if stdx < 19 && stdy < 19 {
		return true
	}
	return false
}

// 计算标准差，描述点的离散程度
func stddev(a []int) float64 {
	mean := 0.0
	for _, v := range a {
		mean += float64(v)
	}
	mean /= float64(len(a))
	variance := 0.0
	for _, v := range a {
		diff := float64(v) - mean
		variance += diff * diff
	}
	return stdmath.Sqrt(variance / float64(len(a)))
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
