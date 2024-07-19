package position

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/nnue"
	"github.com/shaardie/clemens/pkg/types"
)

// GetPiece returns the Piece from the square
func (pos *Position) GetPiece(square uint8) types.Piece {
	return pos.PiecesBoard[square]
}

// SetPiece adds a pieces to the given square
func (pos *Position) SetPiece(p types.Piece, square uint8) {
	pos.PiecesBoard[square] = p
	c := p.Color()
	t := p.Type()
	pos.PiecesBitboard[c][t] |= bitboard.BitBySquares(square)

	// Update zobrist Hash
	pos.zobristUpdatePiece(square, c, t)

	// Update features
	if !pos.needsRefresh[c] {
		if t == types.KING {
			pos.needsRefresh[c] = true
		} else {
			kingSquare := bitboard.LeastSignificantOneBit(pos.PiecesBitboard[c][types.KING])
			pos.removedFeatures = append(pos.addedFeatures, nnue.CalculateIdx(kingSquare, square, p))
		}
	}
}

// DeletePiece deletes the piece on the given square
func (pos *Position) DeletePiece(square uint8) types.Piece {
	// Get Piece from pieceBoard
	p := pos.PiecesBoard[square]
	c := p.Color()
	t := p.Type()

	// Remove Piece from pieceBoard
	pos.PiecesBoard[square] = types.NO_PIECE

	// Remove Piece from Bitboard by generating the difference
	pos.PiecesBitboard[p.Color()][p.Type()] &= ^bitboard.BitBySquares(square)

	// Update zobrist Hash
	pos.zobristUpdatePiece(square, c, t)

	// Update features
	if !pos.needsRefresh[c] {
		if t == types.KING {
			pos.needsRefresh[c] = true
		} else {
			kingSquare := bitboard.LeastSignificantOneBit(pos.PiecesBitboard[c][types.KING])
			pos.removedFeatures = append(pos.removedFeatures, nnue.CalculateIdx(kingSquare, square, p))
		}
	}
	return p
}

func (pos *Position) MovePiece(fromSquare, toSquare uint8) types.Piece {
	p := pos.DeletePiece(fromSquare)
	pos.SetPiece(p, toSquare)
	return p
}
