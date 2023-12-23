package position

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/pieces/bishop"
	"github.com/shaardie/clemens/pkg/pieces/king"
	"github.com/shaardie/clemens/pkg/pieces/knight"
	"github.com/shaardie/clemens/pkg/pieces/pawn"
	"github.com/shaardie/clemens/pkg/pieces/rook"
	. "github.com/shaardie/clemens/pkg/types"
)

type Position struct {
	// Array of Pieces on the Board
	PiecesBoard [SQUARE_NUMBER]Piece
	// Array of bitboards for all pieces
	PiecesBitboard [COLOR_NUMBER][PIECE_TYPE_NUMBER]bitboard.Bitboard
	SideToMove     Color
	// TODO, En Passante, Halfmove Clock
	castling          Castling
	enPassante        int
	halfMoveClock     int
	numberOfFullMoves int
}

func New() Position {
	return Position{
		PiecesBoard: [SQUARE_NUMBER]Piece{
			WHITE_ROOK, WHITE_KNIGHT, WHITE_BISHOP, WHITE_QUEEN, WHITE_KING, WHITE_BISHOP, WHITE_KNIGHT, WHITE_ROOK,
			WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN, WHITE_PAWN,
			NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE,
			NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE,
			NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE,
			NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE, NO_PIECE,
			BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN, BLACK_PAWN,
			BLACK_ROOK, BLACK_KNIGHT, BLACK_BISHOP, BLACK_QUEEN, BLACK_KING, BLACK_BISHOP, BLACK_KNIGHT, BLACK_ROOK,
		},
		PiecesBitboard: [COLOR_NUMBER][PIECE_TYPE_NUMBER]bitboard.Bitboard{
			// White
			{
				// Pawns
				bitboard.BitBySquares(
					SQUARE_A2, SQUARE_B2, SQUARE_C2, SQUARE_D2,
					SQUARE_E2, SQUARE_F2, SQUARE_G2, SQUARE_H2,
				),
				// Knights
				bitboard.BitBySquares(SQUARE_B1, SQUARE_G1),
				// Bishops
				bitboard.BitBySquares(SQUARE_B1, SQUARE_G1),
				// Queens
				bitboard.BitBySquares(SQUARE_D1),
				// King
				bitboard.BitBySquares(SQUARE_E1),
			},
			// Black
			{
				// Pawns
				bitboard.BitBySquares(
					SQUARE_A7, SQUARE_B7, SQUARE_C7, SQUARE_D7,
					SQUARE_E7, SQUARE_F7, SQUARE_G7, SQUARE_H7,
				),
				// Knights
				bitboard.BitBySquares(SQUARE_B8, SQUARE_G8),
				// Bishops
				bitboard.BitBySquares(SQUARE_B8, SQUARE_G8),
				// Queens
				bitboard.BitBySquares(SQUARE_D8),
				// King
				bitboard.BitBySquares(SQUARE_E8),
			},
		},
		SideToMove:        WHITE,
		castling:          WHITE_CASTLING_KING | WHITE_CASTLING_QUEEN | BLACK_CASTLING_QUEEN | BLACK_CASTLING_KING,
		enPassante:        SQUARE_NONE,
		numberOfFullMoves: 1,
	}

}

// SquareAttackedBy returns a bitboard with all pieces attacking the specified square.
// The main idea behind the implementation is to use a piece on the specified square and let it attack all other pieces with all attack pattern,
// then intercept this attacks with the pieces capable of this attack pattern.
func (pos Position) SquareAttackedBy(square int, occupied bitboard.Bitboard) bitboard.Bitboard {
	// Knight attacks
	knights := pos.PiecesBitboard[WHITE][KNIGHT] | pos.PiecesBitboard[BLACK][KNIGHT]
	attacks := knight.AttacksBySquare(square) & knights

	// King attacks
	kings := pos.PiecesBitboard[WHITE][KING] | pos.PiecesBitboard[BLACK][KING]
	attacks |= king.AttacksBySquare(square) & kings

	// Diagonal attacks
	diagonalSlider := pos.PiecesBitboard[WHITE][BISHOP] | pos.PiecesBitboard[BLACK][BISHOP] | pos.PiecesBitboard[WHITE][QUEEN] | pos.PiecesBitboard[BLACK][QUEEN]
	attacks |= bishop.AttacksBySquare(square, occupied) & diagonalSlider

	// Vertical attacks
	verticalAndHorizonalSlider := pos.PiecesBitboard[WHITE][ROOK] | pos.PiecesBitboard[BLACK][ROOK] | pos.PiecesBitboard[WHITE][QUEEN] | pos.PiecesBitboard[BLACK][QUEEN]
	attacks |= rook.AttacksBySquare(square, occupied) & verticalAndHorizonalSlider

	// Pawn attacks, we need to switch color to emuluate that
	attacks |= pawn.AttacksBySquare(WHITE, square) & pos.PiecesBitboard[BLACK][PAWN]
	attacks |= pawn.AttacksBySquare(BLACK, square) & pos.PiecesBitboard[WHITE][PAWN]
	return 0
}

// Empty return true, if there is no piece on the square
func (pos Position) Empty(square int) bool {
	return pos.GetPieceFromSquare(square) == NO_PIECE
}

// GetPieceFromSquare returns the Piece from the square
func (pos Position) GetPieceFromSquare(square int) Piece {
	return pos.PiecesBoard[square]
}

// GetPieceFromSquare returns the Piece from the square
func (pos Position) CanCastle(c Castling) bool {
	return c&pos.castling != NO_CASTLING
}
