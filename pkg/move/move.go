package move

type MoveType int

const (
	NORMAL MoveType = iota
	PROMOTION
	EN_PASSANT
	CASTLING
)

type Move struct {
	SourceSquare      int
	DestinationSquare int
	MoveType          MoveType
}

func (m *Move) GetSourceSquare() int {
	return m.SourceSquare
}
func (m *Move) GetDestinationSquare() int {
	return m.DestinationSquare
}
func (m *Move) GetMoveType() MoveType {
	return m.MoveType
}

func (m *Move) SetSourceSquare(square int) {
	m.SourceSquare = square
}
func (m *Move) SetDestinationSquare(square int) {
	m.DestinationSquare = square
}
func (m *Move) SetMoveType(mt MoveType) {
	m.MoveType = mt
}
