package position

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

// GetPiece returns the Piece from the square
func (pos *Position) GetPiece(square int) types.Piece {
	return pos.piecesBoard[square]
}

// SetPiece adds a pieces to the given square
func (pos *Position) SetPiece(p types.Piece, square int) {
	pos.piecesBoard[square] = p
	pos.piecesBitboard[p.Color()][p.Type()] |= bitboard.BitBySquares(square)
}

// DeletePiece deletes the piece on the given square
func (pos *Position) DeletePiece(square int) types.Piece {
	// Get Piece from pieceBoard
	p := pos.piecesBoard[square]
	// Remove Piece from pieceBoard
	pos.piecesBoard[square] = types.NO_PIECE

	// Remove Piece from Bitboard by generating the difference
	pos.piecesBitboard[p.Color()][p.Type()] &= ^bitboard.BitBySquares(square)

	return p
}

func (pos *Position) MovePiece(fromSquare, toSquare int) {
	p := pos.DeletePiece(fromSquare)
	pos.SetPiece(p, toSquare)
}