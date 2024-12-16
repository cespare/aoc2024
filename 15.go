package main

import (
	"fmt"
	"slices"

	"github.com/cespare/next/container/set"
)

func init() {
	addSolutions(15, problem15)
}

func problem15(ctx *problemContext) {
	var w0 warehouse
	s := ctx.scanner()
	inGrid := true
	for s.scan() {
		line := s.text()
		if inGrid {
			if line == "" {
				inGrid = false
				continue
			}
			w0.g.addRow([]byte(line))
			continue
		}
		for _, c := range line {
			d, ok := caretDirs[c]
			if !ok {
				panic(c)
			}
			w0.moves = append(w0.moves, d)
		}
	}
	w0.init()
	w1 := makeThiccWarehouse(&w0)
	ctx.reportLoad()

	w0.runSimulation()
	ctx.reportPart1(w0.gpsSum())

	w1.runSimulation()
	fmt.Println(byteGridString(&w1.g))
	ctx.reportPart1(w1.gpsSum())
}

// TODO: unify warehouse and thiccWarehouse.

type warehouse struct {
	g     grid[byte]
	cur   vec2
	moves []vec2
}

func (w *warehouse) init() {
	for v, c := range w.g.all() {
		if c == '@' {
			w.cur = v
			w.g.set(v, '.')
			break
		}
	}
}

func (w *warehouse) runSimulation() {
	for len(w.moves) > 0 {
		w.step(SlicePop(&w.moves))
	}
}

func (w *warehouse) step(d vec2) {
	var move func(vec2, vec2, bool) bool
	move = func(v0, d vec2, box bool) bool {
		v := v0.add(d)
		switch w.g.at(v) {
		case '#':
			return false
		case '.':
		case 'O':
			if !move(v, d, true) {
				return false
			}
		default:
			panic("bad")
		}
		if box {
			w.g.set(v, 'O')
		} else {
			w.g.set(v, '.')
		}
		return true
	}
	if move(w.cur, d, false) {
		w.cur = w.cur.add(d)
	}
}

func (w *warehouse) gpsSum() int64 {
	var sum int64
	for v, c := range w.g.all() {
		if c == 'O' {
			sum += v.x + 100*v.y
		}
	}
	return sum
}

type thiccWarehouse struct {
	g     grid[byte]
	cur   vec2
	moves []vec2
}

func makeThiccWarehouse(w0 *warehouse) *thiccWarehouse {
	w := &thiccWarehouse{
		cur:   vec2{w0.cur.x * 2, w0.cur.y},
		moves: slices.Clone(w0.moves),
	}
	w.g.init(w0.g.rows, w0.g.cols*2, 'x')
	for v, c := range w0.g.all() {
		var s string
		switch c {
		case '#':
			s = "##"
		case '.':
			s = ".."
		case 'O':
			s = "[]"
		default:
			panic(c)
		}
		w.g.set(vec2{v.x * 2, v.y}, s[0])
		w.g.set(vec2{v.x*2 + 1, v.y}, s[1])
	}
	return w
}

func (w *thiccWarehouse) runSimulation() {
	for len(w.moves) > 0 {
		// fmt.Println(byteGridString(&w.g))
		w.step(SlicePop(&w.moves))
	}
}

func (w *thiccWarehouse) step(d vec2) {
	var toMove set.Set[vec2]
	var canMove func(vec2, vec2) bool
	canMove = func(v0, d vec2) bool {
		toMove.Add(v0)
		v := v0.add(d)
		switch w.g.at(v) {
		case '#':
			return false
		case '.':
			return true
		case '[':
			va, vb := v, vec2{v.x + 1, v.y}
			if d.x == 0 {
				return canMove(va, d) && canMove(vb, d)
			} else {
				if d.x < 0 {
					panic("bad")
				}
				toMove.Add(va)
				return canMove(vb, d)
			}
		case ']':
			va, vb := vec2{v.x - 1, v.y}, v
			if d.x == 0 {
				return canMove(va, d) && canMove(vb, d)
			} else {
				if d.x > 0 {
					panic("bad")
				}
				toMove.Add(vb)
				return canMove(va, d)
			}
		default:
			panic("bad")
		}
	}
	if !canMove(w.cur, d) {
		return
	}
	w.cur = w.cur.add(d)

	next := make(map[vec2]byte)
	for v := range toMove.All() {
		next[v.add(d)] = w.g.at(v)
	}
	for v := range toMove.All() {
		w.g.set(v, '.')
	}
	for v, c := range next {
		w.g.set(v, c)
	}
}

func (w *thiccWarehouse) gpsSum() int64 {
	var sum int64
	for v, c := range w.g.all() {
		if c == '[' {
			sum += v.x + 100*v.y
		}
	}
	return sum
}

var caretDirs = map[rune]vec2{
	'<': {-1, 0},
	'>': {1, 0},
	'^': {0, -1},
	'v': {0, 1},
}
