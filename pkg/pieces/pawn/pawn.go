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
func AttacksBySquare(c types.Color, square uint8) bitboard.Bitboard {
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

func Pushes(c types.Color, pawns, occupied bitboard.Bitboard) bitboard.Bitboard {
	return singlePushTargets(c, pawns, occupied) | doublePushTargets(c, pawns, occupied)
}

func PushesBySquare(c types.Color, square uint8, occupied bitboard.Bitboard) bitboard.Bitboard {
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

// NumberOfIsolanis returns the number of isolanis from a bitboard of pawns.
func NumberOfIsolanis(pawns bitboard.Bitboard) int {
	return Isolanis(pawns).PopulationCount()
}

// Isolanis returns the isolanis from a bitboard of pawns.
func Isolanis(pawns bitboard.Bitboard) bitboard.Bitboard {
	fileFill := bitboard.FileFill(pawns)
	westAttackFileFill := bitboard.WestOne(fileFill)
	eastAttackFileFill := bitboard.EastOne(fileFill)
	return pawns & ^westAttackFileFill & ^eastAttackFileFill
}

// NumberOfDoubled calculates the number of doubled pawns,
// we do not care that tripple pawns a counted twice, etc.
func NumberOfDoubled(pawns bitboard.Bitboard) int {
	return Doubled(pawns).PopulationCount()
}

// Doubled returns a bitboard with the doubled pawns
// we do not care that tripple pawns a counted twice, etc.
func Doubled(pawns bitboard.Bitboard) bitboard.Bitboard {
	return bitboard.NorthOne(bitboard.NorthFill(pawns)) & pawns
}

func Passed(color types.Color, whitePawns, blackPawns bitboard.Bitboard) bitboard.Bitboard {
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

func NumberOfPassed(color types.Color, whitePawns, blackPawns bitboard.Bitboard) int {
	return Passed(color, whitePawns, blackPawns).PopulationCount()
}

func Backwards(color types.Color, whitePawns, blackPawns bitboard.Bitboard) bitboard.Bitboard {
	if color == types.WHITE {
		stops := bitboard.NorthOne(whitePawns)
		wAttackSpans := bitboard.NorthWestOne(bitboard.NorthFill(whitePawns)) |
			bitboard.NorthEastOne(bitboard.NorthFill(whitePawns))
		bAttacks := attacks(types.BLACK, blackPawns)
		return bitboard.SouthOne(stops & bAttacks & ^wAttackSpans)
	}

	stops := bitboard.SouthOne(blackPawns)
	bAttackSpans := bitboard.SouthWestOne(bitboard.SouthFill(blackPawns)) |
		bitboard.SouthEastOne(bitboard.SouthFill(blackPawns))
	wAttacks := attacks(types.WHITE, whitePawns)
	return bitboard.NorthOne(stops & wAttacks & ^bAttackSpans)
}

func NumberOfBackwards(color types.Color, whitePawns, blackPawns bitboard.Bitboard) int {
	return Backwards(color, whitePawns, blackPawns).PopulationCount()
}

func Supported(color types.Color, pawns bitboard.Bitboard) bitboard.Bitboard {
	return attacks(color, pawns) & pawns
}

func NumberOfSupported(color types.Color, pawns bitboard.Bitboard) int {
	return Supported(color, pawns).PopulationCount()
}

func Phalanx(pawns bitboard.Bitboard) bitboard.Bitboard {
	return pawns & (bitboard.WestOne(pawns) | bitboard.EastOne(pawns))
}

func Opposed(color types.Color, whitePawns, blackPawns bitboard.Bitboard) bitboard.Bitboard {
	if color == types.WHITE {
		return whitePawns & bitboard.SouthFill(blackPawns)
	}
	return blackPawns & bitboard.NorthFill(whitePawns)
}
