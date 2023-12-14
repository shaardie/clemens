package bitboard

import "math/bits"

func PopulationCount(b Bitboard) int {
	return bits.OnesCount64(uint64(b))
}
