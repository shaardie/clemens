package nnue

import (
	"github.com/shaardie/clemens/pkg/types"
)

func CalculateIdx(kingSquare uint8, pieceSquare uint8, piece types.Piece) int {
	pieceIdx := int(piece.Type()*2) + int(piece.Color())
	return int(pieceSquare) + (pieceIdx+int(kingSquare)*10)*64
}
