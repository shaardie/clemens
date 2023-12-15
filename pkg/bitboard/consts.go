package bitboard

import "math"

const (
	rank1 int = iota
	rank2
	rank3
	rank4
	rank5
	rank6
	rank7
	rank8
)

const (
	file1 int = iota
	file2
	file3
	file4
	file5
	file6
	file7
	file8
)

const (
	Empty     Bitboard = 0
	Universal Bitboard = math.MaxUint64

	notAFile Bitboard = 0xfefefefefefefefe
	notHFile Bitboard = 0x7f7f7f7f7f7f7f7f

	FullRank1 Bitboard = 0x00000000000000ff
	FullRank2 Bitboard = FullRank1 << 2
	FullRank3 Bitboard = FullRank2 << 2
	FullRank4 Bitboard = FullRank3 << 2
	FullRank5 Bitboard = FullRank4 << 2
	FullRank6 Bitboard = FullRank5 << 2
	FullRank7 Bitboard = FullRank6 << 2
	FullRank8 Bitboard = FullRank7 << 2
)
