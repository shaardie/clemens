package bitboard

type Knights Bitboard

func (k Knights) Attacks() Bitboard {
	// Attacks 1 west or east and 2 north or south
	east := EastOne(Bitboard(k))
	west := WestOne(Bitboard(k))
	attacks := (west|east)<<16 | (west|east)>>16

	// Attacks 2 west or east and 1 north or south
	east = EastOne(east)
	west = WestOne(west)
	attacks |= NorthOne(west|east) | SouthOne(west|east)

	return attacks
}
