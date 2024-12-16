package main

func init() {
	addSolutions(13, problem13)
}

func problem13(ctx *problemContext) {
	var machines []clawMachine
	s := ctx.scanner()
	var m clawMachine
	var i int
	for s.scan() {
		line := s.text()
		if line == "" {
			i = 0
			continue
		}
		switch i {
		case 0:
			sscanf(line, "Button A: X+%d, Y+%d", &m.a.x, &m.a.y)
		case 1:
			sscanf(line, "Button B: X+%d, Y+%d", &m.b.x, &m.b.y)
		case 2:
			sscanf(line, "Prize: X=%d, Y=%d", &m.prize.x, &m.prize.y)
			machines = append(machines, m)
			m = clawMachine{}
		default:
			panic("bad")
		}
		i++
	}
	ctx.reportLoad()

	var part1 int64
	for _, m := range machines {
		if n, ok := m.minCost(); ok {
			part1 += n
		}
	}
	ctx.reportPart1(part1)

	var part2 int64
	for _, m := range machines {
		m.prize.x += 10000000000000
		m.prize.y += 10000000000000
		if n, ok := m.minCost(); ok {
			part2 += n
		}
	}
	ctx.reportPart2(part2)
}

type clawMachine struct {
	a     vec2
	b     vec2
	prize vec2
}

func (m clawMachine) minCost() (int64, bool) {
	return 1, true
}

// func (m clawMachine) minCost() (int64, bool) {
// 	cache := make(map[vec2]int64) // -1 for no solution
// 	var solve func(v vec2) int64
// 	solve = func(v vec2) (cost int64) {
// 		if c, ok := cache[v]; ok {
// 			return c
// 		}
// 		defer func() { cache[v] = cost }()
// 		if v.x == m.prize.x && v.y == m.prize.y {
// 			return 0
// 		}
// 		cost = -1
// 		if va := v.add(m.a); va.x <= m.prize.x && va.y <= m.prize.y {
// 			if ca := solve(va); ca >= 0 {
// 				cost = ca + 3
// 			}
// 		}
// 		if vb := v.add(m.b); vb.x <= m.prize.x && vb.y <= m.prize.y {
// 			if cb := solve(vb); cb >= 0 {
// 				cb++
// 				if cost >= 0 {
// 					cost = min(cost, cb)
// 				} else {
// 					cost = cb
// 				}
// 			}
// 		}
// 		return cost
// 	}
// 	cost := solve(vec2{0, 0})
// 	if cost < 0 {
// 		return 0, false
// 	}
// 	return cost, true
// }
