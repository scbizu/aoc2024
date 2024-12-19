package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/set"
)

type towel struct {
	stripes *set.Set[string]
	cache   map[string]int
}

func (t *towel) combinations(pattern string) int {
	// 这里虽然是全遍历，但是还是dfs + cache 好， 因为它可以在每次返回值的时候就开始构建cache，效率会高很多
	// 而 bfs 只能对每个 pattern 都构建一次 cache，效率低很多
	var dfs func(remaining string) int
	dfs = func(remaining string) int {
		// fmt.Printf("dfs:  remaining: %s\n", remaining)
		if len(remaining) == 0 {
			t.cache[pattern] += 1
			return 1
		}

		if result, ok := t.cache[remaining]; ok {
			return result
		}
		var total int
		t.stripes.Each(func(stripe string) bool {
			if strings.HasPrefix(remaining, stripe) {
				count := dfs(remaining[len(stripe):])
				if len(remaining[len(stripe):]) > 0 {
					// fmt.Printf("set cache: key: %s, value: %d\n", remaining[len(stripe):], count)
					t.cache[remaining[len(stripe):]] = count
				}
				total += count
			}
			return true
		})

		return total
	}

	return dfs(pattern)
}

func (t *towel) checkPattern(pattern string) bool {
	cache := make(map[string]bool)
	// 这里要用 dfs , 不能用 bfs
	// dfs + cache 的方式命中率要大大高于 bfs + cache 的方式
	// 并且对应的 stripe 越多，dfs 的效果越明显 , 相比之下 bfs 的分支太多了，而且 cache 只在很上层生效, build 的效率太低了
	var dfs func(remaining string) bool
	dfs = func(remaining string) bool {
		if len(remaining) == 0 {
			return true
		}

		if result, ok := cache[remaining]; ok {
			return result
		}

		result := false
		t.stripes.Each(func(stripe string) bool {
			if strings.HasPrefix(remaining, stripe) {
				if dfs(remaining[len(stripe):]) {
					result = true
					return false
				}
			}
			return true
		})

		cache[remaining] = result
		return result
	}

	return dfs(pattern)
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("19.txt")
	t := &towel{
		stripes: set.New[string](),
	}
	var patterns []string
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		// block0 is the stripes
		stripes := block[0]
		for _, stripe := range strings.Split(stripes, ", ") {
			t.stripes.Add(stripe)
		}
		// block1 is the patterns
		patterns = append(patterns, strings.Split(block[1], "\n")...)
		return nil
	})
	var total int
	for _, pattern := range patterns {
		if t.checkPattern(pattern) {
			// fmt.Printf("valid: %s\n", pattern)
			total++
		}
	}
	fmt.Printf("p1: %d\n", total)
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("19.txt")
	t := &towel{
		stripes: set.New[string](),
		cache:   make(map[string]int),
	}
	var patterns []string
	txt.ReadByBlock(ctx, "\n\n", func(block []string) error {
		// block0 is the stripes
		stripes := block[0]
		for _, stripe := range strings.Split(stripes, ", ") {
			t.stripes.Add(stripe)
		}
		// block1 is the patterns
		patterns = append(patterns, strings.Split(block[1], "\n")...)
		return nil
	})
	var total int
	for _, pattern := range patterns {
		total += t.combinations(pattern)
	}
	// fmt.Printf("cache: %v\n", t.cache)
	fmt.Printf("p2: %d\n", total)
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
