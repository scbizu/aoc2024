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
	var ls []string
	var max int
	for c, s := range l.computers {
		s.Each(func(item string) bool {
			other := l.computers[item]
			if other.Has(c) {
				rest := set.Intersection(s, other)
				// fmt.Printf("c: %s, item: %s, rest: %v\n", c, item, rest)
				if rest.Size() > 0 {
					peers := set.New[string]()
					peers.Add(c, item)
					s := rest.List()
					// keep order of origin set
					slices.Sort(s)
					// pop 1
					peers.Add(s[0])
					buffer := set.New[string]()
					buffer.Add(s[0])
					for i := 1; i < len(s); i++ {
						// check if the next contains all the previous
						if l.computers[s[i]].Has(buffer.List()...) {
							buffer.Add(s[i])
							peers.Add(s[i])
						}
					}
					// fmt.Printf("peers: %v\n", peers.List())
					if peers.Size() > max {
						max = peers.Size()
						ls = peers.List()
					}
				}
			}
			return true
		})
	}
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
	fmt.Printf("p2: %v\n", strings.Join(group, ","))
}

func main() {
	ctx := context.Background()
	p1(ctx)
	p2(ctx)
}
