package position

import (
	"fmt"

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
	PiecesBoard [types.SQUARE_NUMBER]types.Piece
	// Array of bitboards for all pieces
	PiecesBitboard [types.COLOR_NUMBER][types.PIECE_TYPE_NUMBER]bitboard.Bitboard

	SideToMove types.Color
	// TODO, En Passante, Halfmove Clock
	castling          Castling
	enPassante        int
	halfMoveClock     int
	numberOfFullMoves int
}

func New() Position {
	pos := Position{
		PiecesBoard: [types.SQUARE_NUMBER]types.Piece{
			types.WHITE_ROOK, types.WHITE_KNIGHT, types.WHITE_BISHOP, types.WHITE_QUEEN, types.WHITE_KING, types.WHITE_BISHOP, types.WHITE_KNIGHT, types.WHITE_ROOK,
			types.WHITE_PAWN, types.WHITE_PAWN, types.WHITE_PAWN, types.WHITE_PAWN, types.WHITE_PAWN, types.WHITE_PAWN, types.WHITE_PAWN, types.WHITE_PAWN,
			types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE,
			types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE,
			types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE,
			types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE, types.NO_PIECE,
			types.BLACK_PAWN, types.BLACK_PAWN, types.BLACK_PAWN, types.BLACK_PAWN, types.BLACK_PAWN, types.BLACK_PAWN, types.BLACK_PAWN, types.BLACK_PAWN,
			types.BLACK_ROOK, types.BLACK_KNIGHT, types.BLACK_BISHOP, types.BLACK_QUEEN, types.BLACK_KING, types.BLACK_BISHOP, types.BLACK_KNIGHT, types.BLACK_ROOK,
		},
		SideToMove:        types.WHITE,
		castling:          WHITE_CASTLING_KING | WHITE_CASTLING_QUEEN | BLACK_CASTLING_QUEEN | BLACK_CASTLING_KING,
		enPassante:        types.SQUARE_NONE,
		numberOfFullMoves: 1,
	}
	pos.boardToBitBoard()
	return pos
}

// SquareAttackedBy returns a bitboard with all pieces attacking the specified square.
// The main idea behind the implementation is to use a piece on the specified square and let it attack all other pieces with all attack pattern,
// then intercept this attacks with the pieces capable of this attack pattern.
func (pos Position) SquareAttackedBy(square int, occupied bitboard.Bitboard) bitboard.Bitboard {
	// Knight attacks
	knights := pos.PiecesBitboard[types.WHITE][types.KNIGHT] | pos.PiecesBitboard[types.BLACK][types.KNIGHT]
	attacks := knight.AttacksBySquare(square) & knights

	// King attacks
	kings := pos.PiecesBitboard[types.WHITE][types.KING] | pos.PiecesBitboard[types.BLACK][types.KING]
	attacks |= king.AttacksBySquare(square) & kings

	// Diagonal attacks
	diagonalSlider := pos.PiecesBitboard[types.WHITE][types.BISHOP] | pos.PiecesBitboard[types.BLACK][types.BISHOP] | pos.PiecesBitboard[types.WHITE][types.QUEEN] | pos.PiecesBitboard[types.BLACK][types.QUEEN]
	attacks |= bishop.AttacksBySquare(square, occupied) & diagonalSlider

	// Vertical attacks
	verticalAndHorizonalSlider := pos.PiecesBitboard[types.WHITE][types.ROOK] | pos.PiecesBitboard[types.BLACK][types.ROOK] | pos.PiecesBitboard[types.WHITE][types.QUEEN] | pos.PiecesBitboard[types.BLACK][types.QUEEN]
	attacks |= rook.AttacksBySquare(square, occupied) & verticalAndHorizonalSlider

	// Pawn attacks, we need to switch color to emuluate that
	attacks |= pawn.AttacksBySquare(types.WHITE, square) & pos.PiecesBitboard[types.BLACK][types.PAWN]
	attacks |= pawn.AttacksBySquare(types.BLACK, square) & pos.PiecesBitboard[types.WHITE][types.PAWN]
	return 0
}

// Empty return true, if there is no piece on the square
func (pos Position) Empty(square int) bool {
	return pos.GetPieceFromSquare(square) == types.NO_PIECE
}

// GetPieceFromSquare returns the Piece from the square
func (pos Position) GetPieceFromSquare(square int) types.Piece {
	return pos.PiecesBoard[square]
}

// GetPieceFromSquare returns the Piece from the square
func (pos Position) CanCastle(c Castling) bool {
	return c&pos.castling != NO_CASTLING
}

func (pos Position) SetPiece(p types.Piece, square int) {
	pos.PiecesBoard[square] = p
	pos.PiecesBitboard[p.Color()][p.Type()] |= bitboard.BitBySquares(square)
}

func (pos Position) boardToBitBoard() {
	for square, piece := range pos.PiecesBoard {
		if piece == types.NO_PIECE {
			continue
		}
		pos.PiecesBitboard[piece.Color()][piece.Type()] |= bitboard.BitBySquares(square)
	}
}

func (pos Position) validate() error {
	// Validate Pieces
	for color, bb := range pos.PiecesBitboard {
		for pieceType, b := range bb {
			idxs := bitboard.SquareIndexSerialization(b)
			fmt.Println(idxs)
			for _, idx := range idxs {
				p := pos.PiecesBoard[idx]
				if p == types.NO_PIECE {
					return fmt.Errorf("no piece on %v", idx)
				}
				if p.Color() != types.Color(color) {
					return fmt.Errorf("piece on %v has different color, board=%v, bitboard=%v", idx, p.Color(), color)
				}
				if p.Type() != types.PieceType(pieceType) {
					return fmt.Errorf("piece on %v has different type, board=%v, bitboard=%v", idx, p.Type(), pieceType)
				}
			}
		}
	}

	return nil
}
