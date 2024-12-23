package main

import (
	"iter"
	"slices"
	"sort"
	"strings"

	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(23, problem23)
}

func problem23(ctx *problemContext) {
	var n network
	s := ctx.scanner()
	for s.scan() {
		c0, c1, ok := strings.Cut(s.text(), "-")
		if !ok {
			panic("bad")
		}
		n.addEdge(c0, c1)
	}
	ctx.reportLoad()

	var part1 int
outer:
	for triple := range n.triples() {
		for _, c := range triple {
			if c[0] == 't' {
				part1++
				continue outer
			}
		}
	}
	ctx.reportPart1(part1)

	ctx.reportPart2(strings.Join(n.lanParty(), ","))
}

type network struct {
	edges map[string]*set.Set[string]
}

func (n *network) addEdge(c0, c1 string) {
	if n.edges == nil {
		n.edges = make(map[string]*set.Set[string])
	}
	e0, ok := n.edges[c0]
	if !ok {
		e0 = new(set.Set[string])
		n.edges[c0] = e0
	}
	e0.Add(c1)
	e1, ok := n.edges[c1]
	if !ok {
		e1 = new(set.Set[string])
		n.edges[c1] = e1
	}
	e1.Add(c0)
}

func (n *network) triples() iter.Seq[[3]string] {
	var seen set.Set[[3]string]
	return func(yield func(triple [3]string) bool) {
		for c0, e0 := range n.edges {
			for c1 := range e0.All() {
				for c2 := range e0.All() {
					if c1 == c2 {
						continue
					}
					if !n.edges[c1].Contains(c2) {
						continue
					}
					triple := [3]string{c0, c1, c2}
					sort.Strings(triple[:])
					if seen.Contains(triple) {
						continue
					}
					seen.Add(triple)
					if !yield(triple) {
						return
					}
				}
			}
		}
	}
}

func (n *network) lanParty() []string {
	var longest []string
	for clique := range n.bronKerbosch() {
		if len(clique) > len(longest) {
			longest = clique
		}
	}
	return longest
}

func (n *network) bronKerbosch() iter.Seq[[]string] {
	return func(yield func([]string) bool) {
		var fn func(r, p, x *set.Set[string]) bool
		fn = func(r, p, x *set.Set[string]) bool {
			if p.Len() == 0 && x.Len() == 0 {
				clique := slices.Sorted(r.All())
				if !yield(clique) {
					return false
				}
			}
			for c := range p.All() {
				r1 := r.Clone()
				r1.Add(c)
				neighbors := n.edges[c]
				p1 := set.Intersection(p, neighbors)
				x1 := set.Intersection(x, neighbors)
				if !fn(r1, p1, x1) {
					return false
				}

				p.Remove(c)
				x.Add(c)
			}
			return true
		}
		r := new(set.Set[string])
		x := new(set.Set[string])
		p := new(set.Set[string])
		for c := range n.edges {
			p.Add(c)
		}
		fn(r, p, x)
	}
}
