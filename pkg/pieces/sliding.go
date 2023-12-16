package pieces

import (
	"github.com/shaardie/clemens/pkg/bitboard"
)

func SlidingAttacks(t pieceType, square int, occupied bitboard.Bitboard) bitboard.Bitboard {
	attacks := bitboard.Empty

	var directions []func(bitboard.Bitboard) bitboard.Bitboard
	switch t {
	case RookType:
		directions = []func(bitboard.Bitboard) bitboard.Bitboard{
			bitboard.NorthOne,
			bitboard.SouthOne,
			bitboard.EastOne,
			bitboard.WestOne,
		}
	case BishopType:
		directions = []func(bitboard.Bitboard) bitboard.Bitboard{
			bitboard.NorthEastOne,
			bitboard.NorthWestOne,
			bitboard.SouthEastOne,
			bitboard.SouthWestOne,
		}
	default:
		panic("Not implemented for this pieceType")
	}

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
