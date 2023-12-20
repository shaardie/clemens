package board

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	. "github.com/shaardie/clemens/pkg/types"
)

type Board struct {
	whitePawns   bitboard.Bitboard
	whiteKnights bitboard.Bitboard
	whiteBishops bitboard.Bitboard
	whiteRooks   bitboard.Bitboard
	whiteQueens  bitboard.Bitboard
	whiteKing    bitboard.Bitboard

	blackPawns   bitboard.Bitboard
	blackKnights bitboard.Bitboard
	blackBishops bitboard.Bitboard
	blackRooks   bitboard.Bitboard
	blackQueens  bitboard.Bitboard
	blackKing    bitboard.Bitboard
}

func NewBoard() Board {
	return Board{
		whitePawns: bitboard.BitBySquares(
			SQUARE_A2,
			SQUARE_B2,
			SQUARE_C2,
			SQUARE_D2,
			SQUARE_E2,
			SQUARE_F2,
			SQUARE_G2,
			SQUARE_H2,
		),
		whiteKnights: bitboard.BitBySquares(
			SQUARE_B1,
			SQUARE_G1,
		),
		whiteBishops: bitboard.BitBySquares(
			SQUARE_C1,
			SQUARE_F1,
		),
		whiteRooks: bitboard.BitBySquares(
			SQUARE_A1,
			SQUARE_H1,
		),
		whiteQueens: bitboard.BitBySquares(
			SQUARE_D1,
		),
		whiteKing: bitboard.BitBySquares(
			SQUARE_E1,
		),
		blackPawns: bitboard.BitBySquares(
			SQUARE_A7,
			SQUARE_B7,
			SQUARE_C7,
			SQUARE_D7,
			SQUARE_E7,
			SQUARE_F7,
			SQUARE_G7,
			SQUARE_H7,
		),
		blackKnights: bitboard.BitBySquares(
			SQUARE_B8,
			SQUARE_G8,
		),
		blackBishops: bitboard.BitBySquares(
			SQUARE_C8,
			SQUARE_F8,
		),
		blackRooks: bitboard.BitBySquares(
			SQUARE_A8,
			SQUARE_H8,
		),
		blackQueens: bitboard.BitBySquares(
			SQUARE_D8,
		),
		blackKing: bitboard.BitBySquares(
			SQUARE_E8,
		),
	}
}
