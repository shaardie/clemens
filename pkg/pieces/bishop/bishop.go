package bishop

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/magic"
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/pieces/utils"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	table  []bitboard.Bitboard
	magics [types.SQUARE_NUMBER]magic.Magic
)

func init() {
	table, magics = magic.Init(attacks)
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

func GenerateMoves(sources, occupied, destinations bitboard.Bitboard) []move.Move {
	moves := []move.Move{}
	for sources != bitboard.Empty {
		source := bitboard.LeastSignificantOneBit(sources)

		// Remove LSB
		source &= source - 1

		as := attacks(source, occupied) & destinations
		for as != bitboard.Empty {
			a := bitboard.LeastSignificantOneBit(as)
			as &= as - 1
			var m move.Move
			m.SetSourceSquare(source)
			m.SetDestinationSquare(a)
			moves = append(moves, m)
		}
	}
	return moves
}
