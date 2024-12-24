package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/set"
)

type lan struct {
	computers map[string]*set.Set[string]
}

func (l lan) findLargestGroup() []string {
	peers := set.New[string]()
	var max int
	for c, s := range l.computers {
		s.Each(func(item string) bool {
			other := l.computers[item]
			if other.Has(c) {
				rest := set.Intersection(s, other)
				fmt.Printf("c: %s, item: %s, rest: %v\n", c, item, rest)
				var full bool
				for i := 0; i < rest.Size(); i++ {
					for j := i + 1; j < rest.Size(); j++ {
						if !l.computers[rest.List()[i]].Has(rest.List()[j]) {
							// TODO: handle join
							return true
						}
					}
				}
				if full {
					if rest.Size() > max {
						max = rest.Size()
						peers.Clear()
						peers.Add(c)
						peers.Add(item)
						peers.Add(rest.List()...)
					}
				}
			}
			return true
		})
	}
	ls := peers.List()
	slices.Sort(ls)
	return ls
}

func (l lan) findPeers() []peer {
	peers := set.New[peer]()
	for c, s := range l.computers {
		s.Each(func(item string) bool {
			other := l.computers[item]
			if other.Has(c) {
				rest := set.Intersection(s, other)
				if rest.Size() > 0 {
					// fmt.Printf("c: %s, item: %s, rest: %v\n", c, item, rest)
					rest.Each(func(item2 string) bool {
						pr := peer{}
						pr[0] = c
						pr[1] = item
						pr[2] = item2
						slices.Sort(pr[:])
						peers.Add(pr)
						return true
					})
				}
			}
			return true
		})
	}
	return peers.List()
}

type link struct {
	l string
	r string
}

type peer [3]string

func buildLanParty(connections []link) lan {
	m := make(map[string]*set.Set[string])
	for _, conn := range connections {
		if _, ok := m[conn.l]; !ok {
			m[conn.l] = set.New[string]()
		}
		if _, ok := m[conn.r]; !ok {
			m[conn.r] = set.New[string]()
		}
		m[conn.l].Add(conn.r)
		m[conn.r].Add(conn.l)
	}
	return lan{computers: m}
}

func p1(ctx context.Context) {
	txt := input.NewTXTFile("23.txt")
	var links []link
	txt.ReadByLineEx(ctx, func(_ int, line string) error {
		parts := strings.Split(line, "-")
		links = append(links, link{l: parts[0], r: parts[1]})
		return nil
	})
	lan := buildLanParty(links)
	peers := lan.findPeers()
	fmt.Printf("peers: %v\n", peers)
	var ts int
	for _, peer := range peers {
		var hasT bool
		for _, c := range peer {
			if strings.HasPrefix(c, "t") {
				hasT = true
				break
			}
		}
		if hasT {
			ts++
		}
	}
	fmt.Printf("p1: %d\n", ts)
}

func p2(ctx context.Context) {
	txt := input.NewTXTFile("23.txt")
	var links []link
	txt.ReadByLineEx(ctx, func(_ int, line string) error {
		parts := strings.Split(line, "-")
		links = append(links, link{l: parts[0], r: parts[1]})
		return nil
	})
	lan := buildLanParty(links)
	group := lan.findLargestGroup()
	fmt.Printf("p2: %v\n", group)
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
