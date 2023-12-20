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
	// Array of bitboards for all pieces
	PiecesBB   [COLOR_NUMBER][PIECE_NUMBER]bitboard.Bitboard
	SideToMove Color
	// TODO Castling, En Passante, Halfmove Clock
}

func New() Position {
	return Position{
		PiecesBB: [COLOR_NUMBER][PIECE_NUMBER]bitboard.Bitboard{
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
		SideToMove: WHITE,
	}

}

// SquareAttackedBy returns a bitboard with all pieces attacking the specified square.
// The main idea behind the implementation is to use a piece on the specified square and let it attack all other pieces with all attack pattern,
// then intercept this attacks with the pieces capable of this attack pattern.
func (pos Position) SquareAttackedBy(square int, occupied bitboard.Bitboard) bitboard.Bitboard {
	// Knight attacks
	knights := pos.PiecesBB[WHITE][KNIGHT] | pos.PiecesBB[BLACK][KNIGHT]
	attacks := knight.AttacksBySquare(square) & knights

	// King attacks
	kings := pos.PiecesBB[WHITE][KING] | pos.PiecesBB[BLACK][KING]
	attacks |= king.AttacksBySquare(square) & kings

	// Diagonal attacks
	diagonalSlider := pos.PiecesBB[WHITE][BISHOP] | pos.PiecesBB[BLACK][BISHOP] | pos.PiecesBB[WHITE][QUEEN] | pos.PiecesBB[BLACK][QUEEN]
	attacks |= bishop.AttacksBySquare(square, occupied) & diagonalSlider

	// Vertical attacks
	verticalAndHorizonalSlider := pos.PiecesBB[WHITE][ROOK] | pos.PiecesBB[BLACK][ROOK] | pos.PiecesBB[WHITE][QUEEN] | pos.PiecesBB[BLACK][QUEEN]
	attacks |= rook.AttacksBySquare(square, occupied) & verticalAndHorizonalSlider

	// Pawn attacks, we need to switch color to emuluate that
	attacks |= pawn.AttacksBySquare(WHITE, square) & pos.PiecesBB[BLACK][PAWN]
	attacks |= pawn.AttacksBySquare(BLACK, square) & pos.PiecesBB[WHITE][PAWN]
	return 0
}
