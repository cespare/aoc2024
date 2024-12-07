package main

import (
	"strconv"
	"strings"
)

func init() {
	addSolutions(7, problem7)
}

func problem7(ctx *problemContext) {
	var equations []calibEquation
	s := ctx.scanner()
	for s.scan() {
		equations = append(equations, parseCalibEquation(s.text()))
	}
	ctx.reportLoad()

	var part1 int64
	for _, eq := range equations {
		if eq.ok() {
			part1 += eq.testVal
		}
	}
	ctx.reportPart1(part1)

	var part2 int64
	for _, eq := range equations {
		if eq.ok2() {
			part2 += eq.testVal
		}
	}
	ctx.reportPart2(part2)
}

type calibEquation struct {
	testVal int64
	args    []int64
}

func parseCalibEquation(s string) calibEquation {
	head, tail, ok := strings.Cut(s, ": ")
	if !ok {
		panic("bad")
	}
	eq := calibEquation{
		testVal: parseInt(head),
		args:    SliceMap(strings.Fields(tail), parseInt),
	}
	for _, n := range eq.args {
		if n <= 0 {
			panic(n)
		}
	}
	return eq
}

func (e calibEquation) ok() bool {
	return eqCanEqual(0, e.args, e.testVal)
}

func eqCanEqual(total int64, rem []int64, targ int64) bool {
	if len(rem) == 0 {
		return total == targ
	}
	next, rem := rem[0], rem[1:]
	if x := total + next; x <= targ {
		if eqCanEqual(x, rem, targ) {
			return true
		}
	}
	if x := total * next; x <= targ {
		if eqCanEqual(x, rem, targ) {
			return true
		}
	}
	return false
}

func (e calibEquation) ok2() bool {
	return eqCanEqual2(0, e.args, e.testVal)
}

func eqCanEqual2(total int64, rem []int64, targ int64) bool {
	if len(rem) == 0 {
		return total == targ
	}
	next, rem := rem[0], rem[1:]
	if x := total + next; x <= targ {
		if eqCanEqual2(x, rem, targ) {
			return true
		}
	}
	if x := total * next; x <= targ {
		if eqCanEqual2(x, rem, targ) {
			return true
		}
	}
	x := parseInt(strconv.FormatInt(total, 10) + strconv.FormatInt(next, 10))
	if x <= targ {
		if eqCanEqual2(x, rem, targ) {
			return true
		}
	}
	return false
}
