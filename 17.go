package main

import (
	"slices"
	"strconv"
	"strings"
)

func init() {
	addSolutions(17, problem17)
}

/*

Register A: ?
Register B: 0
Register C: 0

Program: 2,4,1,1,7,5,1,5,4,3,5,5,0,3,3,0

bst 4 // B = A&0b111		  B1 = A0 & 111
bxl 1 // B ^= 0b001               B2 = B1 ^ 001 = (A0 & 111) ^ 001
cdv 5 // C = A >> B               C1 = A0 >> B2 = A0 >> ((A0 & 111) ^ 001)
bxl 5 // B ^= 0b101               B3 = B2 ^ 101 = (A0 & 111) ^ 100
bxc 3 // B ^= C                   B4 = B3 ^ C1  = ((A0 & 111) ^ 100) ^ (A0 >> ((A0 & 111) ^ 001))
out 5 // output B&0b111           output B4 & 111
adv 3 // A = A >> 3               A1 = A0 >> 3
jnz 0 // if A!=0, jump to 0

consider last round (output 0)

A1 = 0, so A0 < 8
B4 = 0 = ((A0 & 111) ^ 100) ^ (A0 >> ((A0 & 111) ^ 001))
       = (A0^100) ^ (A0 >> (A0 ^ 001))
       = (A0 >> (A0 ^ 001)) ^ 100 ^ A0
       => A0 = 4 (only solution)

second-to-last round (output 3):

A1 = 4, so A0 = 0b010xxx
Let's say Z = xxx
B4 = 3 = (Z ^ 100) ^ (A0 >> (Z ^ 001))
       = (A0 >> (Z^001)) ^ 100 ^ Z
       => Z = 7, A0 = 39 (only solution)

...

*/

func problem17(ctx *problemContext) {
	comp := parseComputer(string(ctx.readAll()))
	ctx.reportLoad()

	for !comp.step() {
	}
	ctx.reportPart1(comp.output())

	// Solve in reverse.
	candidates := []uint64{0}
	for _, inst := range slices.Backward(comp.insts) {
		// fmt.Printf("\033[01;34m>>>> a: %v\x1B[m\n", a)
		want := uint64(inst)
		// Work out the next 3 bits.
		var next []uint64
		for _, a := range candidates {
			for a0 := a << 3; a0 < a<<3+8; a0++ {
				b := a0 & 0b111
				b ^= 0b001
				c := a0 >> b
				b ^= 0b101
				b ^= c
				got := b & 0b111
				if got == want {
					next = append(next, a0)
				}
			}
		}
		if len(next) == 0 {
			panic("no solution")
		}
		candidates = next
	}
	ctx.reportPart2(slices.Min(candidates))
}

type computer struct {
	regs  [3]uint64
	insts []uint8
	ip    int

	out []uint64
}

func parseComputer(s string) *computer {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) != 5 || lines[3] != "" {
		panic("bad")
	}
	var c computer
	sscanf(lines[0], "Register A: %d", &c.regs[0])
	sscanf(lines[1], "Register B: %d", &c.regs[1])
	sscanf(lines[2], "Register C: %d", &c.regs[2])
	nums := strings.TrimPrefix(lines[4], "Program: ")
	c.insts = SliceMap(strings.Split(nums, ","), func(s string) uint8 {
		n, err := strconv.ParseUint(s, 10, 3)
		if err != nil {
			panic(err)
		}
		return uint8(n)
	})
	return &c
}

func (c *computer) step() (done bool) {
	if c.ip >= len(c.insts) {
		return true
	}
	opcode, operand := c.insts[c.ip], c.insts[c.ip+1]

	switch opcode {
	case 0: // adv
		c.regs[0] = c.regs[0] / (1 << c.combo(operand))
	case 1: // bxl
		c.regs[1] ^= uint64(operand)
	case 2: // bst
		c.regs[1] = c.combo(operand) & 0b111
	case 3: // jnz
		if c.regs[0] != 0 {
			c.ip = int(operand)
			return false
		}
	case 4: // bxc
		c.regs[1] ^= c.regs[2]
	case 5: // out
		c.out = append(c.out, c.combo(operand)&0b111)
	case 6: // bdv
		c.regs[1] = c.regs[0] / (1 << c.combo(operand))
	case 7: // cdv
		c.regs[2] = c.regs[0] / (1 << c.combo(operand))
	default:
		panic("cannot happen")
	}
	c.ip += 2
	return false
}

func (c *computer) combo(operand uint8) uint64 {
	switch operand {
	case 0, 1, 2, 3:
		return uint64(operand)
	case 4:
		return c.regs[0]
	case 5:
		return c.regs[1]
	case 6:
		return c.regs[2]
	case 7:
		panic("7")
	default:
		panic("shouldn't happen")
	}
}

func (c *computer) output() string {
	return strings.Join(
		SliceMap(c.out, func(n uint64) string {
			return strconv.FormatUint(n, 10)
		}),
		",",
	)
}
