package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
)

func init() {
	addSolutions(24, problem24)
}

/*

x00 XOR y00 -> z00
x00 AND y00 -> gwq // carry from z00

x01 XOR y01 -> qgt
gwq XOR qgt -> z01
gwq AND qgt -> cgt
x01 AND y01 -> pvw
pvw OR  cgt -> mct // carry from z01


x02 XOR y02 -> wvk
mct XOR wvk -> z02
mct AND wvk -> mtb
x02 AND y02 -> pcq
pcq OR mtb -> vgc // carry from z02

...


swaps:

kfp,hbs
dhq,z18
pdg,z22
jcp,z27

dhq,hbs,jcp,kfp,pdg,z18,z22,z27

*/

func problem24(ctx *problemContext) {
	var d logicDevice
	s := ctx.scanner()
	phase := 0
	for s.scan() {
		if s.text() == "" {
			phase++
			continue
		}
		switch phase {
		case 0:
			wire, bitStr, ok := strings.Cut(s.text(), ": ")
			if !ok {
				panic("bad")
			}
			bit := parseUint(bitStr)
			if d.init == nil {
				d.init = make(map[string]uint64)
			}
			d.init[wire] = bit
		case 1:
			d.gates = append(d.gates, parseLogicGate(s.text()))
		default:
			panic("bad")
		}
	}
	ctx.reportLoad()

	ctx.reportPart1(d.eval())

	// Part 2 is solved semi-manually by incrementally running validate and
	// fixing each swap as it comes up.
	//
	// if err := d.validate(); err != nil {
	// 	panic(err)
	// }
}

func (d *logicDevice) validate() error {
	gates := make(map[[3]string]logicGate)
	gatesByOut := make(map[string]logicGate)
	for _, g := range d.gates {
		gates[[3]string{g.left, g.op, g.right}] = g
		gatesByOut[g.out] = g
	}
	find := func(left, op, right string) string {
		g, ok := gates[[3]string{left, op, right}]
		if !ok {
			g, ok = gates[[3]string{right, op, left}]
			if !ok {
				panic(fmt.Sprintf(
					"%s %s %s not found",
					left, op, right,
				))
			}
		}
		fmt.Println(g)
		return g.out
	}
	carry := find("x00", "AND", "y00")
	for pos := 1; pos <= 44; pos++ {
		log.Printf("Checking %d", pos)
		x := fmt.Sprintf("x%02d", pos)
		y := fmt.Sprintf("y%02d", pos)
		tmp0 := find(x, "XOR", y)
		z := find(carry, "XOR", tmp0)
		if !strings.HasPrefix(z, "z") {
			return fmt.Errorf(
				"%s XOR %s -> %s; want z%02d",
				carry, tmp0, z, pos,
			)
		}
		tmp1 := find(carry, "AND", tmp0)
		tmp2 := find(x, "AND", y)

		carry = find(tmp1, "OR", tmp2)
	}
	return nil
}

type logicDevice struct {
	init  map[string]uint64
	gates []logicGate
}

type logicGate struct {
	left  string
	op    string
	right string
	out   string
}

func (g logicGate) String() string {
	return fmt.Sprintf("%s %s %s -> %s", g.left, g.op, g.right, g.out)
}

func parseLogicGate(s string) logicGate {
	var g logicGate
	sscanf(s, "%s %s %s -> %s", &g.left, &g.op, &g.right, &g.out)
	if g.left > g.right {
		g.left, g.right = g.right, g.left
	}
	return g
}

func (d *logicDevice) eval() uint64 {
	chs := make(map[string]*lazyVal)
	for wire, bit := range d.init {
		v := newLazyVal()
		v.set(bit)
		chs[wire] = v
	}
	var wg sync.WaitGroup
	var result atomic.Uint64
	var fns []func()
	for _, g := range d.gates {
		left, ok := chs[g.left]
		if !ok {
			left = newLazyVal()
			chs[g.left] = left
		}
		right, ok := chs[g.right]
		if !ok {
			right = newLazyVal()
			chs[g.right] = right
		}
		out, ok := chs[g.out]
		if !ok {
			out = newLazyVal()
			chs[g.out] = out
		}
		fn := func() {
			var r uint64
			switch g.op {
			case "AND":
				r = left.get() & right.get()
			case "OR":
				r = left.get() | right.get()
			case "XOR":
				r = left.get() ^ right.get()
			default:
				panic("bad")
			}
			out.set(r)

			if ns, ok := strings.CutPrefix(g.out, "z"); ok {
				if r == 1 {
					result.Or(1 << parseUint(ns))
				}
			}
		}
		fns = append(fns, fn)
	}
	for _, fn := range fns {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn()
		}()
	}
	wg.Wait()
	return result.Load()
}

type lazyVal struct {
	done chan struct{}
	val  uint64
}

func newLazyVal() *lazyVal {
	return &lazyVal{done: make(chan struct{})}
}

func (v *lazyVal) get() uint64 {
	<-v.done
	return v.val
}

func (v *lazyVal) set(n uint64) {
	v.val = n
	close(v.done)
}
