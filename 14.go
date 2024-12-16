package main

import (
	"strings"

	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(14, problem14)
}

func problem14(ctx *problemContext) {
	m := &robotMap{size: vec2{101, 103}}
	s := ctx.scanner()
	for s.scan() {
		m.bots = append(m.bots, parseRobot(s.text()))
	}
	ctx.reportLoad()

	for range 100 {
		m.advance()
	}
	ctx.reportPart1(m.safety())

	for {
		if m.isTree() {
			// fmt.Printf("After %d steps:", m.step)
			// fmt.Println(m)
			ctx.reportPart2(m.step)
			return
		}
		m.advance()
	}
}

type robotMap struct {
	bots []robot
	size vec2
	step int64
}

type robot struct {
	p vec2
	v vec2
}

func parseRobot(s string) robot {
	var b robot
	sscanf(s, "p=%d,%d v=%d,%d", &b.p.x, &b.p.y, &b.v.x, &b.v.y)
	return b
}

func (m *robotMap) advance() {
	m.bots = SliceMap(m.bots, func(b robot) robot {
		vx := b.v.x
		for vx < 0 {
			vx += m.size.x
		}
		vy := b.v.y
		for vy < 0 {
			vy += m.size.y
		}
		return robot{
			p: vec2{
				(b.p.x + vx) % m.size.x,
				(b.p.y + vy) % m.size.y,
			},
			v: b.v,
		}
	})
	m.step++
}

func (m robotMap) String() string {
	var s set.Set[vec2]
	for _, b := range m.bots {
		s.Add(b.p)
	}
	var builder strings.Builder
	for y := int64(0); y < m.size.y; y++ {
		for x := int64(0); x < m.size.x; x++ {
			if s.Contains(vec2{x, y}) {
				builder.WriteByte('#')
			} else {
				builder.WriteByte(' ')
			}
		}
		builder.WriteByte('\n')
	}
	return builder.String()
}

func (m robotMap) safety() int64 {
	var quad [4]int64
	for _, b := range m.bots {
		var i int
		if b.p.x > m.size.x/2 {
			i |= 1
		} else if b.p.x == m.size.x/2 {
			continue
		}
		i <<= 1
		if b.p.y > m.size.y/2 {
			i |= 1
		} else if b.p.y == m.size.y/2 {
			continue
		}
		quad[i]++
	}
	return quad[0] * quad[1] * quad[2] * quad[3]
}

func (m robotMap) isTree() bool {
	var s set.Set[vec2]
	for _, b := range m.bots {
		s.Add(b.p)
	}
	return s.Len() == len(m.bots)
}
