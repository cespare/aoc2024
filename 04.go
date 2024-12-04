package main

func init() {
	addSolutions(4, problem4)
}

func problem4(ctx *problemContext) {
	var g grid[byte]
	s := ctx.scanner()
	for s.scan() {
		g.addRow([]byte(s.text()))
	}
	ctx.reportLoad()

	var part1 int
	for v := range g.vecs() {
		for _, d := range box8 {
			if xmasAt(&g, v, d) {
				part1++
			}
		}
	}
	ctx.reportPart1(part1)

	var part2 int
	for y := int64(0); y < g.rows-2; y++ {
		for x := int64(0); x < g.cols-2; x++ {
			v := vec2{x, y}
			if xDashMasAt(&g, v) {
				part2++
			}
		}
	}
	ctx.reportPart2(part2)
}

func xmasAt(g *grid[byte], v, d vec2) bool {
	for _, c := range []byte("XMAS") {
		if !g.contains(v) || g.at(v) != c {
			return false
		}
		v = v.add(d)
	}
	return true
}

func xDashMasAt(g *grid[byte], v vec2) bool {
	ds := []vec2{
		{0, 0},
		{2, 0},
		{1, 1},
		{0, 2},
		{2, 2},
	}
	x := SliceMap(ds, func(d vec2) byte { return g.at(v.add(d)) })
	switch string(x) {
	case "MMASS", "MSAMS", "SMASM", "SSAMM":
		return true
	}
	return false
}
