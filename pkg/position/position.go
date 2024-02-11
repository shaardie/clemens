package position

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/pieces/bishop"
	"github.com/shaardie/clemens/pkg/pieces/king"
	"github.com/shaardie/clemens/pkg/pieces/knight"
	"github.com/shaardie/clemens/pkg/pieces/pawn"
	"github.com/shaardie/clemens/pkg/pieces/rook"
	"github.com/shaardie/clemens/pkg/types"
)

type Position struct {
	// Array of Pieces on the Board
	piecesBoard [types.SQUARE_NUMBER]types.Piece
	// Array of bitboards for all pieces
	piecesBitboard [types.COLOR_NUMBER][types.PIECE_TYPE_NUMBER]bitboard.Bitboard
	// Color of the side to move
	SideToMove types.Color
	// Castling possibilities
	castling Castling
	// En passant square
	enPassant     int
	HalfMoveClock uint8
	ply           int

	ZobristHash      uint64
	allPieces        bitboard.Bitboard
	allPiecesByColor [types.COLOR_NUMBER]bitboard.Bitboard
}

func New() *Position {
	pos := &Position{
		piecesBoard: [types.SQUARE_NUMBER]types.Piece{
			types.WHITE_ROOK, types.WHITE_KNIGHT, types.WHITE_BISHOP, types.WHITE_QUEEN, types.WHITE_KING, types.WHITE_BISHOP, types.WHITE_KNIGHT, types.WHITE_ROOK,
			types.WHITE_PAWN, types.WHITE_PAWN, types.WHITE_PAWN, types.WHITE_PAWN, types.WHITE_PAWN, types.WHITE_PAWN, types.WHITE_PAWN, types.WHITE_PAWN,
			types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE,
			types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE,
			types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE,
			types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE,
			types.BLACK_PAWN, types.BLACK_PAWN, types.BLACK_PAWN, types.BLACK_PAWN, types.BLACK_PAWN, types.BLACK_PAWN, types.BLACK_PAWN, types.BLACK_PAWN,
			types.BLACK_ROOK, types.BLACK_KNIGHT, types.BLACK_BISHOP, types.BLACK_QUEEN, types.BLACK_KING, types.BLACK_BISHOP, types.BLACK_KNIGHT, types.BLACK_ROOK,
		},
		SideToMove: types.WHITE,
		castling:   WHITE_CASTLING_KING | WHITE_CASTLING_QUEEN | BLACK_CASTLING_QUEEN | BLACK_CASTLING_KING,
		enPassant:  types.SQUARE_NONE,
		ply:        0,
	}
	pos.boardToBitBoard()
	pos.generateHelperBitboards()

	// Create initial zobrist hash
	pos.initZobristHash()

	return pos
}

// SquareAttackedBy returns a bitboard with all pieces attacking the specified square.
// The main idea behind the implementation is to use a piece on the specified square and let it attack all other pieces with all attack pattern,
// then intercept this attacks with the pieces capable of this attack pattern.
func (pos *Position) SquareAttackedBy(square int) bitboard.Bitboard {
	occupied := pos.AllPieces()

	// Knight attacks
	knights := pos.piecesBitboard[types.WHITE][types.KNIGHT] | pos.piecesBitboard[types.BLACK][types.KNIGHT]
	attacks := knight.AttacksBySquare(square) & knights

	// King attacks
	kings := pos.piecesBitboard[types.WHITE][types.KING] | pos.piecesBitboard[types.BLACK][types.KING]
	attacks |= king.AttacksBySquare(square) & kings

	// Diagonal attacks
	diagonalSlider := pos.piecesBitboard[types.WHITE][types.BISHOP] | pos.piecesBitboard[types.BLACK][types.BISHOP] | pos.piecesBitboard[types.WHITE][types.QUEEN] | pos.piecesBitboard[types.BLACK][types.QUEEN]
	attacks |= bishop.AttacksBySquare(square, occupied) & diagonalSlider

	// Vertical attacks
	verticalAndHorizonalSlider := pos.piecesBitboard[types.WHITE][types.ROOK] | pos.piecesBitboard[types.BLACK][types.ROOK] | pos.piecesBitboard[types.WHITE][types.QUEEN] | pos.piecesBitboard[types.BLACK][types.QUEEN]
	attacks |= rook.AttacksBySquare(square, occupied) & verticalAndHorizonalSlider

	// Pawn attacks, we need to switch color to emuluate that
	attacks |= pawn.AttacksBySquare(types.WHITE, square) & pos.piecesBitboard[types.BLACK][types.PAWN]
	attacks |= pawn.AttacksBySquare(types.BLACK, square) & pos.piecesBitboard[types.WHITE][types.PAWN]
	return attacks
}

// Empty return true, if there is no piece on the square
func (pos *Position) Empty(square int) bool {
	return pos.GetPiece(square) == types.NO_PIECE
}

func (pos *Position) boardToBitBoard() {
	for square, piece := range pos.piecesBoard {
		if piece == types.NO_PIECE {
			continue
		}
		pos.piecesBitboard[piece.Color()][piece.Type()] |= bitboard.BitBySquares(square)
	}
}

func (pos *Position) generateHelperBitboards() {
	pos.allPieces = bitboard.Empty
	for _, c := range []types.Color{types.WHITE, types.BLACK} {
		bb := bitboard.Empty
		for _, piece := range pos.piecesBitboard[c] {
			bb |= piece
		}
		pos.allPiecesByColor[c] = bb
		pos.allPieces |= bb
	}
}

func (pos *Position) AllPieces() bitboard.Bitboard {
	return pos.allPieces
}

func (pos *Position) AllPiecesByColor(c types.Color) bitboard.Bitboard {
	return pos.allPiecesByColor[c]
}

func (pos *Position) IsLegal() bool {
	return !pos.IsInCheck(types.SwitchColor(pos.SideToMove))
}
