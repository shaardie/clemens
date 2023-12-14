package bitboard

type Board struct {
	whitePawns   Bitboard
	whiteKnights Bitboard
	whiteBishops Bitboard
	whiteRools   Bitboard
	whiteQueens  Bitboard
	whiteKing    Bitboard

	blackPawns   Bitboard
	blackKnights Bitboard
	blackBishops Bitboard
	blackRools   Bitboard
	blackQueens  Bitboard
	blackKing    Bitboard
}

func (b Board) String() string {
	fields := [64]string{}
	for _, v := range []struct {
		b Bitboard
		c string
	}{
		{
			b.whitePawns,
			"wP",
		},
	} {
		for _, idx := range SquareIndexSerialization(v.b) {
			fields[idx] = v.c
		}
	}
	return ""
}
