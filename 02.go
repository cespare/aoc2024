package main

import (
	"cmp"
	"slices"
)

func init() {
	addSolutions(2, problem2)
}

func problem2(ctx *problemContext) {
	var levels [][]int64
	s := ctx.scanner()
	for s.scan() {
		var level []int64
		for _, field := range s.fields() {
			level = append(level, parseInt(field))
		}
		levels = append(levels, level)
	}
	ctx.reportLoad()

	var part1 int
	for _, level := range levels {
		if checkLevel1(level) {
			part1++
		}
	}
	ctx.reportPart1(part1)

	var part2 int
	for _, level := range levels {
		if checkLevel2(level) {
			part2++
		}
	}
	ctx.reportPart2(part2)
}

func levelInc(level []int64) (inc, ok bool) {
	var numInc, numDec int
	for i := 1; i < len(level); i++ {
		switch cmp.Compare(level[i], level[i-1]) {
		case -1:
			numDec++
		case 1:
			numInc++
		}
	}
	if numInc > 1 {
		if numDec > 1 {
			return false, false
		}
		return true, true
	}
	return false, true
}

func findLevelError(level []int64, inc bool) int {
	for i := 1; i < len(level); i++ {
		d := level[i] - level[i-1]
		if !inc {
			d = -d
		}
		if d < 1 || d > 3 {
			return i
		}
	}
	return -1
}

func checkLevel1(level []int64) bool {
	inc, ok := levelInc(level)
	if !ok {
		return false
	}
	return findLevelError(level, inc) < 0
}

func checkLevel2(level []int64) bool {
	inc, ok := levelInc(level)
	if !ok {
		return false
	}
	i := findLevelError(level, inc)
	if i < 0 {
		return true
	}
	if findLevelError(slices.Delete(slices.Clone(level), i-1, i), inc) < 0 {
		return true
	}
	if findLevelError(slices.Delete(slices.Clone(level), i, i+1), inc) < 0 {
		return true
	}
	return false
}
