package main

import "github.com/cespare/next/container/set"

func init() {
	addSolutions(6, problem6)
}

func problem6(ctx *problemContext) {
	var p patrol
	s := ctx.scanner()
	for s.scan() {
		p.addRow([]byte(s.text()))
	}
	p.init()
	ctx.reportLoad()

	p1 := p.clone()
	for !p1.advance() {
	}
	var part1 int
	for _, c := range p1.g.all() {
		if c == 'X' {
			part1++
		}
	}
	ctx.reportPart1(part1)

	ctx.reportPart2(p.countObstacles())
}

type patrol struct {
	g   grid[byte]
	pos vec2
	d   vec2
}

func (p *patrol) clone() *patrol {
	p1 := *p
	p1.g = *p1.g.clone()
	return &p1
}

func (p *patrol) addRow(b []byte) {
	p.g.addRow(b)
}

func (p *patrol) init() {
	p.d = vec2{0, -1}
	for v, c := range p.g.all() {
		if c == '^' {
			p.pos = v
			p.g.set(v, '.')
			return
		}
	}
	panic("no guard")
}

func (p *patrol) advance() (done bool) {
	next := p.pos.add(p.d)
	if !p.g.contains(next) {
		return true
	}
	switch p.g.at(next) {
	case '.', 'X':
		p.g.set(next, 'X')
		p.pos = next
		return false
	case '#':
		p.d = p.d.matMul(rotations[3])
		return false
	default:
		panic("bad")
	}
}

func (p *patrol) countObstacles() int {
	seen := set.Of([2]vec2{p.pos, p.d})
	var total int
	newObstacles := set.Of(p.pos)
	for {
		next := p.pos.add(p.d)
		if !p.g.contains(next) {
			return total
		}
		switch p.g.at(next) {
		case '.':
			if !newObstacles.Contains(next) {
				// Place an obstacle.
				p1 := p.clone()
				p1.g.set(next, '#')
				if !p1.terminates() {
					total++
				}
				newObstacles.Add(next)
			}
			p.pos = next
		case '#':
			p.d = p.d.matMul(rotations[3])
		default:
			panic("bad")
		}
		state := [2]vec2{p.pos, p.d}
		if seen.Contains(state) {
			panic("bad")
		}
		seen.Add(state)
	}
}

func (p *patrol) terminates() bool {
	seen := set.Of([2]vec2{p.pos, p.d})
	for {
		if p.advance() {
			return true
		}
		state := [2]vec2{p.pos, p.d}
		if seen.Contains(state) {
			return false
		}
		seen.Add(state)
	}
}
