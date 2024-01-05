package position

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

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

type CastlingSide int

const (
	CASTLING_QUEEN CastlingSide = iota
	CASTLING_KING
)

func (c Castling) Side() CastlingSide {
	switch c {
	case WHITE_CASTLING_KING:
		return CASTLING_KING
	case BLACK_CASTLING_KING:
		return CASTLING_KING
	case WHITE_CASTLING_QUEEN:
		return CASTLING_QUEEN
	case BLACK_CASTLING_QUEEN:
		return CASTLING_QUEEN
	}
	panic("unknown castling")
}
func (c Castling) Color() types.Color {
	switch c {
	case WHITE_CASTLING_KING:
		return types.WHITE
	case WHITE_CASTLING_QUEEN:
		return types.WHITE
	case BLACK_CASTLING_KING:
		return types.BLACK
	case BLACK_CASTLING_QUEEN:
		return types.BLACK
	}
	panic("unknown castling")
}

// CanCastle returns true, if castling is possible possible in theory,
// because rook and kind did not move yet.
func (pos *Position) CanCastle(c Castling) bool {
	return c&pos.castling != NO_CASTLING
}

// CanCastleNow returns true, if castling is now.
// It checks, if the path between the pieces is free and not attacked
// and there is no check.
func (pos *Position) CanCastleNow(c Castling) bool {
	// Can Castle is possible in theory
	if !pos.CanCastle(c) {
		return false
	}

	// Color matches
	if c.Color() != pos.sideToMove {
		return false
	}

	// No check
	if pos.IsCheck() {
		return false
	}

	var rank int
	switch c.Color() {
	case types.WHITE:
		rank = types.RANK_1
	case types.BLACK:
		rank = types.RANK_8
	default:
		panic("unkown color")
	}

	var files []int
	switch c.Side() {
	case CASTLING_QUEEN:
		files = []int{
			types.FILE_B,
			types.FILE_C,
			types.FILE_D,
		}
	case CASTLING_KING:
		files = []int{
			types.FILE_F,
			types.FILE_G,
			types.FILE_H,
		}
	default:
		panic("unkown color")
	}

	for _, file := range files {
		square := bitboard.SquareFromRankAndFile(rank, file)
		// Check if squares between are empty
		if file != types.FILE_H && !pos.Empty(square) {
			return false
		}
		// Check if squares between are attacked
		if pos.SquareAttackedBy(square)&pos.AllPiecesByColor(types.SwitchColor(pos.sideToMove)) != bitboard.Empty {
			return false
		}
	}

	return true
}
