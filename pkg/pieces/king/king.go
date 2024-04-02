package king

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

var attackTable [types.SQUARE_NUMBER]bitboard.Bitboard

// init initializes the attack table for knights for all squares
func init() {
	for square := types.SQUARE_A1; square < types.SQUARE_NUMBER; square++ {
		attackTable[square] = attacks(
			bitboard.BitBySquares(square))
	}
}

// AttacksBySquare returns the attacks for a given square.
// This is done by lookup.
func AttacksBySquare(square uint8) bitboard.Bitboard {
	return attackTable[square]
}

func attacks(king bitboard.Bitboard) bitboard.Bitboard {
	// Attacks to West and East
	attacks := bitboard.WestOne(king) | bitboard.EastOne(king)

	// Use attacks to also calculate attacks in north and south,
	// so it also includes northeast, etc.
	kings := king | attacks
	attacks |= bitboard.NorthOne(kings) | bitboard.SouthOne(kings)

	return attacks
}
