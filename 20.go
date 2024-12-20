package main

func init() {
	addSolutions(20, problem20)
}

func problem20(ctx *problemContext) {
	var t racetrack
	s := ctx.scanner()
	for s.scan() {
		t.g.addRow([]byte(s.text()))
	}
	t.init()
	ctx.reportLoad()

	ctx.reportPart1(t.bestCheats(2, 100))
	ctx.reportPart2(t.bestCheats(20, 100))
}

type racetrack struct {
	g     grid[byte]
	start vec2
	end   vec2
}

func (t *racetrack) init() {
	for v, c := range t.g.all() {
		switch c {
		case 'S':
			t.start = v
			t.g.set(v, '.')
		case 'E':
			t.end = v
			t.g.set(v, '.')
		}
	}
}

func (t *racetrack) bestPaths(target vec2) map[vec2]int64 {
	type state struct {
		v    vec2
		cost int64
	}
	best := map[vec2]int64{target: 0}
	q := []state{{v: target}}
	for len(q) > 0 {
		s := SlicePop(&q)
		for _, n := range s.v.neighbors4() {
			if !t.g.contains(n) || t.g.at(n) != '.' {
				continue
			}
			if _, ok := best[n]; ok {
				continue
			}
			best[n] = s.cost + 1
			q = append(q, state{n, s.cost + 1})
		}
	}
	return best
}

func (t *racetrack) bestCheats(limit, minSavings int64) int64 {
	var (
		bestToStart = t.bestPaths(t.start)
		bestToEnd   = t.bestPaths(t.end)
		threshold   = bestToEnd[t.start] - minSavings
		result      int64
	)
	for v, c := range t.g.all() {
		if c != '.' {
			continue
		}
		costToStart := bestToStart[v]
		// Try all the possible cheats which start here.
		for n, d := range t.g.ball(v, limit+1) {
			costToEnd, ok := bestToEnd[n]
			if !ok {
				// Not a valid end of cheat.
				continue
			}
			total := costToStart + d + costToEnd
			if total <= threshold {
				result++
			}
		}
	}
	return result
}
