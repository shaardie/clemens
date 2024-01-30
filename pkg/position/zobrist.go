package position

import (
	"math/bits"
	"math/rand"

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

func (pos *Position) initZobristHash() {
	pos.ZobristHash = 0

	for square, piece := range pos.piecesBoard {
		if piece != types.NO_PIECE {
			pos.zobristUpdatePiece(square, piece.Color(), piece.Type())
		}
	}

	if pos.SideToMove == types.BLACK {
		pos.zobristUpdateColor()
	}

	for _, castling := range []Castling{WHITE_CASTLING_KING, WHITE_CASTLING_QUEEN, BLACK_CASTLING_KING, BLACK_CASTLING_QUEEN} {
		if pos.castling|castling != 0 {
			pos.zobristUpdateCastling(castling)
		}
	}

	if pos.enPassant != types.SQUARE_NONE {
		pos.ZobristHash ^= z.enPassant[types.FileOfSquare(pos.enPassant)]
	}
}

func (pos *Position) zobristUpdateColor() {
	pos.ZobristHash ^= z.sideToMoveIsBlack
}

func (pos *Position) zobristUpdateCastling(c Castling) {
	pos.ZobristHash ^= z.castling[bits.TrailingZeros(uint(c))]
}

func (pos *Position) zobristUpdatePiece(square int, color types.Color, pieceType types.PieceType) {
	pos.ZobristHash ^= z.piecesOnSquares[square][color][pieceType]
}

func (pos *Position) zobristUpdateEnPassant(square int) {
	pos.ZobristHash ^= z.enPassant[types.FileOfSquare(square)]
}
