package evaluation

import (
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	rookPair   int16 = -16
	knightPair int16 = -8
	bishopPair int16 = 30
)

func (e *eval) evalPairs(pos *position.Position) {
	// Pairs bonus/malus, see for example https://www.chessprogramming.org/Bishop_Pair
	if pos.PiecesBitboard[types.WHITE][types.BISHOP].PopulationCount() > 1 {
		e.baseScore += bishopPair
	}
	if pos.PiecesBitboard[types.BLACK][types.BISHOP].PopulationCount() > 1 {
		e.baseScore -= bishopPair
	}
	if pos.PiecesBitboard[types.WHITE][types.KNIGHT].PopulationCount() > 1 {
		e.baseScore += knightPair
	}
	if pos.PiecesBitboard[types.BLACK][types.KNIGHT].PopulationCount() > 1 {
		e.baseScore -= knightPair
	}
	if pos.PiecesBitboard[types.WHITE][types.ROOK].PopulationCount() > 1 {
		e.baseScore += rookPair
	}
	if pos.PiecesBitboard[types.BLACK][types.ROOK].PopulationCount() > 1 {
		e.baseScore -= rookPair
	}
}
