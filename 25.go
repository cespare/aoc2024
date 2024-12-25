package main

import (
	"strings"
)

func init() {
	addSolutions(25, problem25)
}

func problem25(ctx *problemContext) {
	var locks [][]int
	var keys [][]int
	for _, chunk := range strings.Split(string(ctx.readAll()), "\n\n") {
		pins, isLock := parseLockOrKey(chunk)
		if isLock {
			locks = append(locks, pins)
		} else {
			keys = append(keys, pins)
		}
	}
	ctx.reportLoad()

	var part1 int64
	for _, key := range keys {
		for _, lock := range locks {
			if keyFitsLock(key, lock) {
				part1++
			}
		}
	}
	ctx.reportPart1(part1)
}

func parseLockOrKey(s string) (pins []int, isLock bool) {
	var g grid[byte]
	for i, row := range strings.Split(strings.TrimSpace(s), "\n") {
		if i == 0 && row == "#####" {
			isLock = true
		}
		g.addRow([]byte(row))
	}
	pins = make([]int, 5)
	for x := range pins {
		for y := 0; y < 7; y++ {
			c := g.at(vec2{int64(x), int64(y)})
			if isLock {
				if c != '#' {
					pins[x] = y - 1
					break
				}
			} else {
				if c != '.' {
					pins[x] = 6 - y
					break
				}
			}
		}
	}
	return pins, isLock
}

func keyFitsLock(key, lock []int) bool {
	for i, k := range key {
		p := lock[i]
		if k+p > 5 {
			return false
		}
	}
	return true
}
