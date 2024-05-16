package evaluation

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	// These tables are from https://github.com/nescitus/cpw-engine/blob/master/eval_init.cpp
	pieceTables = [types.PIECE_TYPE_NUMBER][game_number][types.SQUARE_NUMBER]int16{
		// Pawn
		{
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				50, 50, 50, 50, 50, 50, 50, 50,
				10, 10, 20, 30, 30, 20, 10, 10,
				5, 5, 10, 25, 25, 10, 5, 5,
				0, 0, 0, 20, 20, 0, 0, 0,
				5, -5, -10, 0, 0, -10, -5, 5,
				5, 10, 10, -20, -20, 10, 10, 5,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				50, 50, 50, 50, 50, 50, 50, 50,
				10, 10, 20, 30, 30, 20, 10, 10,
				5, 5, 10, 25, 25, 10, 5, 5,
				0, 0, 0, 20, 20, 0, 0, 0,
				5, -5, -10, 0, 0, -10, -5, 5,
				5, 10, 10, -20, -20, 10, 10, 5,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
		},
		// Knight
		{
			{
				-50, -40, -30, -30, -30, -30, -40, -50,
				-40, -20, 0, 0, 0, 0, -20, -40,
				-30, 0, 10, 15, 15, 10, 0, -30,
				-30, 5, 15, 20, 20, 15, 5, -30,
				-30, 0, 15, 20, 20, 15, 0, -30,
				-30, 5, 10, 15, 15, 10, 5, -30,
				-40, -20, 0, 5, 5, 0, -20, -40,
				-50, -40, -30, -30, -30, -30, -40, -50,
			},
			{
				-50, -40, -30, -30, -30, -30, -40, -50,
				-40, -20, 0, 0, 0, 0, -20, -40,
				-30, 0, 10, 15, 15, 10, 0, -30,
				-30, 5, 15, 20, 20, 15, 5, -30,
				-30, 0, 15, 20, 20, 15, 0, -30,
				-30, 5, 10, 15, 15, 10, 5, -30,
				-40, -20, 0, 5, 5, 0, -20, -40,
				-50, -40, -30, -30, -30, -30, -40, -50,
			},
		},
		// Bishop
		{
			{
				-20, -10, -10, -10, -10, -10, -10, -20,
				-10, 0, 0, 0, 0, 0, 0, -10,
				-10, 0, 5, 10, 10, 5, 0, -10,
				-10, 5, 5, 10, 10, 5, 5, -10,
				-10, 0, 10, 10, 10, 10, 0, -10,
				-10, 10, 10, 10, 10, 10, 10, -10,
				-10, 5, 0, 0, 0, 0, 5, -10,
				-20, -10, -10, -10, -10, -10, -10, -20,
			},
			{
				-20, -10, -10, -10, -10, -10, -10, -20,
				-10, 0, 0, 0, 0, 0, 0, -10,
				-10, 0, 5, 10, 10, 5, 0, -10,
				-10, 5, 5, 10, 10, 5, 5, -10,
				-10, 0, 10, 10, 10, 10, 0, -10,
				-10, 10, 10, 10, 10, 10, 10, -10,
				-10, 5, 0, 0, 0, 0, 5, -10,
				-20, -10, -10, -10, -10, -10, -10, -20,
			},
		},
		// Rook
		{
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				5, 10, 10, 10, 10, 10, 10, 5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				0, 0, 0, 5, 5, 0, 0, 0,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				5, 10, 10, 10, 10, 10, 10, 5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				0, 0, 0, 5, 5, 0, 0, 0,
			},
		},
		// Queen
		{
			{
				-20, -10, -10, -5, -5, -10, -10, -20,
				-10, 0, 0, 0, 0, 0, 0, -10,
				-10, 0, 5, 5, 5, 5, 0, -10,
				-5, 0, 5, 5, 5, 5, 0, -5,
				0, 0, 5, 5, 5, 5, 0, -5,
				-10, 5, 5, 5, 5, 5, 0, -10,
				-10, 0, 5, 0, 0, 0, 0, -10,
				-20, -10, -10, -5, -5, -10, -10, -20,
			},
			{
				-20, -10, -10, -5, -5, -10, -10, -20,
				-10, 0, 0, 0, 0, 0, 0, -10,
				-10, 0, 5, 5, 5, 5, 0, -10,
				-5, 0, 5, 5, 5, 5, 0, -5,
				0, 0, 5, 5, 5, 5, 0, -5,
				-10, 5, 5, 5, 5, 5, 0, -10,
				-10, 0, 5, 0, 0, 0, 0, -10,
				-20, -10, -10, -5, -5, -10, -10, -20,
			},
		},
		// King
		{
			{
				-30, -40, -40, -50, -50, -40, -40, -30,
				-30, -40, -40, -50, -50, -40, -40, -30,
				-30, -40, -40, -50, -50, -40, -40, -30,
				-30, -40, -40, -50, -50, -40, -40, -30,
				-20, -30, -30, -40, -40, -30, -30, -20,
				-10, -20, -20, -20, -20, -20, -20, -10,
				20, 20, 0, 0, 0, 0, 20, 20,
				20, 30, 10, 0, 0, 10, 30, 20,
			},
			{
				-50, -40, -30, -20, -20, -30, -40, -50,
				-30, -20, -10, 0, 0, -10, -20, -30,
				-30, -10, 20, 30, 30, 20, -10, -30,
				-30, -10, 30, 40, 40, 30, -10, -30,
				-30, -10, 30, 40, 40, 30, -10, -30,
				-30, -10, 20, 30, 30, 20, -10, -30,
				-30, -30, 0, 0, 0, 0, -30, -30,
				-50, -30, -30, -30, -30, -30, -30, -50,
			},
		},
	}
	midgamePieceSquareTables [types.COLOR_NUMBER][types.PIECE_TYPE_NUMBER][types.SQUARE_NUMBER]int16
	endgamePieceSquareTables [types.COLOR_NUMBER][types.PIECE_TYPE_NUMBER][types.SQUARE_NUMBER]int16
)

func init() {
	// Set Piece Square Tables
	for square := types.SQUARE_A1; square < types.SQUARE_NUMBER; square++ {
		for _, color := range []types.Color{types.WHITE, types.BLACK} {
			colorAwareSquare := square
			if types.WHITE == color {
				// Flipping square, see https://www.chessprogramming.org/Color_Flipping#Flipping_an_8x8_Board
				colorAwareSquare = square ^ 56
			}
			for pieceType := types.PAWN; pieceType < types.PIECE_TYPE_NUMBER; pieceType++ {
				midgamePieceSquareTables[color][pieceType][square] = pieceTables[pieceType][midgame][colorAwareSquare]
				endgamePieceSquareTables[color][pieceType][square] = pieceTables[pieceType][endgame][colorAwareSquare]
			}
		}
	}
}

// evalPieceSquareTables evaluates the position of each piece based on its current square.
func (e *eval) evalPieceSquareTables(pos *position.Position) {
	bb := bitboard.Empty

	// piece tables
	for pieceType := types.PAWN; pieceType < types.PIECE_TYPE_NUMBER; pieceType++ {
		bb = pos.PiecesBitboard[types.WHITE][pieceType]
		for bb != 0 {
			square := bitboard.SquareIndexSerializationNextSquare(&bb)
			e.phaseScores[midgame] += midgamePieceSquareTables[types.WHITE][pieceType][square]
			e.phaseScores[endgame] += endgamePieceSquareTables[types.WHITE][pieceType][square]
		}
		bb = pos.PiecesBitboard[types.BLACK][pieceType]
		for bb != 0 {
			square := bitboard.SquareIndexSerializationNextSquare(&bb)
			e.phaseScores[midgame] -= midgamePieceSquareTables[types.BLACK][pieceType][square]
			e.phaseScores[endgame] -= endgamePieceSquareTables[types.BLACK][pieceType][square]
		}
	}
}
