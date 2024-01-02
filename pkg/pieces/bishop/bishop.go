package bishop

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/magic"
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/pieces/utils"
	"github.com/shaardie/clemens/pkg/types"
)

var magics [types.SQUARE_NUMBER]magic.Magic

func init() {
	magics = magic.Init(AttacksByBitboard)
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

// AttacksByBitboard calculates the AttacksByBitboard of the bishop for the given square and occupation
func AttacksByBitboard(square int, occupied bitboard.Bitboard) bitboard.Bitboard {
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

// GenerateMoves generates all moves for all squares to all destinations by a given occupation
func GenerateMoves(squares, occupied, destinations bitboard.Bitboard) []move.Move {
	return utils.GenerateMoves(
		squares,
		occupied,
		destinations,
		AttacksByBitboard,
	)
}
