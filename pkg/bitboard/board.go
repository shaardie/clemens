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
