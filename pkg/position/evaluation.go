package position

import (
	"github.com/shaardie/clemens/pkg/types"
)

const (
	kingScalar     = 2000
	queenScalar    = 90
	rookScalar     = 50
	bishopScalar   = 30
	knightScalar   = 30
	pawnScalar     = 10
	mobilityScalar = 1
)

func (pos *Position) Evaluation() int {
	score :=
		kingScalar*(pos.piecesBitboard[types.WHITE][types.KING].PopulationCount()-pos.piecesBitboard[types.BLACK][types.KING].PopulationCount()) +
			queenScalar*(pos.piecesBitboard[types.WHITE][types.QUEEN].PopulationCount()-pos.piecesBitboard[types.BLACK][types.QUEEN].PopulationCount()) +
			rookScalar*(pos.piecesBitboard[types.WHITE][types.ROOK].PopulationCount()-pos.piecesBitboard[types.BLACK][types.ROOK].PopulationCount()) +
			bishopScalar*(pos.piecesBitboard[types.WHITE][types.BISHOP].PopulationCount()-pos.piecesBitboard[types.BLACK][types.BISHOP].PopulationCount()) +
			knightScalar*(pos.piecesBitboard[types.WHITE][types.KNIGHT].PopulationCount()-pos.piecesBitboard[types.BLACK][types.KNIGHT].PopulationCount()) +
			pawnScalar*(pos.piecesBitboard[types.WHITE][types.PAWN].PopulationCount()-pos.piecesBitboard[types.BLACK][types.PAWN].PopulationCount()) +
			mobilityScalar*len(pos.GeneratePseudoLegalMoves())

	if pos.sideToMove == types.BLACK {
		score *= -1
	}

	return score
}
