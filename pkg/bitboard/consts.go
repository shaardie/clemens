package bitboard

import "math"

const (
	SQUARE_A1 int = iota
	SQUARE_B1
	SQUARE_C1
	SQUARE_D1
	SQUARE_E1
	SQUARE_F1
	SQUARE_G1
	SQUARE_H1
	SQUARE_A2
	SQUARE_B2
	SQUARE_C2
	SQUARE_D2
	SQUARE_E2
	SQUARE_F2
	SQUARE_G2
	SQUARE_H2
	SQUARE_A3
	SQUARE_B3
	SQUARE_C3
	SQUARE_D3
	SQUARE_E3
	SQUARE_F3
	SQUARE_G3
	SQUARE_H3
	SQUARE_A4
	SQUARE_B4
	SQUARE_C4
	SQUARE_D4
	SQUARE_E4
	SQUARE_F4
	SQUARE_G4
	SQUARE_H4
	SQUARE_A5
	SQUARE_B5
	SQUARE_C5
	SQUARE_D5
	SQUARE_E5
	SQUARE_F5
	SQUARE_G5
	SQUARE_H5
	SQUARE_A6
	SQUARE_B6
	SQUARE_C6
	SQUARE_D6
	SQUARE_E6
	SQUARE_F6
	SQUARE_G6
	SQUARE_H6
	SQUARE_A7
	SQUARE_B7
	SQUARE_C7
	SQUARE_D7
	SQUARE_E7
	SQUARE_F7
	SQUARE_G7
	SQUARE_H7
	SQUARE_A8
	SQUARE_B8
	SQUARE_C8
	SQUARE_D8
	SQUARE_E8
	SQUARE_F8
	SQUARE_G8
	SQUARE_H8
)

const (
	Rank1 int = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
	RankNumber
)

const (
	FileA int = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
	FileNumber
)

const (
	Empty     Bitboard = 0
	One       Bitboard = 1
	Universal Bitboard = math.MaxUint64

	RankMask1 Bitboard = 0x00000000000000ff
	RankMask2 Bitboard = RankMask1 << 8
	RankMask3 Bitboard = RankMask2 << 8
	RankMask4 Bitboard = RankMask3 << 8
	RankMask5 Bitboard = RankMask4 << 8
	RankMask6 Bitboard = RankMask5 << 8
	RankMask7 Bitboard = RankMask6 << 8
	RankMask8 Bitboard = RankMask7 << 8

	FileMaskA Bitboard = 0x0101010101010101
	FileMaskB Bitboard = FileMaskA << 1
	FileMaskC Bitboard = FileMaskB << 1
	FileMaskD Bitboard = FileMaskC << 1
	FileMaskE Bitboard = FileMaskD << 1
	FileMaskF Bitboard = FileMaskE << 1
	FileMaskG Bitboard = FileMaskF << 1
	FileMaskH Bitboard = FileMaskG << 1

	notAFile Bitboard = ^FileMaskA
	notHFile Bitboard = ^FileMaskH
)
