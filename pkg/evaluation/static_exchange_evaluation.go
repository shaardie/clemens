package evaluation

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/pieces/bishop"
	"github.com/shaardie/clemens/pkg/pieces/pawn"
	"github.com/shaardie/clemens/pkg/pieces/rook"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

// https://www.chessprogramming.org/SEE_-_The_Swap_Algorithm#Traversal_of_To-Attacks
func getLeastValuablePiece(pos *position.Position, attacks bitboard.Bitboard, color types.Color, attacker *types.PieceType) bitboard.Bitboard {
	for *attacker = types.PAWN; *attacker < types.PIECE_TYPE_NUMBER; *attacker++ {
		subset := attacks & pos.PiecesBitboard[color][*attacker]
		if subset > 0 {
			return subset & -subset // single bit
		}
	}
	return 0 // empty set
}

// Static Exchange Evaluation, https://www.chessprogramming.org/Static_Exchange_Evaluation
func StaticExchangeEvaluation(pos *position.Position, m *move.Move) int16 {
	targetSquare := m.GetTargetSquare()
	sourceSquare := m.GetSourceSquare()
	sourceSquareBB := bitboard.BitBySquares(sourceSquare)
	targetType := pos.PiecesBoard[targetSquare].Type()
	attackerType := pos.PiecesBoard[sourceSquare].Type()
	occupied := pos.AllPieces
	sideToMove := pos.SideToMove
	gain := [32]int16{}
	d := 0
	maxXray := pos.PiecesBitboard[types.WHITE][types.PAWN] | pos.PiecesBitboard[types.BLACK][types.PAWN] |
		pos.PiecesBitboard[types.WHITE][types.BISHOP] | pos.PiecesBitboard[types.BLACK][types.BISHOP] |
		pos.PiecesBitboard[types.WHITE][types.ROOK] | pos.PiecesBitboard[types.BLACK][types.ROOK] |
		pos.PiecesBitboard[types.WHITE][types.QUEEN] | pos.PiecesBitboard[types.BLACK][types.QUEEN]
	attacks := pos.SquareAttackedBy(targetSquare)

	gain[d] = PieceValue[targetType]
	alreadyAttacked := bitboard.BitBySquares(sourceSquare)

	for {
		d++
		gain[d] = PieceValue[attackerType] - gain[d-1] // speculative store, if defended
		if max(-gain[d-1], gain[d]) < 0 {
			break // pruning does not influence the result
		}
		attacks ^= sourceSquareBB  // reset bit in set to traverse
		occupied ^= sourceSquareBB // reset bit in temporary occupancy (for x-Rays)
		alreadyAttacked |= sourceSquareBB
		if sourceSquareBB&maxXray > 0 {
			attacks |= considerXrays(pos, targetSquare, occupied, alreadyAttacked)
		}
		sideToMove = types.SwitchColor(sideToMove)
		sourceSquareBB = getLeastValuablePiece(pos, attacks, sideToMove, &attackerType)
		if sourceSquareBB == bitboard.Empty {
			break
		}
	}
	for {
		d--
		if d == 0 {
			break
		}
		gain[d-1] = -max(-gain[d-1], gain[d])
	}
	return gain[0]
}

func considerXrays(pos *position.Position, square uint8, occupied, alreadyAttacked bitboard.Bitboard) bitboard.Bitboard {
	attacks := bitboard.Empty

	// Diagonal attacks
	diagonalSlider := pos.PiecesBitboard[types.WHITE][types.BISHOP] |
		pos.PiecesBitboard[types.BLACK][types.BISHOP] |
		pos.PiecesBitboard[types.WHITE][types.QUEEN] |
		pos.PiecesBitboard[types.BLACK][types.QUEEN]
	attacks |= bishop.AttacksBySquare(square, occupied) & diagonalSlider

	// Vertical attacks
	verticalAndHorizonalSlider := pos.PiecesBitboard[types.WHITE][types.ROOK] |
		pos.PiecesBitboard[types.BLACK][types.ROOK] |
		pos.PiecesBitboard[types.WHITE][types.QUEEN] |
		pos.PiecesBitboard[types.BLACK][types.QUEEN]
	attacks |= rook.AttacksBySquare(square, occupied) & verticalAndHorizonalSlider

	// Pawn attacks, we need to switch color to emuluate that
	attacks |= pawn.AttacksBySquare(types.WHITE, square) & pos.PiecesBitboard[types.BLACK][types.PAWN]
	attacks |= pawn.AttacksBySquare(types.BLACK, square) & pos.PiecesBitboard[types.WHITE][types.PAWN]
	return attacks &^ alreadyAttacked
}
