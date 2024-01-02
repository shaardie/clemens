package utils

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/move"
)

func SlidingAttacks(square int, directions []func(bitboard.Bitboard) bitboard.Bitboard, occupied bitboard.Bitboard) bitboard.Bitboard {
	attacks := bitboard.Empty

	for _, direction := range directions {
		b := bitboard.BitBySquares(square)
		for {
			nextMove := direction(b)
			if nextMove == bitboard.Empty || b&occupied != bitboard.Empty {
				break
			}
			attacks |= nextMove
			b = nextMove
		}
	}
	return attacks
}

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
