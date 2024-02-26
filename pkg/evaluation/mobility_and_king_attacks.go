package evaluation

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/pieces/bishop"
	"github.com/shaardie/clemens/pkg/pieces/king"
	"github.com/shaardie/clemens/pkg/pieces/knight"
	"github.com/shaardie/clemens/pkg/pieces/pawn"
	"github.com/shaardie/clemens/pkg/pieces/queen"
	"github.com/shaardie/clemens/pkg/pieces/rook"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	// scalar adjustments for attacking squares near the king
	kingAttValue = [types.PIECE_TYPE_NUMBER]int{1, 2, 2, 3, 4, 1}
)

func (e *eval) evalMobilityAndKingAttackValue(pos *position.Position) {
	e.baseScore += evalMobilityAndKingAttackValueByColor(pos, types.WHITE) - evalMobilityAndKingAttackValueByColor(pos, types.BLACK)
}

func evalMobilityAndKingAttackValueByColor(pos *position.Position, we types.Color) int {
	var mobility bitboard.Bitboard
	var square int
	var val int
	them := types.SwitchColor(we)
	destination := ^pos.AllPiecesByColor[them]
	kingSquares := king.AttacksBySquare(bitboard.LeastSignificantOneBit(pos.PiecesBitboard[we][types.KING]))
	for pt := types.PAWN; pt < types.PIECE_TYPE_NUMBER; pt++ {
		it := bitboard.SquareIndexSerializationIterator(pos.PiecesBitboard[them][pt])
		for {
			square = it()
			if square == types.SQUARE_NONE {
				break
			}
			switch pt {
			case types.PAWN:
				val += (pawn.PushesBySquare(we, square, pos.AllPieces) & destination).PopulationCount()
				mobility = pawn.AttacksBySquare(we, square)
			case types.BISHOP:
				mobility = bishop.AttacksBySquare(square, pos.AllPieces)
			case types.KNIGHT:
				mobility = knight.AttacksBySquare(square)
			case types.ROOK:
				mobility = rook.AttacksBySquare(square, pos.AllPieces)
			case types.QUEEN:
				mobility = queen.AttacksBySquare(square, pos.AllPieces)
			case types.KING:
				mobility = king.AttacksBySquare(square)
			}
			// Bonus for mobility
			mobility &= destination
			val += mobility.PopulationCount()

			// Bonus for pieces attacking the squares next to the king
			val += kingAttValue[pt] * (mobility & kingSquares).PopulationCount()
		}
	}
	return val
}
