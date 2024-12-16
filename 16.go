package main

import (
	"github.com/cespare/next/container/heap"
	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(16, problem16)
}

func problem16(ctx *problemContext) {
	var m reindeerMaze
	s := ctx.scanner()
	for s.scan() {
		m.g.addRow([]byte(s.text()))
	}
	m.init()
	ctx.reportLoad()

	best, visited := m.shortestPath()
	ctx.reportPart1(best)
	ctx.reportPart2(visited.Len())
}

type reindeerMaze struct {
	g     grid[byte]
	start vec2
	end   vec2
}

func (m *reindeerMaze) init() {
	for v, c := range m.g.all() {
		switch c {
		case 'S':
			m.start = v
			m.g.set(v, '.')
		case 'E':
			m.end = v
			m.g.set(v, '.')
		}
	}
}

func (m *reindeerMaze) shortestPath() (int64, *set.Set[vec2]) {
	type state struct {
		p vec2
		d vec2
	}
	type path struct {
		cost    int64
		visited *set.Set[vec2]
	}
	type queueState struct {
		state
		cost     int64
		cameFrom *state
	}
	less := func(s0, s1 queueState) bool { return s0.cost < s1.cost }
	q := heap.New(less)
	s0 := state{p: m.start, d: vec2{1, 0}}
	q.Push(queueState{state: s0})
	paths := make(map[state]path)
	leastCost := int64(-1)
	for q.Len() > 0 {
		s := q.Pop()
		if best, ok := paths[s.state]; ok {
			if s.cost < best.cost {
				panic("shouldn't happen")
			}
			if s.cost == best.cost {
				if s.cameFrom != nil {
					best.visited.AddSet(paths[*s.cameFrom].visited)
				}
			}
			continue
		}
		visited := set.Of(s.p)
		if s.cameFrom != nil {
			visited.AddSet(paths[*s.cameFrom].visited)
		}
		paths[s.state] = path{
			cost:    s.cost,
			visited: visited,
		}

		if s.p == m.end {
			if leastCost >= 0 && s.cost > leastCost {
				// Done.
				break
			}
			leastCost = s.cost
		}

		p1 := s.p.add(s.d)
		if m.g.at(p1) == '.' {
			q.Push(queueState{
				state:    state{p: p1, d: s.d},
				cost:     s.cost + 1,
				cameFrom: &s.state,
			})
		}
		cw := s.d.matMul(turnCW)
		q.Push(queueState{
			state:    state{p: s.p, d: cw},
			cost:     s.cost + 1000,
			cameFrom: &s.state,
		})
		ccw := s.d.matMul(turnCCW)
		q.Push(queueState{
			state:    state{p: s.p, d: ccw},
			cost:     s.cost + 1000,
			cameFrom: &s.state,
		})
	}
	if leastCost < 0 {
		panic("no solution")
	}
	// Tally up all the equally good ways to get there.
	var visited set.Set[vec2]
	for s, p := range paths {
		if s.p == m.end && p.cost == leastCost {
			visited.AddSet(p.visited)
		}
	}
	return leastCost, &visited
}
