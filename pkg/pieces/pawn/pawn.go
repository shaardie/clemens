package pawn

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	attackTable [types.COLOR_NUMBER][types.SQUARE_NUMBER]bitboard.Bitboard
)

// init initializes the attack table for knights for all squares
func init() {
	for square := types.SQUARE_A1; square < types.SQUARE_NUMBER; square++ {
		attackTable[types.WHITE][square] = attacks(
			types.WHITE, bitboard.BitBySquares(square),
		)
		attackTable[types.BLACK][square] = attacks(
			types.BLACK, bitboard.BitBySquares(square),
		)
	}
}

// AttacksBySquare returns the attacks for a given square.
// This is done by lookup.
func AttacksBySquare(c types.Color, square int) bitboard.Bitboard {
	return attackTable[c][square]
}

// attacks calculates all attacks for the given bitboard
func attacks(c types.Color, pawns bitboard.Bitboard) bitboard.Bitboard {
	switch c {
	case types.WHITE:
		return bitboard.NorthEastOne(pawns) | bitboard.NorthWestOne(pawns)
	case types.BLACK:
		return bitboard.SouthEastOne(pawns) | bitboard.SouthWestOne(pawns)
	}
	panic("unknown color")
}

func PushesBySquare(c types.Color, square int, occupied bitboard.Bitboard) bitboard.Bitboard {
	pawn := bitboard.BitBySquares(square)
	return singlePushTargets(c, pawn, occupied) | doublePushTargets(c, pawn, occupied)
}

func singlePushTargets(c types.Color, pawns, occupied bitboard.Bitboard) bitboard.Bitboard {
	switch c {
	case types.WHITE:
		return bitboard.NorthOne(pawns) & ^occupied
	case types.BLACK:
		return bitboard.SouthOne(pawns) & ^occupied
	}
	panic("unknown color")
}

func doublePushTargets(c types.Color, pawns, occupied bitboard.Bitboard) bitboard.Bitboard {
	switch c {
	case types.WHITE:
		// Mandatory condition that single push is possible
		singlePushTargets := singlePushTargets(c, pawns, occupied)
		// White Double Push only possible on empty fields on rank 4
		return bitboard.NorthOne(singlePushTargets) & ^occupied & bitboard.RankMask4
	case types.BLACK:
		// Mandatory condition that single push is possible
		singlePushTargets := singlePushTargets(c, pawns, occupied)
		// Black Double Push only possible on empty fields on rank 5
		return bitboard.SouthOne(singlePushTargets) & ^occupied & bitboard.RankMask5
	}
	panic("unknown color")
}
