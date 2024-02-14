package position

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

func (pos *Position) IsInCheck(c types.Color) bool {
	square := bitboard.LeastSignificantOneBit(pos.PiecesBitboard[c][types.KING])
	attacks := pos.SquareAttackedBy(square)
	filtered := attacks & pos.AllPiecesByColor[types.SwitchColor(c)]
	return filtered != 0
}
