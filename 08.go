package main

import (
	"iter"

	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(8, problem8)
}

func problem8(ctx *problemContext) {
	var m antennaMap
	s := ctx.scanner()
	for s.scan() {
		m.g.addRow([]byte(s.text()))
	}
	m.fillFreqs()
	ctx.reportLoad()

	ctx.reportPart1(m.countAntinodes())
	ctx.reportPart2(m.countAntinodes2())
}

type antennaMap struct {
	g     grid[byte]
	freqs map[byte][]vec2
}

func (m *antennaMap) fillFreqs() {
	m.freqs = make(map[byte][]vec2)
	for v, c := range m.g.all() {
		if c == '.' {
			continue
		}
		m.freqs[c] = append(m.freqs[c], v)
	}
}

func (m *antennaMap) countAntinodes() int {
	var antinodes set.Set[vec2]
	for _, vs := range m.freqs {
		for i, v0 := range vs {
			for j := i + 1; j < len(vs); j++ {
				v1 := vs[j]
				for vn := range m.antinodesFor(v0, v1) {
					antinodes.Add(vn)
				}
			}
		}
	}
	return antinodes.Len()
}

func (m *antennaMap) antinodesFor(v0, v1 vec2) iter.Seq[vec2] {
	d0 := v1.sub(v0)
	d1 := v0.sub(v1)
	return func(yield func(vec2) bool) {
		if vn := v1.add(d0); m.g.contains(vn) {
			if !yield(vn) {
				return
			}
		}
		if vn := v0.add(d1); m.g.contains(vn) {
			if !yield(vn) {
				return
			}
		}
	}
}

func (m *antennaMap) countAntinodes2() int {
	var antinodes set.Set[vec2]
	for _, vs := range m.freqs {
		for i, v0 := range vs {
			for j := i + 1; j < len(vs); j++ {
				v1 := vs[j]
				for vn := range m.antinodesFor2(v0, v1) {
					antinodes.Add(vn)
				}
			}
		}
	}
	return antinodes.Len()
}

func (m *antennaMap) antinodesFor2(v0, v1 vec2) iter.Seq[vec2] {
	d0 := v1.sub(v0)
	d1 := v0.sub(v1)
	return func(yield func(vec2) bool) {
		for vn := v0; m.g.contains(vn); vn = vn.add(d0) {
			if !yield(vn) {
				return
			}
		}
		for vn := v1; m.g.contains(vn); vn = vn.add(d1) {
			if !yield(vn) {
				return
			}
		}
	}
}
