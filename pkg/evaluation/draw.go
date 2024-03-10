package evaluation

import (
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

// Draw Evaluation, https://www.chessprogramming.org/Draw_Evaluation
func (e *eval) evalDraw(pos *position.Position) bool {
	// There are only kings left, it is a draw
	if pos.AllPieces.PopulationCount() == 2 {
		return true
	}

	// If there is any Pawn, Rook or Queen, it is no draw
	if (pos.PiecesBitboard[types.WHITE][types.PAWN] |
		pos.PiecesBitboard[types.BLACK][types.PAWN] |
		pos.PiecesBitboard[types.WHITE][types.ROOK] |
		pos.PiecesBitboard[types.BLACK][types.ROOK] |
		pos.PiecesBitboard[types.WHITE][types.QUEEN] |
		pos.PiecesBitboard[types.BLACK][types.QUEEN]).PopulationCount() > 0 {
		return false
	}
	// At this point there are now Pawns, Rooks or Queens on the board

	numberOfPieces := [types.COLOR_NUMBER]int{
		pos.AllPiecesByColor[types.WHITE].PopulationCount(),
		pos.AllPiecesByColor[types.BLACK].PopulationCount(),
	}

	// If both side have only one minor piece
	if numberOfPieces[types.WHITE] == 2 && numberOfPieces[types.BLACK] == 2 {
		return true
	}

	// If both side have more than one minor piece
	if numberOfPieces[types.WHITE] > 2 && numberOfPieces[types.BLACK] > 2 {
		return false
	}

	// One Side has more than 2 minor pieces
	if numberOfPieces[types.WHITE] > 3 || numberOfPieces[types.BLACK] > 3 {
		return false
	}

	// At this point one side has 2 minor pieces and the other one has only one minor piece.
	// Everything is now a draw, except two Bishops against a Knight.
	for color := types.WHITE; color < types.COLOR_NUMBER; color++ {
		we := color
		them := types.SwitchColor(we)

		// To Bishops agains another minor piece is a draw except
		if pos.PiecesBitboard[we][types.BISHOP].PopulationCount() == 2 {
			return pos.PiecesBitboard[them][types.BISHOP].PopulationCount() == 1
		}
	}

	return true

}
