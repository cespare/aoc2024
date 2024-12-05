package main

import (
	"strings"

	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(5, problem5)
}

func problem5(ctx *problemContext) {
	m := &safetyManual{
		deps: make(map[int64][]int64),
	}
	s := ctx.scanner()
	inRules := true
	for s.scan() {
		line := s.text()
		if line == "" {
			inRules = false
			continue
		}
		if inRules {
			s0, s1, ok := strings.Cut(line, "|")
			if !ok {
				panic("bad")
			}
			n0, n1 := parseInt(s0), parseInt(s1)
			m.deps[n1] = append(m.deps[n1], n0)
			continue
		}
		upd := SliceMap(strings.Split(line, ","), parseInt)
		m.updates = append(m.updates, upd)
	}
	ctx.reportLoad()

	var part1 int64
	for _, upd := range m.updates {
		if m.check(upd) {
			part1 += upd[len(upd)/2]
		}
	}
	ctx.reportPart1(part1)

	var part2 int64
	for _, upd := range m.updates {
		if m.check(upd) {
			continue
		}
		ord := topologicalSort(m.depSubset(upd))
		part2 += ord[len(ord)/2]
	}
	ctx.reportPart1(part2)
}

type safetyManual struct {
	deps    map[int64][]int64
	updates [][]int64
}

func (m *safetyManual) check(upd []int64) bool {
	all := set.Of(upd...)
	var seen set.Set[int64]
	for _, n := range upd {
		for _, dep := range m.deps[n] {
			if !all.Contains(dep) {
				continue
			}
			if !seen.Contains(dep) {
				return false
			}
		}
		seen.Add(n)
	}
	return true
}

func (m *safetyManual) depSubset(upd []int64) map[int64][]int64 {
	all := set.Of(upd...)
	subset := make(map[int64][]int64)
	for n, deps := range m.deps {
		if !all.Contains(n) {
			continue
		}
		var subsetDeps []int64
		for _, dep := range deps {
			if all.Contains(dep) {
				subsetDeps = append(subsetDeps, dep)
			}
		}
		if len(subsetDeps) > 0 {
			subset[n] = subsetDeps
		}
	}
	return subset
}

func topologicalSort(deps map[int64][]int64) []int64 {
	var sorted []int64
	visited := make(map[int64]bool) // false for visiting, true when complete
	var visit func(n int64)
	visit = func(n int64) {
		v, ok := visited[n]
		if ok {
			if !v {
				panic("cycle")
			}
			return // visited already
		}
		visited[n] = false
		for _, d := range deps[n] {
			visit(d)
		}
		visited[n] = true
		sorted = append(sorted, n)
	}
	for n := range deps {
		visit(n)
	}
	return sorted
}
