package bitboard

import "math/bits"

func LeastSignificantOneBit(b Bitboard) uint8 {
	if b == Empty {
		panic("Bitboard empty")
	}
	return uint8(bits.TrailingZeros64(uint64(b)))
}

func MostSignificantOneBit(b Bitboard) uint8 {
	if b == Empty {
		panic("Bitboard empty")
	}
	return uint8(bits.LeadingZeros64(uint64(b)))
}
