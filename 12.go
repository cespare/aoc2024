package main

import (
	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(12, problem12)
}

func problem12(ctx *problemContext) {
	var g garden
	s := ctx.scanner()
	for s.scan() {
		g.g.addRow([]byte(s.text()))
	}
	ctx.reportLoad()

	g.findRegions()

	ctx.reportPart1(g.totalPrice1())
	ctx.reportPart2(g.totalPrice2())
}

type garden struct {
	g       grid[byte]
	regions [][]vec2
}

func (g *garden) findRegions() {
	regions := make(map[vec2]int)
	var i int
	for v, c := range g.g.all() {
		if _, ok := regions[v]; !ok {
			g.markRegion(v, c, regions, i)
			i++
		}
	}
	g.regions = make([][]vec2, i)
	for v, i := range regions {
		g.regions[i] = append(g.regions[i], v)
	}
}

func (g *garden) markRegion(v vec2, c byte, regions map[vec2]int, i int) {
	regions[v] = i
	for _, n := range v.neighbors4() {
		if !g.g.contains(n) {
			continue
		}
		if g.g.at(n) != c {
			continue
		}
		if _, ok := regions[n]; ok {
			continue
		}
		g.markRegion(n, c, regions, i)
	}
}

func (g *garden) totalPrice1() int64 {
	var sum int64
	for _, reg := range g.regions {
		sum += int64(len(reg)) * int64(edgesForRegion(reg).Len())
	}
	return sum
}

func edgesForRegion(region []vec2) *set.Set[edge] {
	var s set.Set[edge]
	for _, v := range region {
		for _, fence := range []edge{
			{v, edgeHoriz | edgeAbove},
			{v, edgeAbove},
			{vec2{v.x, v.y + 1}, edgeHoriz},
			{vec2{v.x + 1, v.y}, 0},
		} {
			if m := fence.mirror(); s.Contains(m) {
				s.Remove(m)
			} else {
				s.Add(fence)
			}
		}
	}
	return &s
}

type edge struct {
	v    vec2
	attr edgeAttr
}

type edgeAttr uint8

const (
	edgeHoriz edgeAttr = 1 << iota
	edgeAbove          // left (vert) or above (horiz)
)

func (e edge) mirror() edge {
	e.attr ^= edgeAbove
	return e
}

func (g *garden) totalPrice2() int64 {
	var sum int64
	for _, reg := range g.regions {
		sum += int64(len(reg)) * countSides(reg)
	}
	return sum
}

func countSides(region []vec2) int64 {
	edges := edgesForRegion(region)
	edgeGroups := make(map[edge]int)
	var id int
	for e := range edges.All() {
		if _, ok := edgeGroups[e]; ok {
			continue
		}
		markEdge(e, edges, edgeGroups, id)
		id++
	}
	return int64(id)
}

func markEdge(e edge, edges *set.Set[edge], groups map[edge]int, id int) {
	groups[e] = id
	var neigh []edge
	if e.attr&edgeHoriz != 0 {
		neigh = []edge{
			{vec2{e.v.x - 1, e.v.y}, e.attr},
			{vec2{e.v.x + 1, e.v.y}, e.attr},
		}
	} else {
		neigh = []edge{
			{vec2{e.v.x, e.v.y - 1}, e.attr},
			{vec2{e.v.x, e.v.y + 1}, e.attr},
		}
	}
	for _, ne := range neigh {
		if _, ok := groups[ne]; ok {
			continue
		}
		if edges.Contains(ne) {
			markEdge(ne, edges, groups, id)
		}
	}
}
