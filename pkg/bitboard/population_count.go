package bitboard

import "math/bits"

func (b Bitboard) PopulationCount() int {
	return bits.OnesCount64(uint64(b))
}
