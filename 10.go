package main

import (
	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(10, problem10)
}

func problem10(ctx *problemContext) {
	var m topoMap
	s := ctx.scanner()
	for s.scan() {
		m.g.addRow([]byte(s.text()))
	}
	ctx.reportLoad()

	ctx.reportPart1(m.trailheadSum())
	ctx.reportPart2(m.trailheadRatingSum())
}

type topoMap struct {
	g grid[byte]
}

func (m *topoMap) trailheadSum() int64 {
	var sum int64
	for v, c := range m.g.all() {
		if c == '0' {
			sum += m.trailheadScore(v, '0', new(set.Set[vec2]))
		}
	}
	return sum
}

func (m *topoMap) trailheadScore(v vec2, c byte, seen *set.Set[vec2]) int64 {
	var total int64
	for _, n := range v.neighbors4() {
		if !m.g.contains(n) {
			continue
		}
		nc := m.g.at(n)
		if nc != c+1 {
			continue
		}
		if seen.Contains(n) {
			continue
		}
		seen.Add(n)
		if nc == '9' {
			total++
			continue
		}
		total += m.trailheadScore(n, nc, seen)
	}
	return total
}

func (m *topoMap) trailheadRatingSum() int64 {
	var sum int64
	for v, c := range m.g.all() {
		if c == '0' {
			sum += m.trailheadRating(v, '0')
		}
	}
	return sum
}

func (m *topoMap) trailheadRating(v vec2, c byte) int64 {
	var total int64
	for _, n := range v.neighbors4() {
		if !m.g.contains(n) {
			continue
		}
		nc := m.g.at(n)
		if nc != c+1 {
			continue
		}
		if nc == '9' {
			total++
			continue
		}
		total += m.trailheadRating(n, nc)
	}
	return total
}
