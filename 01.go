package main

import "slices"

func init() {
	addSolutions(1, problem1)
}

func problem1(ctx *problemContext) {
	var list0, list1 []int64
	s := ctx.scanner()
	for s.scan() {
		parts := s.fields()
		list0 = append(list0, parseInt(parts[0]))
		list1 = append(list1, parseInt(parts[1]))
	}
	ctx.reportLoad()

	sorted0 := slices.Sorted(slices.Values(list0))
	sorted1 := slices.Sorted(slices.Values(list1))
	var part1 int64
	for i, n0 := range sorted0 {
		n1 := sorted1[i]
		part1 += abs(n1 - n0)
	}
	ctx.reportPart1(part1)

	counts := make(map[int64]int64)
	for _, n1 := range list1 {
		counts[n1]++
	}
	var part2 int64
	for _, n0 := range list0 {
		part2 += n0 * counts[n0]
	}
	ctx.reportPart2(part2)
}
