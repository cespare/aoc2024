package main

import (
	"iter"
)

func init() {
	addSolutions(22, problem22)
}

func problem22(ctx *problemContext) {
	var nums []uint64
	s := ctx.scanner()
	for s.scan() {
		nums = append(nums, parseUint(s.text()))
	}
	ctx.reportLoad()

	const seqLength = 2000

	priceSeqs := make([][]int, len(nums))
	var part1 uint64
	for i, n := range nums {
		prices := make([]int, seqLength+1)
		prices[0] = int(n % 10)
		for j := range seqLength {
			n = nextSecretNumber(n)
			prices[j+1] = int(n % 10)
		}
		priceSeqs[i] = prices
		part1 += n
	}
	ctx.reportPart1(part1)

	seqIndices := make(map[seq4]map[int]int) // seq -> price idx -> first appearance
	for seqIdx, prices := range priceSeqs {
		for s, i := range allSeq4s(prices) {
			indices, ok := seqIndices[s]
			if !ok {
				indices = make(map[int]int)
				seqIndices[s] = indices
			}
			if _, ok := indices[seqIdx]; ok {
				continue
			}
			indices[seqIdx] = i
		}
	}
	var mostBananas int
	for _, indices := range seqIndices {
		var bananas int
		for seqIdx, i := range indices {
			bananas += priceSeqs[seqIdx][i]
		}
		mostBananas = max(mostBananas, bananas)
	}
	ctx.reportPart2(mostBananas)
}

func nextSecretNumber(n uint64) uint64 {
	n = ((n << 6) ^ n) & pruneMask
	n = ((n >> 5) ^ n) & pruneMask
	n = ((n << 11) ^ n) & pruneMask
	return n
}

const pruneMask = uint64(1)<<24 - 1

type seq4 [4]int

func allSeq4s(prices []int) iter.Seq2[seq4, int] {
	return func(yield func(seq4, int) bool) {
		for i := 1; i < len(prices)-4; i++ {
			var s seq4
			for j := range s {
				s[j] = prices[i+j] - prices[i+j-1]
			}
			if !yield(s, i+3) {
				return
			}
		}
	}
}
