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
		n, ok := m.minCost()
		if ok {
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
	// System of equations:
	// A * a.x + B * b.x = prize.x
	// A * a.y + B * b.y = prize.y
	//
	// Rename:
	//
	// A*Xa + B*Xb = Xp
	// A*Ya + B*Yb = Yp
	//
	// => A = (Xp - B*Xb) / Xa
	//
	// => ((Xp - B*Xb)/Xa) * Ya + B*Yb = Yp
	//    (Ya*Xp - Ya*B*Xb)/Xa + Xa*B*Yb/Xa = Yp
	//    Ya*Xp - Ya*B*Xb + Xa*B*Yb = Yp*Xa
	//    Xa*B*Yb - Ya*B*Xb = Yp*Xa - Ya*Xp
	//    B(Xa*Yb - Ya*Xb) = Yp*Xa - Ya*Xp
	//    B = (Yp*Xa - Xp*Ya) / (Xa*Yb - Ya*Xb)
	//
	// => A = (Yp*Xb - Xp*Yb) / (Xb*Ya - Yb*Xa) (by symmetry)

	numA := m.prize.y*m.b.x - m.prize.x*m.b.y
	denomA := m.b.x*m.a.y - m.b.y*m.a.x
	if denomA == 0 || numA%denomA != 0 {
		return 0, false
	}
	a := numA / denomA

	numB := m.prize.y*m.a.x - m.prize.x*m.a.y
	denomB := m.a.x*m.b.y - m.a.y*m.b.x
	if denomB == 0 || numB%denomB != 0 {
		return 0, false
	}
	b := numB / denomB

	return a*3 + b, true
}
