package main

import (
	"strings"
)

func init() {
	addSolutions(19, problem19)
}

func problem19(ctx *problemContext) {
	var o onsen
	s := ctx.scanner()
	var stanza int
	for s.scan() {
		line := s.text()
		if line == "" {
			stanza++
			continue
		}
		switch stanza {
		case 0:
			o.towels = strings.Split(line, ", ")
		case 1:
			o.designs = append(o.designs, line)
		default:
			panic("bad")
		}
	}
	o.init()
	ctx.reportLoad()

	ctx.reportPart1(o.countPossible())
	ctx.reportPart2(o.countWays())
}

type onsen struct {
	towels  []string
	designs []string

	byFirstLetter map[byte][]string
}

func (o *onsen) init() {
	o.byFirstLetter = make(map[byte][]string)
	for _, t := range o.towels {
		o.byFirstLetter[t[0]] = append(o.byFirstLetter[t[0]], t)
	}
}

func (o *onsen) countPossible() int {
	var possible int
	for _, d := range o.designs {
		if o.ways(d, make(map[string]int64)) > 0 {
			possible++
		}
	}
	return possible
}

func (o *onsen) countWays() int64 {
	var ways int64
	for _, d := range o.designs {
		ways += o.ways(d, make(map[string]int64))
	}
	return ways
}

func (o *onsen) ways(d string, memo map[string]int64) int64 {
	if d == "" {
		return 1
	}
	if n, ok := memo[d]; ok {
		return n
	}
	var ways int64
	defer func() { memo[d] = ways }()
	c := d[0]
	for _, t := range o.byFirstLetter[c] {
		rest, ok := strings.CutPrefix(d, t)
		if !ok {
			continue
		}
		ways += o.ways(rest, memo)
	}
	return ways
}
