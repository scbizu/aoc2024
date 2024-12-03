package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
)

type sec struct {
	line string
}

func atoi(s string) int {
	var numStr string
	for _, c := range s {
		if numStr != "" && (c < '0' || c > '9') {
			return 0
		}
		if c >= '0' && c <= '9' {
			numStr += string(c)
		}
	}
	return input.Atoi(numStr)
}

func (s sec) parseMul() int {
	parts := strings.Split(s.line, "mul")
	var total int
	for _, part := range parts {
		var leftIndex, mdIndex, rightIndex int = -1, -1, -1
		fmt.Println(part)
		for i, c := range part {
			if leftIndex != -1 && rightIndex != -1 && mdIndex != -1 {
				break
			}
			if c == '(' && leftIndex == -1 && i == 0 {
				leftIndex = i
			}
			if c == ')' && rightIndex == -1 {
				rightIndex = i
			}
			if c == ',' && mdIndex == -1 {
				mdIndex = i
			}
		}
		fmt.Printf("left: %d,md: %d,right: %d\n", leftIndex, mdIndex, rightIndex)
		if leftIndex != -1 && rightIndex != -1 && mdIndex != -1 &&
			rightIndex > mdIndex && mdIndex > leftIndex {
			leftNumber := atoi(part[leftIndex+1 : mdIndex])
			rightNumber := atoi(part[mdIndex+1 : rightIndex])
			fmt.Printf("%d,%d\n", leftNumber, rightNumber)
			total += leftNumber * rightNumber
		}
	}
	return total
}

func p1() {
	txt := input.NewTXTFile("03.txt")
	var section []sec
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		section = append(section, sec{
			line: line,
		})
		return nil
	})
	var total int
	for _, s := range section {
		total += s.parseMul()
	}
	println("p1: ", total)
}

func main() {
	p1()
}
