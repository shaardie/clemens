package bitboard

import "math"

const (
	Empty     Bitboard = 0
	Universal Bitboard = math.MaxUint64

	notAFile Bitboard = 0xfefefefefefefefe
	notHFile Bitboard = 0x7f7f7f7f7f7f7f7f

	rank1 Bitboard = 0x00000000000000ff
	rank2 Bitboard = rank1 << 2
	rank3 Bitboard = rank2 << 2
	rank4 Bitboard = rank3 << 2
	rank5 Bitboard = rank4 << 2
	rank6 Bitboard = rank5 << 2
	rank7 Bitboard = rank6 << 2
	rank8 Bitboard = rank7 << 2
)
