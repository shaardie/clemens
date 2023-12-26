package pawn

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	WHITE_IDX = iota
	BLACK_IDX
	COLOR_NUMBER
)

var (
	attackTable [COLOR_NUMBER][types.SQUARE_NUMBER]bitboard.Bitboard
)

// init initializes the attack table for knights for all squares
func init() {
	for square := types.SQUARE_A1; square < types.SQUARE_NUMBER; square++ {
		attackTable[WHITE_IDX][square] = AttacksByBitboard(
			types.WHITE, bitboard.BitBySquares(square),
		)
		attackTable[BLACK_IDX][square] = AttacksByBitboard(
			types.BLACK, bitboard.BitBySquares(square),
		)
	}
}

// AttacksBySquare returns the attacks for a given square.
// This is done by lookup.
func AttacksBySquare(c types.Color, square int) bitboard.Bitboard {
	return attackTable[c][square]
}

// AttacksByBitboard calculates all attacks for the given bitboard
func AttacksByBitboard(c types.Color, pawns bitboard.Bitboard) bitboard.Bitboard {
	switch c {
	case types.WHITE:
		return bitboard.NorthEastOne(pawns) | bitboard.NorthWestOne(pawns)
	case types.BLACK:
		return bitboard.SouthEastOne(pawns) | bitboard.SouthWestOne(pawns)
	}
	panic("unknown color")
}

type BlackPawns bitboard.Bitboard
type WhitePawns bitboard.Bitboard

func (p WhitePawns) SinglePushTargets(emptySquares bitboard.Bitboard) bitboard.Bitboard {
	return bitboard.NorthOne(bitboard.Bitboard(p)) & emptySquares
}

func (p WhitePawns) DoublePushTargets(emptySquares bitboard.Bitboard) bitboard.Bitboard {
	// Mandatory condition that single push is possible
	singlePushTargets := p.SinglePushTargets(emptySquares)
	// White Double Push only possible on empty fields on rank 4
	return bitboard.SouthOne(singlePushTargets) & emptySquares & bitboard.RankMask4
}

func (p BlackPawns) SinglePushTargets(emptySquares bitboard.Bitboard) bitboard.Bitboard {
	return bitboard.SouthOne(bitboard.Bitboard(p)) & emptySquares
}

func (p BlackPawns) DoublePushTargets(emptySquares bitboard.Bitboard) bitboard.Bitboard {
	// Mandatory condition that single push is possible
	singlePushTargets := p.SinglePushTargets(emptySquares)
	// Black Double Push only possible on empty fields on rank 5
	return bitboard.SouthOne(singlePushTargets) & emptySquares & bitboard.RankMask5
}
