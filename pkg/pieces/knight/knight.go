package knight

import (
	"github.com/shaardie/clemens/pkg/bitboard"
)

var attackTable [bitboard.SQUARE_NUMBER]bitboard.Bitboard

// init initializes the attack table for knights for all squares
func init() {
	for square := bitboard.SQUARE_A1; square < bitboard.SQUARE_NUMBER; square++ {
		attackTable[square] = AttacksByBitboard(
			bitboard.BitBySquares(square))
	}
}

// AttacksBySquare returns the attacks for a given square.
// This is done by lookup.
func AttacksBySquare(square int) bitboard.Bitboard {
	return attackTable[square]
}

// AttacksByBitboard calculates all attacks for knights on the given bitboard
func AttacksByBitboard(knights bitboard.Bitboard) bitboard.Bitboard {
	// Attacks 1 west or east and 2 north or south
	east := bitboard.EastOne(knights)
	west := bitboard.WestOne(knights)
	attacks := (west|east)<<16 | (west|east)>>16

	// Attacks 2 west or east and 1 north or south
	east = bitboard.EastOne(east)
	west = bitboard.WestOne(west)
	attacks |= bitboard.NorthOne(west|east) | bitboard.SouthOne(west|east)

	return attacks
}
