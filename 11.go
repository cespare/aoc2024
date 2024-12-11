package main

import (
	"strings"
)

func init() {
	addSolutions(11, problem11)
}

func problem11(ctx *problemContext) {
	s := make(stones)
	for _, field := range strings.Fields(string(ctx.readAll())) {
		s[parseInt(field)]++
	}
	ctx.reportLoad()

	for range 25 {
		s = s.step()
	}
	ctx.reportPart1(s.count())

	for range 50 {
		s = s.step()
	}
	ctx.reportPart2(s.count())
}

type stones map[int64]int64

func (s stones) step() stones {
	next := make(stones)
outer:
	for n, c := range s {
		if n == 0 {
			next[1] += c
			continue
		}
		var digits int
		half := int64(1)
		for pow10 := int64(10); pow10 < 1e18; pow10 *= 10 {
			digits++
			if digits%2 == 0 {
				half *= 10
			}
			if n >= pow10 {
				continue
			}
			if digits%2 == 0 {
				next[n/half] += c
				next[n%half] += c
			} else {
				next[n*2024] += c
			}
			continue outer
		}
		panic("too big")
	}
	return next
}

func (s stones) count() int64 {
	var total int64
	for _, c := range s {
		total += c
	}
	return total
}
