package move

import (
	"fmt"

	"github.com/shaardie/clemens/pkg/types"
)

type MoveType int

const (
	NORMAL MoveType = iota
	PROMOTION
	EN_PASSANT
	CASTLING
)

// Move represents a move from one position to another
// 0-5 is the source square
// 6-11 is the destination square
// 12-13 is the Move Type
// 14-15 is the Promotion Piece Type
// 16-31 are for the score
type Move uint32

const (
	NullMove Move = 0
)

func (m Move) String() string {
	r := fmt.Sprintf(
		"%s%s",
		types.SquareToString(m.GetSourceSquare()),
		types.SquareToString(m.GetTargetSquare()),
	)

	if m.GetMoveType() == PROMOTION {
		r = fmt.Sprintf("%s%v", r, m.GetPromitionPieceType())
	}
	return r
}

func (m *Move) GetSourceSquare() uint8 {
	return uint8(*m & 0b111111)
}

func (m *Move) SetSourceSquare(square uint8) *Move {
	*m |= Move(square)
	return m
}

func (m *Move) GetTargetSquare() uint8 {
	return uint8(*m >> 6 & 0b111111)
}

func (m *Move) SetTargetSquare(square uint8) *Move {
	*m |= Move(square) << 6
	return m
}

func (m *Move) GetMoveType() MoveType {
	return MoveType(*m >> 12 & 0b11)
}

func (m *Move) SetMoveType(mt MoveType) *Move {
	*m |= Move(mt << 12)
	return m
}

func (m *Move) GetPromitionPieceType() types.PieceType {
	return types.PieceType(*m>>14&0b11) + 1
}

func (m *Move) SetPromitionPieceType(pt types.PieceType) *Move {
	*m |= Move(0 << 14)
	*m |= (Move(pt) - 1) << 14
	return m
}

func (m *Move) GetScore() uint16 {
	return uint16(*m >> 16)
}

func (m *Move) SetScore(s uint16) {
	*m |= Move(s) << 16
}
