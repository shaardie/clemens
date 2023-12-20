package board

import (
	"github.com/shaardie/clemens/pkg/bitboard"
)

type Board struct {
	whitePawns   bitboard.Bitboard
	whiteKnights bitboard.Bitboard
	whiteBishops bitboard.Bitboard
	whiteRools   bitboard.Bitboard
	whiteQueens  bitboard.Bitboard
	whiteKing    bitboard.Bitboard

	blackPawns   bitboard.Bitboard
	blackKnights bitboard.Bitboard
	blackBishops bitboard.Bitboard
	blackRools   bitboard.Bitboard
	blackQueens  bitboard.Bitboard
	blackKing    bitboard.Bitboard
}

func NewBoard() Board {
	return Board{
		whiteKnights: bitboard.Bitboard(
			bitboard.BitBySquares(
				bitboard.SQUARE_B1,
				bitboard.SQUARE_G1,
			),
		),
		blackKnights: bitboard.Bitboard(
			bitboard.BitBySquares(
				bitboard.SQUARE_B8,
				bitboard.SQUARE_G8,
			),
		),
	}
}
