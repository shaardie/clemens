package position

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/nnue"
	"github.com/shaardie/clemens/pkg/pieces/bishop"
	"github.com/shaardie/clemens/pkg/pieces/king"
	"github.com/shaardie/clemens/pkg/pieces/knight"
	"github.com/shaardie/clemens/pkg/pieces/pawn"
	"github.com/shaardie/clemens/pkg/pieces/rook"
	"github.com/shaardie/clemens/pkg/types"
)

type Position struct {
	// Array of bitboards for all pieces
	PiecesBitboard   [types.COLOR_NUMBER][types.PIECE_TYPE_NUMBER]bitboard.Bitboard
	ZobristHash      uint64
	AllPieces        bitboard.Bitboard
	AllPiecesByColor [types.COLOR_NUMBER]bitboard.Bitboard
	// Array of Pieces on the Board
	PiecesBoard [types.SQUARE_NUMBER]types.Piece
	// Color of the side to move
	SideToMove types.Color
	// Castling possibilities
	Castling Castling
	// En passant square
	EnPassant     uint8
	HalfMoveClock uint8
	Ply           uint8

	Accumulator     nnue.Accumulator
	addedFeatures   []int
	removedFeatures []int
	needsRefresh    [types.COLOR_NUMBER]bool
}

func New() *Position {
	pos := &Position{
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
		SideToMove: types.WHITE,
		Castling:   WHITE_CASTLING_KING | WHITE_CASTLING_QUEEN | BLACK_CASTLING_QUEEN | BLACK_CASTLING_KING,
		EnPassant:  types.SQUARE_NONE,
		Ply:        0,
	}
	pos.boardToBitBoard()
	pos.generateHelperBitboards()

	// Create initial zobrist hash
	pos.initZobristHash()

	// Init accumulator
	pos.Accumulator.Refresh(pos.activeFeatures(), pos.SideToMove)

	return pos
}

// SquareAttackedBy returns a bitboard with all pieces attacking the specified square.
// The main idea behind the implementation is to use a piece on the specified square and let it attack all other pieces with all attack pattern,
// then intercept this attacks with the pieces capable of this attack pattern.
func (pos *Position) SquareAttackedBy(square uint8) bitboard.Bitboard {
	occupied := pos.AllPieces

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
	return attacks
}

// Empty return true, if there is no piece on the square
func (pos *Position) Empty(square uint8) bool {
	return pos.GetPiece(square) == types.NO_PIECE
}

func (pos *Position) boardToBitBoard() {
	for square, piece := range pos.PiecesBoard {
		if piece == types.NO_PIECE {
			continue
		}
		pos.PiecesBitboard[piece.Color()][piece.Type()] |= bitboard.BitBySquares(uint8(square))
	}
}

func (pos *Position) generateHelperBitboards() {
	pos.AllPieces = bitboard.Empty
	for _, c := range []types.Color{types.WHITE, types.BLACK} {
		bb := bitboard.Empty
		for _, piece := range pos.PiecesBitboard[c] {
			bb |= piece
		}
		pos.AllPiecesByColor[c] = bb
		pos.AllPieces |= bb
	}
}

func (pos *Position) IsLegal() bool {
	return !pos.IsInCheck(types.SwitchColor(pos.SideToMove))
}

func (pos *Position) activeFeatures() []int {
	features := make([]int, 0, pos.AllPieces.PopulationCount())
	for c := range types.COLOR_NUMBER {
		kingSquare := bitboard.LeastSignificantOneBit(pos.PiecesBitboard[c][types.KING])
		for t := range types.KING {
			bb := pos.PiecesBitboard[c][t]
			for bb != 0 {
				square := bitboard.SquareIndexSerializationNextSquare(&bb)
				features = append(features, nnue.CalculateIdx(kingSquare, square, types.NewPiece(c, t)))
			}
		}
	}
	return features
}
