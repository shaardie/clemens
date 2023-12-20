package bitboard

import "math"

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
