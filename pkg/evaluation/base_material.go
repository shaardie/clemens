package evaluation

import (
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	pieceValue = [types.PIECE_TYPE_NUMBER]int{100, 300, 300, 500, 800, 2000}
)

func (e *eval) evalBaseMaterial(pos *position.Position) {
	// Basic Material Score
	for pieceType := types.PAWN; pieceType < types.PIECE_TYPE_NUMBER; pieceType++ {
		e.baseScore += pieceValue[pieceType] * (pos.PiecesBitboard[types.WHITE][pieceType].PopulationCount() - pos.PiecesBitboard[types.BLACK][pieceType].PopulationCount())
	}
}
