package position

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

func (pos *Position) IsInCheck(c types.Color) bool {
	square := bitboard.SquareIndexSerialization(pos.piecesBitboard[c][types.KING])[0]
	attacks := pos.SquareAttackedBy(square)
	filtered := attacks & pos.AllPiecesByColor(types.SwitchColor(c))
	return filtered != 0
}
