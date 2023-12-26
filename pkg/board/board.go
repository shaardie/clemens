package board

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
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
			types.SQUARE_A2,
			types.SQUARE_B2,
			types.SQUARE_C2,
			types.SQUARE_D2,
			types.SQUARE_E2,
			types.SQUARE_F2,
			types.SQUARE_G2,
			types.SQUARE_H2,
		),
		whiteKnights: bitboard.BitBySquares(
			types.SQUARE_B1,
			types.SQUARE_G1,
		),
		whiteBishops: bitboard.BitBySquares(
			types.SQUARE_C1,
			types.SQUARE_F1,
		),
		whiteRooks: bitboard.BitBySquares(
			types.SQUARE_A1,
			types.SQUARE_H1,
		),
		whiteQueens: bitboard.BitBySquares(
			types.SQUARE_D1,
		),
		whiteKing: bitboard.BitBySquares(
			types.SQUARE_E1,
		),
		blackPawns: bitboard.BitBySquares(
			types.SQUARE_A7,
			types.SQUARE_B7,
			types.SQUARE_C7,
			types.SQUARE_D7,
			types.SQUARE_E7,
			types.SQUARE_F7,
			types.SQUARE_G7,
			types.SQUARE_H7,
		),
		blackKnights: bitboard.BitBySquares(
			types.SQUARE_B8,
			types.SQUARE_G8,
		),
		blackBishops: bitboard.BitBySquares(
			types.SQUARE_C8,
			types.SQUARE_F8,
		),
		blackRooks: bitboard.BitBySquares(
			types.SQUARE_A8,
			types.SQUARE_H8,
		),
		blackQueens: bitboard.BitBySquares(
			types.SQUARE_D8,
		),
		blackKing: bitboard.BitBySquares(
			types.SQUARE_E8,
		),
	}
}
