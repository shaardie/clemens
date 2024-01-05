package position

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

func (pos *Position) IsCheck() bool {
	square := bitboard.SquareIndexSerialization(pos.piecesBitboard[pos.sideToMove][types.KING])[0]
	attacks := pos.SquareAttackedBy(square)
	filtered := attacks & pos.AllPiecesByColor(types.SwitchColor(pos.sideToMove))
	return filtered != 0
}
