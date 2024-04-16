package evaluation

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	// These tables are from https://github.com/nescitus/cpw-engine/blob/master/eval_init.cpp
	pieceTables = [types.PIECE_TYPE_NUMBER][game_number][types.SQUARE_NUMBER]int16{
		/******************************************************************************
		 *                           PAWN PCSQ                                         *
		 *                                                                             *
		 *  Unlike TSCP, CPW generally doesn't want to advance its pawns. Its piece/   *
		 *  square table for pawns takes into account the following factors:           *
		 *                                                                             *
		 *  - file-dependent component, encouraging program to capture                 *
		 *    towards the center                                                       *
		 *  - small bonus for staying on the 2nd rank                                  *
		 *  - small bonus for standing on a3/h3                                        *
		 *  - penalty for d/e pawns on their initial squares                           *
		 *  - bonus for occupying the center                                           *
		 ******************************************************************************/
		{
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				-6, -4, 1, 1, 1, 1, -4, -6,
				-6, -4, 1, 2, 2, 1, -4, -6,
				-6, -4, 2, 8, 8, 2, -4, -6,
				-6, -4, 5, 10, 10, 5, -4, -6,
				-4, -4, 1, 5, 5, 1, -4, -4,
				-6, -4, 1, -24, -24, 1, -4, -6,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				-6, -4, 1, 1, 1, 1, -4, -6,
				-6, -4, 1, 2, 2, 1, -4, -6,
				-6, -4, 2, 8, 8, 2, -4, -6,
				-6, -4, 5, 10, 10, 5, -4, -6,
				-4, -4, 1, 5, 5, 1, -4, -4,
				-6, -4, 1, -24, -24, 1, -4, -6,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
		},
		/******************************************************************************
		 *    KNIGHT PCSQ                                                              *
		 *                                                                             *
		 *   - centralization bonus                                                    *
		 *   - rim and back rank penalty, including penalty for not being developed    *
		 ******************************************************************************/
		{
			{
				-8, -8, -8, -8, -8, -8, -8, -8,
				-8, 0, 0, 0, 0, 0, 0, -8,
				-8, 0, 4, 6, 6, 4, 0, -8,
				-8, 0, 6, 8, 8, 6, 0, -8,
				-8, 0, 6, 8, 8, 6, 0, -8,
				-8, 0, 4, 6, 6, 4, 0, -8,
				-8, 0, 1, 2, 2, 1, 0, -8,
				-16, -12, -8, -8, -8, -8, -12, -16,
			},
			{
				-8, -8, -8, -8, -8, -8, -8, -8,
				-8, 0, 0, 0, 0, 0, 0, -8,
				-8, 0, 4, 6, 6, 4, 0, -8,
				-8, 0, 6, 8, 8, 6, 0, -8,
				-8, 0, 6, 8, 8, 6, 0, -8,
				-8, 0, 4, 6, 6, 4, 0, -8,
				-8, 0, 1, 2, 2, 1, 0, -8,
				-16, -12, -8, -8, -8, -8, -12, -16,
			},
		},
		/******************************************************************************
		 *                BISHOP PCSQ                                                  *
		 *                                                                             *
		 *   - centralization bonus, smaller than for knight                           *
		 *   - penalty for not being developed                                         *
		 *   - good squares on the own half of the board                               *
		 ******************************************************************************/
		{
			{
				-4, -4, -4, -4, -4, -4, -4, -4,
				-4, 0, 0, 0, 0, 0, 0, -4,
				-4, 0, 2, 4, 4, 2, 0, -4,
				-4, 0, 4, 6, 6, 4, 0, -4,
				-4, 0, 4, 6, 6, 4, 0, -4,
				-4, 1, 2, 4, 4, 2, 1, -4,
				-4, 2, 1, 1, 1, 1, 2, -4,
				-4, -4, -12, -4, -4, -12, -4, -4,
			},
			{
				-4, -4, -4, -4, -4, -4, -4, -4,
				-4, 0, 0, 0, 0, 0, 0, -4,
				-4, 0, 2, 4, 4, 2, 0, -4,
				-4, 0, 4, 6, 6, 4, 0, -4,
				-4, 0, 4, 6, 6, 4, 0, -4,
				-4, 1, 2, 4, 4, 2, 1, -4,
				-4, 2, 1, 1, 1, 1, 2, -4,
				-4, -4, -12, -4, -4, -12, -4, -4,
			},
		},
		/******************************************************************************
		*                        ROOK PCSQ                                            *
		*                                                                             *
		*    - bonus for 7th and 8th ranks                                            *
		*    - penalty for a/h columns                                                *
		*    - small centralization bonus                                             *
		******************************************************************************/
		{
			{
				5, 5, 5, 5, 5, 5, 5, 5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				0, 0, 0, 2, 2, 0, 0, 0,
			},
			{
				5, 5, 5, 5, 5, 5, 5, 5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				0, 0, 0, 2, 2, 0, 0, 0,
			},
		},
		/******************************************************************************
		*                     QUEEN PCSQ                                              *
		*                                                                             *
		* - small bonus for centralization in the endgame                             *
		* - penalty for staying on the 1st rank, between rooks in the midgame         *
		******************************************************************************/
		{
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 1, 1, 1, 1, 0, 0,
				0, 0, 1, 2, 2, 1, 0, 0,
				0, 0, 2, 3, 3, 2, 0, 0,
				0, 0, 2, 3, 3, 2, 0, 0,
				0, 0, 1, 2, 2, 1, 0, 0,
				0, 0, 1, 1, 1, 1, 0, 0,
				-5, -5, -5, -5, -5, -5, -5, -5,
			},
			{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 1, 1, 1, 1, 0, 0,
				0, 0, 1, 2, 2, 1, 0, 0,
				0, 0, 2, 3, 3, 2, 0, 0,
				0, 0, 2, 3, 3, 2, 0, 0,
				0, 0, 1, 2, 2, 1, 0, 0,
				0, 0, 1, 1, 1, 1, 0, 0,
				-5, -5, -5, -5, -5, -5, -5, -5,
			},
		},
		/******************************************************************************
			*                     King PCSQ                                               *
		 	*                                                                             *
		 	******************************************************************************/
		{
			{
				-40, -30, -50, -70, -70, -50, -30, -40,
				-30, -20, -40, -60, -60, -40, -20, -30,
				-20, -10, -30, -50, -50, -30, -10, -20,
				-10, 0, -20, -40, -40, -20, 0, -10,
				0, 10, -10, -30, -30, -10, 10, 0,
				10, 20, 0, -20, -20, 0, 20, 10,
				30, 40, 20, 0, 0, 20, 40, 30,
				40, 50, 30, 10, 10, 30, 50, 40,
			},
			{
				-72, -48, -36, -24, -24, -36, -48, -72,
				-48, -24, -12, 0, 0, -12, -24, -48,
				-36, -12, 0, 12, 12, 0, -12, -36,
				-24, 0, 12, 24, 24, 12, 0, -24,
				-24, 0, 12, 24, 24, 12, 0, -24,
				-36, -12, 0, 12, 12, 0, -12, -36,
				-48, -24, -12, 0, 0, -12, -24, -48,
				-72, -48, -36, -24, -24, -36, -48, -72,
			},
		},
	}
	midgamePieceSquareTables [types.COLOR_NUMBER][types.PIECE_TYPE_NUMBER][types.SQUARE_NUMBER]int16
	endgamePieceSquareTables [types.COLOR_NUMBER][types.PIECE_TYPE_NUMBER][types.SQUARE_NUMBER]int16
)

func init() {
	initPieceSquareTables()
}

func initPieceSquareTables() {
	// Set Piece Square Tables
	for square := types.SQUARE_A1; square < types.SQUARE_NUMBER; square++ {
		for _, color := range []types.Color{types.WHITE, types.BLACK} {
			colorAwareSquare := square
			if types.BLACK == color {
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
