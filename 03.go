package main

import (
	"regexp"
)

func init() {
	addSolutions(3, problem3)
}

func problem3(ctx *problemContext) {
	input := string(ctx.readAll())
	ctx.reportLoad()

	re := regexp.MustCompile(`mul\(([0-9]{1,3}),([0-9]{1,3})\)`)
	var part1 int64
	for _, inst := range re.FindAllStringSubmatch(input, -1) {
		x := parseInt(inst[1])
		y := parseInt(inst[2])
		part1 += x * y
	}
	ctx.reportPart1(part1)

	re = regexp.MustCompile(`do\(\)|don't\(\)|` + re.String())
	var part2 int64
	enabled := true
	for _, inst := range re.FindAllStringSubmatch(input, -1) {
		switch inst[0] {
		case "do()":
			enabled = true
		case "don't()":
			enabled = false
		default:
			if enabled {
				x := parseInt(inst[1])
				y := parseInt(inst[2])
				part2 += x * y
			}
		}
	}
	ctx.reportPart2(part2)
}
