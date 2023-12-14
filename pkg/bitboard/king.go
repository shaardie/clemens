package bitboard

type King Bitboard

func (k King) Attacks() Bitboard {
	// Attacks to West and East
	attacks := WestOne(Bitboard(k)) | EastOne(Bitboard(k))
	// Use attackes to also calculate attacks in north and south,
	// so it also includes northeast, etc.
	attacks |= NorthOne(attacks) | SouthOne(attacks)

	return attacks
}
