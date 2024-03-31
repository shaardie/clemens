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

// NumberOfIsolanis returns tthe number of isolanis from a bitboard of pawns.
func NumberOfIsolanis(pawns bitboard.Bitboard) int {
	fileFill := bitboard.FileFill(pawns)
	westAttackFileFill := bitboard.WestOne(fileFill)
	eastAttackFileFill := bitboard.EastOne(fileFill)
	r := pawns & ^westAttackFileFill & ^eastAttackFileFill
	return r.PopulationCount()
}

// NumberOfDoubledPawns calculates the number of doubled pawns,
// we do not care that tripple pawns a counted twice, etc.
func NumberOfDoubledPawns(pawns bitboard.Bitboard) int {
	return (bitboard.NorthOne(bitboard.NorthFill(pawns)) & pawns).PopulationCount()
}

func PassedPawns(color types.Color, whitePawns, blackPawns bitboard.Bitboard) bitboard.Bitboard {
	// White
	if color == types.WHITE {
		allFrontSpans := bitboard.SouthFill(blackPawns)
		allFrontSpans |= bitboard.EastOne(allFrontSpans) | bitboard.WestOne(allFrontSpans)
		return whitePawns & ^allFrontSpans
	}

	// Black
	allFrontSpans := bitboard.NorthFill(whitePawns)
	allFrontSpans |= bitboard.EastOne(allFrontSpans) | bitboard.WestOne(allFrontSpans)
	return blackPawns & ^allFrontSpans
}
