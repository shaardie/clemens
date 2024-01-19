package position

import (
	"github.com/shaardie/clemens/pkg/types"
)

const (
	kingScalar   = 200
	queenScalar  = 9
	rookScalar   = 5
	bishopScalar = 3
	knightScalar = 3
	pawnScalar   = 1
)

func (pos *Position) Evaluation() int {
	materialScore :=
		kingScalar*(pos.piecesBitboard[types.WHITE][types.KING].PopulationCount()-pos.piecesBitboard[types.BLACK][types.KING].PopulationCount()) +
			queenScalar*(pos.piecesBitboard[types.WHITE][types.QUEEN].PopulationCount()-pos.piecesBitboard[types.BLACK][types.QUEEN].PopulationCount()) +
			rookScalar*(pos.piecesBitboard[types.WHITE][types.ROOK].PopulationCount()-pos.piecesBitboard[types.BLACK][types.ROOK].PopulationCount()) +
			bishopScalar*(pos.piecesBitboard[types.WHITE][types.BISHOP].PopulationCount()-pos.piecesBitboard[types.BLACK][types.BISHOP].PopulationCount()) +
			knightScalar*(pos.piecesBitboard[types.WHITE][types.KNIGHT].PopulationCount()-pos.piecesBitboard[types.BLACK][types.KNIGHT].PopulationCount()) +
			pawnScalar*(pos.piecesBitboard[types.WHITE][types.PAWN].PopulationCount()-pos.piecesBitboard[types.BLACK][types.PAWN].PopulationCount())

	if pos.sideToMove == types.BLACK {
		materialScore *= -1
	}

	return materialScore
}
