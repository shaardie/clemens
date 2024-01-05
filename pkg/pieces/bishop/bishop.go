package bishop

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/magic"
	"github.com/shaardie/clemens/pkg/pieces/utils"
	"github.com/shaardie/clemens/pkg/types"
)

var magics [types.SQUARE_NUMBER]magic.Magic

func init() {
	magics = magic.Init(attacks)
}

// AttacksBySquare returns the attacks for a given square.
// This is done by magic lookup.
func AttacksBySquare(square int, occupied bitboard.Bitboard) bitboard.Bitboard {
	// Get magic for square
	m := magics[square]

	// Calucate index of the occupation
	idx := m.Index(occupied)

	// Return attacks for the given occupation
	return m.Attacks[idx]
}

// attacks calculates the attacks of the bishop for the given square and occupation
func attacks(square int, occupied bitboard.Bitboard) bitboard.Bitboard {
	return utils.SlidingAttacks(
		square,
		[]func(bitboard.Bitboard) bitboard.Bitboard{
			bitboard.NorthEastOne,
			bitboard.NorthWestOne,
			bitboard.SouthEastOne,
			bitboard.SouthWestOne,
		},
		occupied,
	)
}
