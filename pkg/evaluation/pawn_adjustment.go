package evaluation

import (
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	/* adjustements of piece value based on the number of own pawns */
	knight_pawn_adjustment = [9]int16{-20, -16, -12, -8, -4, 0, 4, 8, 12}
	rook_pawn_adjustment   = [9]int16{15, 12, 9, 6, 3, 0, -3, -6, -9}
)

// Adjustments based on the number of pawns
func (e *eval) evalPawnAdjustment(pos *position.Position) {
	numberOfWhitePawns := pos.PiecesBitboard[types.WHITE][types.PAWN].PopulationCount()
	numberOfBlackPawns := pos.PiecesBitboard[types.BLACK][types.PAWN].PopulationCount()
	e.baseScore += knight_pawn_adjustment[numberOfWhitePawns] * int16(pos.PiecesBitboard[types.WHITE][types.KNIGHT].PopulationCount())
	e.baseScore -= knight_pawn_adjustment[numberOfBlackPawns] * int16(pos.PiecesBitboard[types.BLACK][types.KNIGHT].PopulationCount())
	e.baseScore += rook_pawn_adjustment[numberOfWhitePawns] * int16(pos.PiecesBitboard[types.WHITE][types.ROOK].PopulationCount())
	e.baseScore -= rook_pawn_adjustment[numberOfBlackPawns] * int16(pos.PiecesBitboard[types.BLACK][types.ROOK].PopulationCount())
}
