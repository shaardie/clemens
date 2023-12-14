package bitboard

import "math/bits"

func LeastSignificantOneBit(b Bitboard) int {
	if b == Empty {
		panic("Bitboard empty")
	}
	return bits.TrailingZeros64(uint64(b))
}

func MostSignificantOneBit(b Bitboard) int {
	if b == Empty {
		panic("Bitboard empty")
	}
	return bits.LeadingZeros64(uint64(b))
}
