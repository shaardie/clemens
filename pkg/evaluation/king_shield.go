package evaluation

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	shield2Value = 10
	shield3Value = 5
)

func (e *eval) evalKingShield(pos *position.Position) {
	// Kind Shield Evalutation
	// White
	kingFile := types.FileOfSquare(bitboard.LeastSignificantOneBit(pos.PiecesBitboard[types.WHITE][types.KING]))
	// King Side
	if kingFile > types.FILE_E {

		if pos.GetPiece(types.SQUARE_F2) == types.WHITE_PAWN {
			e.phaseScores[midgame] += shield2Value
		} else if pos.GetPiece(types.SQUARE_F3) == types.WHITE_PAWN {
			e.phaseScores[midgame] += shield3Value
		}

		if pos.GetPiece(types.SQUARE_G2) == types.WHITE_PAWN {
			e.phaseScores[midgame] += shield2Value
		} else if pos.GetPiece(types.SQUARE_G3) == types.WHITE_PAWN {
			e.phaseScores[midgame] += shield3Value
		}

		if pos.GetPiece(types.SQUARE_H2) == types.WHITE_PAWN {
			e.phaseScores[midgame] += shield2Value
		} else if pos.GetPiece(types.SQUARE_H3) == types.WHITE_PAWN {
			e.phaseScores[midgame] += shield3Value
		}
	} else
	// Queen Side
	if kingFile < types.FILE_D {
		if pos.GetPiece(types.SQUARE_A2) == types.WHITE_PAWN {
			e.phaseScores[midgame] += shield2Value
		} else if pos.GetPiece(types.SQUARE_A3) == types.WHITE_PAWN {
			e.phaseScores[midgame] += shield3Value
		}

		if pos.GetPiece(types.SQUARE_B2) == types.WHITE_PAWN {
			e.phaseScores[midgame] += shield2Value
		} else if pos.GetPiece(types.SQUARE_B3) == types.WHITE_PAWN {
			e.phaseScores[midgame] += shield3Value
		}

		if pos.GetPiece(types.SQUARE_C2) == types.WHITE_PAWN {
			e.phaseScores[midgame] += shield2Value
		} else if pos.GetPiece(types.SQUARE_C3) == types.WHITE_PAWN {
			e.phaseScores[midgame] += shield3Value
		}
	}

	// Black
	kingFile = types.FileOfSquare(bitboard.LeastSignificantOneBit(pos.PiecesBitboard[types.BLACK][types.KING]))
	// King Side
	if kingFile > types.FILE_E {

		if pos.GetPiece(types.SQUARE_F7) == types.BLACK_PAWN {
			e.phaseScores[midgame] -= shield2Value
		} else if pos.GetPiece(types.SQUARE_F6) == types.BLACK_PAWN {
			e.phaseScores[midgame] -= shield3Value
		}

		if pos.GetPiece(types.SQUARE_G7) == types.BLACK_PAWN {
			e.phaseScores[midgame] -= shield2Value
		} else if pos.GetPiece(types.SQUARE_G6) == types.BLACK_PAWN {
			e.phaseScores[midgame] -= shield3Value
		}

		if pos.GetPiece(types.SQUARE_H7) == types.BLACK_PAWN {
			e.phaseScores[midgame] -= shield2Value
		} else if pos.GetPiece(types.SQUARE_H6) == types.BLACK_PAWN {
			e.phaseScores[midgame] -= shield3Value
		}
	} else
	// Queen Side
	if kingFile < types.FILE_D {
		if pos.GetPiece(types.SQUARE_A7) == types.WHITE_PAWN {
			e.phaseScores[midgame] -= shield2Value
		} else if pos.GetPiece(types.SQUARE_A6) == types.WHITE_PAWN {
			e.phaseScores[midgame] -= shield3Value
		}

		if pos.GetPiece(types.SQUARE_B7) == types.WHITE_PAWN {
			e.phaseScores[midgame] -= shield2Value
		} else if pos.GetPiece(types.SQUARE_B6) == types.WHITE_PAWN {
			e.phaseScores[midgame] -= shield3Value
		}

		if pos.GetPiece(types.SQUARE_C7) == types.WHITE_PAWN {
			e.phaseScores[midgame] -= shield2Value
		} else if pos.GetPiece(types.SQUARE_C6) == types.WHITE_PAWN {
			e.phaseScores[midgame] -= shield3Value
		}
	}

}
