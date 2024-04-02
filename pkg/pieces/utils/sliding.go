package utils

import (
	"github.com/shaardie/clemens/pkg/bitboard"
)

func SlidingAttacks(square uint8, directions []func(bitboard.Bitboard) bitboard.Bitboard, occupied bitboard.Bitboard) bitboard.Bitboard {
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
