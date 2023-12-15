package pieces

import "github.com/shaardie/clemens/pkg/bitboard"

type King bitboard.Bitboard

func (k King) Attacks() bitboard.Bitboard {
	// Attacks to West and East
	attacks := bitboard.WestOne(bitboard.Bitboard(k)) | bitboard.EastOne(bitboard.Bitboard(k))
	// Use attackes to also calculate attacks in north and south,
	// so it also includes northeast, etc.
	attacks |= bitboard.NorthOne(attacks) | bitboard.SouthOne(attacks)

	return attacks
}
