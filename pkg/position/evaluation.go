package position

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	midgame int = iota
	endgame
	game_number
)

var (
	// These tables are from https://github.com/nescitus/cpw-engine/blob/master/eval_init.cpp
	pieceTables = [types.PIECE_TYPE_NUMBER][game_number][types.SQUARE_NUMBER]int{
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
	midgamePieceSquareTables [types.COLOR_NUMBER][types.PIECE_TYPE_NUMBER][types.SQUARE_NUMBER]int
	endgamePieceSquareTables [types.COLOR_NUMBER][types.PIECE_TYPE_NUMBER][types.SQUARE_NUMBER]int

	pieceValue = [types.PIECE_TYPE_NUMBER]int{100, 300, 300, 500, 800, 2000}

	/* adjustements of piece value based on the number of own pawns */
	knight_pawn_adjustment = [9]int{-20, -16, -12, -8, -4, 0, 4, 8, 12}
	rook_pawn_adjustment   = [9]int{15, 12, 9, 6, 3, 0, -3, -6, -9}
)

const (
	shield2Value = 10
	shield3Value = 5

	bishopPair = 30
	knightPair = -8
	rookPair   = -16

	rookSeventhMidgame = 20
	rookSeventhEndgame = 30
	rookOpenFile       = 10
	rookHalfOpenFile   = 5
)

func init() {
	// Set Piece Square Tables
	for square := types.SQUARE_A1; square < types.SQUARE_NUMBER; square++ {
		for _, color := range []types.Color{types.WHITE, types.BLACK} {
			var colorAwareSquare int
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

func (pos *Position) Evaluation() int {
	bb := bitboard.Empty
	var scores [game_number]int

	// Calculate the game phase based on the number of specific PieceTypes,
	// maxed by 24 to a a better linear interpolation later.
	var gamePhase int
	for color := types.WHITE; color < types.COLOR_NUMBER; color++ {
		gamePhase += pos.piecesBitboard[color][types.BISHOP].PopulationCount()
		gamePhase += pos.piecesBitboard[color][types.KNIGHT].PopulationCount()
		gamePhase += 2 * pos.piecesBitboard[color][types.ROOK].PopulationCount()
		gamePhase += 4 * pos.piecesBitboard[color][types.QUEEN].PopulationCount()
	}
	if gamePhase > 24 {
		gamePhase = 24
	}

	// piece tables
	for pieceType := types.PAWN; pieceType < types.PIECE_TYPE_NUMBER; pieceType++ {
		bb = pos.piecesBitboard[types.WHITE][pieceType]
		for bb != 0 {
			square := bitboard.SquareIndexSerializationNextSquare(&bb)
			scores[midgame] += midgamePieceSquareTables[types.WHITE][pieceType][square]
			scores[endgame] += endgamePieceSquareTables[types.WHITE][pieceType][square]
		}
		bb = pos.piecesBitboard[types.BLACK][pieceType]
		for bb != 0 {
			square := bitboard.SquareIndexSerializationNextSquare(&bb)
			scores[midgame] -= midgamePieceSquareTables[types.BLACK][pieceType][square]
			scores[endgame] -= endgamePieceSquareTables[types.BLACK][pieceType][square]
		}
	}

	// Kind Shield Evalutation
	// White
	kingFile := types.FileOfSquare(bitboard.LeastSignificantOneBit(pos.piecesBitboard[types.WHITE][types.KING]))
	// King Side
	if kingFile > types.FILE_E {

		if pos.GetPiece(types.SQUARE_F2) == types.WHITE_PAWN {
			scores[midgame] += shield2Value
		} else if pos.GetPiece(types.SQUARE_F3) == types.WHITE_PAWN {
			scores[midgame] += shield3Value
		}

		if pos.GetPiece(types.SQUARE_G2) == types.WHITE_PAWN {
			scores[midgame] += shield2Value
		} else if pos.GetPiece(types.SQUARE_G3) == types.WHITE_PAWN {
			scores[midgame] += shield3Value
		}

		if pos.GetPiece(types.SQUARE_H2) == types.WHITE_PAWN {
			scores[midgame] += shield2Value
		} else if pos.GetPiece(types.SQUARE_H3) == types.WHITE_PAWN {
			scores[midgame] += shield3Value
		}
	} else
	// Queen Side
	if kingFile < types.FILE_D {
		if pos.GetPiece(types.SQUARE_A2) == types.WHITE_PAWN {
			scores[midgame] += shield2Value
		} else if pos.GetPiece(types.SQUARE_A3) == types.WHITE_PAWN {
			scores[midgame] += shield3Value
		}

		if pos.GetPiece(types.SQUARE_B2) == types.WHITE_PAWN {
			scores[midgame] += shield2Value
		} else if pos.GetPiece(types.SQUARE_B3) == types.WHITE_PAWN {
			scores[midgame] += shield3Value
		}

		if pos.GetPiece(types.SQUARE_C2) == types.WHITE_PAWN {
			scores[midgame] += shield2Value
		} else if pos.GetPiece(types.SQUARE_C3) == types.WHITE_PAWN {
			scores[midgame] += shield3Value
		}
	}

	// Black
	kingFile = types.FileOfSquare(bitboard.LeastSignificantOneBit(pos.piecesBitboard[types.BLACK][types.KING]))
	// King Side
	if kingFile > types.FILE_E {

		if pos.GetPiece(types.SQUARE_F7) == types.BLACK_PAWN {
			scores[midgame] -= shield2Value
		} else if pos.GetPiece(types.SQUARE_F6) == types.BLACK_PAWN {
			scores[midgame] -= shield3Value
		}

		if pos.GetPiece(types.SQUARE_G7) == types.BLACK_PAWN {
			scores[midgame] -= shield2Value
		} else if pos.GetPiece(types.SQUARE_G6) == types.BLACK_PAWN {
			scores[midgame] -= shield3Value
		}

		if pos.GetPiece(types.SQUARE_H7) == types.BLACK_PAWN {
			scores[midgame] -= shield2Value
		} else if pos.GetPiece(types.SQUARE_H6) == types.BLACK_PAWN {
			scores[midgame] -= shield3Value
		}
	} else
	// Queen Side
	if kingFile < types.FILE_D {
		if pos.GetPiece(types.SQUARE_A7) == types.WHITE_PAWN {
			scores[midgame] -= shield2Value
		} else if pos.GetPiece(types.SQUARE_A6) == types.WHITE_PAWN {
			scores[midgame] -= shield3Value
		}

		if pos.GetPiece(types.SQUARE_B7) == types.WHITE_PAWN {
			scores[midgame] -= shield2Value
		} else if pos.GetPiece(types.SQUARE_B6) == types.WHITE_PAWN {
			scores[midgame] -= shield3Value
		}

		if pos.GetPiece(types.SQUARE_C7) == types.WHITE_PAWN {
			scores[midgame] -= shield2Value
		} else if pos.GetPiece(types.SQUARE_C6) == types.WHITE_PAWN {
			scores[midgame] -= shield3Value
		}
	}

	// Rook Evaluation
	// Rook on the seventh with either the king on the eigth or enemy pawns on the seven
	// https://www.chessprogramming.org/Rook_on_Seventh
	if pos.piecesBitboard[types.WHITE][types.ROOK]&bitboard.RankMask7 > 0 && (pos.piecesBitboard[types.BLACK][types.KING]&bitboard.RankMask8 > 0 ||
		pos.piecesBitboard[types.BLACK][types.PAWN]&bitboard.RankMask7 > 0) {
		scores[midgame] += rookSeventhMidgame
		scores[endgame] += rookSeventhEndgame
	}
	// vice versa
	if pos.piecesBitboard[types.BLACK][types.ROOK]&bitboard.RankMask2 > 0 && (pos.piecesBitboard[types.WHITE][types.KING]&bitboard.RankMask1 > 0 ||
		pos.piecesBitboard[types.WHITE][types.PAWN]&bitboard.RankMask2 > 0) {
		scores[midgame] -= rookSeventhMidgame
		scores[endgame] -= rookSeventhEndgame
	}

	// Rooks are on an open or half open file
	// https://www.chessprogramming.org/Rook_on_Open_File
	rooks := pos.piecesBitboard[types.WHITE][types.ROOK]
	for rooks != 0 {
		fileMask := bitboard.FileMaskOfSquare(bitboard.SquareIndexSerializationNextSquare(&rooks))
		if pos.piecesBitboard[types.BLACK][types.PAWN]&fileMask > 0 {
			if pos.piecesBitboard[types.WHITE][types.PAWN]&fileMask > 0 {
				scores[midgame] += rookOpenFile
				continue
			}
			scores[endgame] += rookHalfOpenFile
		}
	}
	rooks = pos.piecesBitboard[types.BLACK][types.ROOK]
	for rooks != 0 {
		fileMask := bitboard.FileMaskOfSquare(bitboard.SquareIndexSerializationNextSquare(&rooks))
		if pos.piecesBitboard[types.WHITE][types.PAWN]&fileMask > 0 {
			if pos.piecesBitboard[types.BLACK][types.PAWN]&fileMask > 0 {
				scores[midgame] -= rookOpenFile
				continue
			}
			scores[endgame] -= rookHalfOpenFile
		}
	}

	// Merge midgame and endgame value
	score := (scores[midgame]*gamePhase + scores[endgame]*(24-gamePhase)) / 24

	// Basic Material Score
	for pieceType := types.PAWN; pieceType < types.PIECE_TYPE_NUMBER; pieceType++ {
		score += pieceValue[pieceType] * (pos.piecesBitboard[types.WHITE][pieceType].PopulationCount() - pos.piecesBitboard[types.BLACK][pieceType].PopulationCount())
	}

	// Adjust piece type values based on different factor
	// Pairs give bonus or malus, see for example https://www.chessprogramming.org/Bishop_Pair
	if pos.piecesBitboard[types.WHITE][types.BISHOP].PopulationCount() > 1 {
		score += bishopPair
	}
	if pos.piecesBitboard[types.BLACK][types.BISHOP].PopulationCount() > 1 {
		score -= bishopPair
	}
	if pos.piecesBitboard[types.WHITE][types.KNIGHT].PopulationCount() > 1 {
		score += knightPair
	}
	if pos.piecesBitboard[types.BLACK][types.KNIGHT].PopulationCount() > 1 {
		score -= knightPair
	}
	if pos.piecesBitboard[types.WHITE][types.ROOK].PopulationCount() > 1 {
		score += rookPair
	}
	if pos.piecesBitboard[types.BLACK][types.ROOK].PopulationCount() > 1 {
		score -= rookPair
	}

	// Adjustments based on the number of pawns
	numberOfWhitePawns := pos.piecesBitboard[types.WHITE][types.PAWN].PopulationCount()
	numberOfBlackPawns := pos.piecesBitboard[types.BLACK][types.PAWN].PopulationCount()
	score += knight_pawn_adjustment[numberOfWhitePawns] * pos.piecesBitboard[types.WHITE][types.KNIGHT].PopulationCount()
	score -= knight_pawn_adjustment[numberOfBlackPawns] * pos.piecesBitboard[types.BLACK][types.KNIGHT].PopulationCount()
	score += rook_pawn_adjustment[numberOfWhitePawns] * pos.piecesBitboard[types.WHITE][types.ROOK].PopulationCount()
	score -= rook_pawn_adjustment[numberOfBlackPawns] * pos.piecesBitboard[types.BLACK][types.ROOK].PopulationCount()

	// Make the result side aware
	if pos.SideToMove == types.BLACK {
		score *= -1
	}

	return score
}
