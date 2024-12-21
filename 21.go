package main

import (
	"fmt"
	"strings"
)

func init() {
	addSolutions(21, problem21)
}

func problem21(ctx *problemContext) {
	var codes []string
	s := ctx.scanner()
	for s.scan() {
		codes = append(codes, s.text())
	}
	ctx.reportLoad()

	var part1 int64
	for _, code := range codes {
		part1 += codeComplexity(code, 2)
	}
	ctx.reportPart1(part1)

	var part2 int64
	for _, code := range codes {
		part2 += codeComplexity(code, 25)
	}
	ctx.reportPart2(part2)
}

func codeComplexity(code string, dirLevels int) int64 {
	set := codeToMoveSet(code)
	for range dirLevels {
		set = set.levelUp()
	}
	n := parseInt(strings.TrimRight(strings.TrimLeft(code, "0"), "A"))
	return set.total() * n
}

type keypadMoveSet map[keypadMove]int64

func (s keypadMoveSet) total() int64 {
	var sum int64
	for move, n := range s {
		sum += int64(move.repeat) * n
	}
	return sum
}

func codeToMoveSet(code string) keypadMoveSet {
	set := make(keypadMoveSet)
	v0 := vec2{2, 3}
	for i := 0; i < len(code); i++ {
		v1 := codeToKeypad(code[i])
		repeat := 1
		for i < len(code)-1 && code[i+1] == code[i] {
			repeat++
			i++
		}
		seq := keypadPress(v0, v1, vec2{0, 3}, repeat)
		move := keypadMove{start: byte('A')}
		for len(seq) > 0 {
			move.end = seq[0]
			seq = seq[1:]
			move.repeat = 1
			for len(seq) > 0 && seq[0] == move.end {
				move.repeat++
				seq = seq[1:]
			}
			set[move]++
			move.start = move.end
		}
		v0 = v1
	}
	return set
}

func (s keypadMoveSet) levelUp() keypadMoveSet {
	next := make(keypadMoveSet)
	for move, n := range s {
		v0 := dirToKeypad(move.start)
		v1 := dirToKeypad(move.end)
		seq := keypadPress(v0, v1, vec2{0, 0}, move.repeat)
		if seq == "" {
			panic("asdf")
		}
		newMove := keypadMove{start: byte('A')}
		for len(seq) > 0 {
			newMove.end = seq[0]
			seq = seq[1:]
			newMove.repeat = 1
			for len(seq) > 0 && seq[0] == newMove.end {
				newMove.repeat++
				seq = seq[1:]
			}
			next[newMove] += n
			newMove.start = newMove.end
		}
	}
	return next
}

type keypadMove struct {
	start  byte
	end    byte
	repeat int
}

func (m keypadMove) String() string {
	return fmt.Sprintf("%c%c%d", m.start, m.end, m.repeat)
}

func keypadPress(v0, v1, forbidden vec2, repeat int) string {
	if v0 == v1 {
		return ""
	}
	var moves []vec2
	d := v1.sub(v0)
	// Move left first; the left arrow is in a bad location so we want to do
	// a double move to get there for efficiency.
	if d.x < 0 {
		moves = append(moves, vec2{d.x, 0})
	}
	if d.y != 0 {
		moves = append(moves, vec2{0, d.y})
	}
	if d.x > 0 {
		moves = append(moves, vec2{d.x, 0})
	}
	if len(moves) == 2 {
		if v0.add(moves[0]) == forbidden {
			moves[0], moves[1] = moves[1], moves[0]
		}
	}
	var b strings.Builder
	for _, move := range moves {
		var c byte
		switch {
		case move.x < 0:
			c = '<'
		case move.x > 0:
			c = '>'
		case move.y < 0:
			c = '^'
		case move.y > 0:
			c = 'v'
		default:
			panic("bad")
		}
		for range move.mag() {
			b.WriteByte(c)
		}
	}
	for range repeat {
		b.WriteByte('A')
	}
	return b.String()
}

func codeToKeypad(c byte) vec2 {
	switch c {
	case '0':
		return vec2{1, 3}
	case '1':
		return vec2{0, 2}
	case '2':
		return vec2{1, 2}
	case '3':
		return vec2{2, 2}
	case '4':
		return vec2{0, 1}
	case '5':
		return vec2{1, 1}
	case '6':
		return vec2{2, 1}
	case '7':
		return vec2{0, 0}
	case '8':
		return vec2{1, 0}
	case '9':
		return vec2{2, 0}
	case 'A':
		return vec2{2, 3}
	default:
		panic("bad")
	}
}

func dirToKeypad(c byte) vec2 {
	switch c {
	case '<':
		return vec2{0, 1}
	case '^':
		return vec2{1, 0}
	case '>':
		return vec2{2, 1}
	case 'v':
		return vec2{1, 1}
	case 'A':
		return vec2{2, 0}
	default:
		panic("bad")
	}
}
