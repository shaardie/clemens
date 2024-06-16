package evaluation

import (
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	PieceValue = [types.PIECE_TYPE_NUMBER]int16{100, 310, 310, 510, 910, 0}
)

func (e *eval) evalBaseMaterial(pos *position.Position) {
	// Basic Material Score
	for pieceType := types.PAWN; pieceType < types.PIECE_TYPE_NUMBER; pieceType++ {
		e.baseScore += PieceValue[pieceType] * int16(pos.PiecesBitboard[types.WHITE][pieceType].PopulationCount()-pos.PiecesBitboard[types.BLACK][pieceType].PopulationCount())
	}
}
