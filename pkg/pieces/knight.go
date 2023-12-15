package pieces

import "github.com/shaardie/clemens/pkg/bitboard"

type Knights bitboard.Bitboard

func (k Knights) Attacks() bitboard.Bitboard {
	// Attacks 1 west or east and 2 north or south
	east := bitboard.EastOne(bitboard.Bitboard(k))
	west := bitboard.WestOne(bitboard.Bitboard(k))
	attacks := (west|east)<<16 | (west|east)>>16

	// Attacks 2 west or east and 1 north or south
	east = bitboard.EastOne(east)
	west = bitboard.WestOne(west)
	attacks |= bitboard.NorthOne(west|east) | bitboard.SouthOne(west|east)

	return attacks
}
