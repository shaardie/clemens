package position

type Castling int

const (
	NO_CASTLING          Castling = 0
	WHITE_CASTLING_KING  Castling = 1
	WHITE_CASTLING_QUEEN Castling = 2
	BLACK_CASTLING_KING  Castling = 4
	BLACK_CASTLING_QUEEN Castling = 8
	ANY_CASTLING                  = WHITE_CASTLING_KING | WHITE_CASTLING_QUEEN | BLACK_CASTLING_KING | BLACK_CASTLING_QUEEN
	CASTLING_NUMBER      int      = 4
)

// CanCastle returns true, if castling is possible
func (pos *Position) CanCastle(c Castling) bool {
	return c&pos.castling != NO_CASTLING
}
