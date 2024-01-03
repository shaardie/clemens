package utils

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/move"
)

func GenerateMoves(squares, occupied, destinations bitboard.Bitboard, attacks func(square int, occupied bitboard.Bitboard) bitboard.Bitboard) []move.Move {
	moves := []move.Move{}
	for squares != bitboard.Empty {
		source := bitboard.LeastSignificantOneBit(squares)

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
