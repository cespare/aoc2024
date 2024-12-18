package main

import (
	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(18, problem18)
}

func problem18(ctx *problemContext) {
	var badBytes []vec2
	s := ctx.scanner()
	for s.scan() {
		var v vec2
		sscanf(s.text(), "%d,%d", &v.x, &v.y)
		badBytes = append(badBytes, v)
	}
	ctx.reportLoad()

	var m1 byteMaze
	m1.g.init(71, 71, '.')
	for _, v := range badBytes[:1024] {
		m1.g.set(v, '#')
	}
	ctx.reportPart1(m1.shortestPath())

	var m2 byteMaze
	m2.g.init(71, 71, '.')
	for _, v := range badBytes {
		m2.g.set(v, '#')
		if m2.shortestPath() < 0 {
			ctx.reportPart2(v)
			return
		}
	}
}

type byteMaze struct {
	g grid[byte]
	p vec2
}

func (m *byteMaze) shortestPath() int64 {
	type state struct {
		p    vec2
		cost int64
	}
	seen := set.Of(m.p)
	q := []state{{p: m.p}}
	for len(q) > 0 {
		s := SlicePop(&q)
		for _, n := range s.p.neighbors4() {
			if !m.g.contains(n) || seen.Contains(n) || m.g.at(n) != '.' {
				continue
			}
			seen.Add(n)
			cost := s.cost + 1
			if n.x == m.g.cols-1 && n.y == m.g.rows-1 {
				return cost
			}
			q = append(q, state{p: n, cost: cost})
		}
	}
	return -1
}
