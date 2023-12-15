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

	fullRank1 Bitboard = 0x00000000000000ff
	fullRank2 Bitboard = fullRank1 << 2
	fullRank3 Bitboard = fullRank2 << 2
	fullRank4 Bitboard = fullRank3 << 2
	fullRank5 Bitboard = fullRank4 << 2
	fullRank6 Bitboard = fullRank5 << 2
	fullRank7 Bitboard = fullRank6 << 2
	fullRank8 Bitboard = fullRank7 << 2
)
