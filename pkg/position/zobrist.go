package position

import (
	"math/rand"

	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	rnd rand.Source64 = rand.New(rand.NewSource(281954))
)

type zobrist struct {
	piecesOnSquares   [types.SQUARE_NUMBER][types.COLOR_NUMBER][types.PIECE_TYPE_NUMBER]uint64
	sideToMoveIsBlack uint64
	castling          [CASTLING_NUMBER]uint64
	enPassant         [types.FILE_NUMBER]uint64
}

var z = zobrist{}

func init() {
	for i, pieces := range z.piecesOnSquares {
		for ii, pieceTypes := range pieces {
			for iii := range pieceTypes {
				z.piecesOnSquares[i][ii][iii] = rnd.Uint64()
			}
		}
	}

	z.sideToMoveIsBlack = rnd.Uint64()
	for i := range z.castling {
		z.castling[i] = rnd.Uint64()
	}

	for i := range z.enPassant {
		z.enPassant[i] = rnd.Uint64()
	}
}

func (pos *Position) ZobristHash() uint64 {
	var hash uint64
	for square, piece := range pos.piecesBoard {
		if piece != types.NO_PIECE {
			hash ^= z.piecesOnSquares[square][piece.Color()][piece.Type()]
		}
	}
	if pos.sideToMove == types.BLACK {
		hash ^= z.sideToMoveIsBlack
	}

	for i, castling := range []Castling{WHITE_CASTLING_KING, WHITE_CASTLING_QUEEN, BLACK_CASTLING_KING, BLACK_CASTLING_QUEEN} {
		if pos.castling|castling != 0 {
			hash ^= z.castling[i]
		}
	}

	if pos.enPassant != types.SQUARE_NONE {
		hash ^= z.enPassant[bitboard.FileOfSquare(pos.enPassant)]
	}
	return hash
}
