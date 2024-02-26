package evaluation

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	rookSeventhMidgame = 20
	rookSeventhEndgame = 30
	rookOpenFile       = 10
	rookHalfOpenFile   = 5
)

func (e *eval) evalRooks(pos *position.Position) {
	var numberOfRooks [types.COLOR_NUMBER]int

	// Rook Evaluation
	// Rook on the seventh with either the king on the eigth or enemy pawns on the seven
	// https://www.chessprogramming.org/Rook_on_Seventh
	if pos.PiecesBitboard[types.WHITE][types.ROOK]&bitboard.RankMask7 > 0 && (pos.PiecesBitboard[types.BLACK][types.KING]&bitboard.RankMask8 > 0 ||
		pos.PiecesBitboard[types.BLACK][types.PAWN]&bitboard.RankMask7 > 0) {
		e.phaseScores[midgame] += rookSeventhMidgame
		e.phaseScores[endgame] += rookSeventhEndgame
	}
	// vice versa
	if pos.PiecesBitboard[types.BLACK][types.ROOK]&bitboard.RankMask2 > 0 && (pos.PiecesBitboard[types.WHITE][types.KING]&bitboard.RankMask1 > 0 ||
		pos.PiecesBitboard[types.WHITE][types.PAWN]&bitboard.RankMask2 > 0) {
		e.phaseScores[midgame] -= rookSeventhMidgame
		e.phaseScores[endgame] -= rookSeventhEndgame
	}

	// Rooks are on an open or half open file
	// https://www.chessprogramming.org/Rook_on_Open_File
	it := bitboard.SquareIndexSerializationIterator(pos.PiecesBitboard[types.WHITE][types.ROOK])
	for {
		square := it()
		if square == types.SQUARE_NONE {
			break
		}
		numberOfRooks[types.WHITE]++
		fileMask := bitboard.FileMaskOfSquare(square)
		if pos.PiecesBitboard[types.BLACK][types.PAWN]&fileMask > 0 {
			if pos.PiecesBitboard[types.WHITE][types.PAWN]&fileMask > 0 {
				e.phaseScores[midgame] += rookOpenFile
				continue
			}
			e.phaseScores[endgame] += rookHalfOpenFile
		}
	}
	it = bitboard.SquareIndexSerializationIterator(pos.PiecesBitboard[types.BLACK][types.ROOK])
	for {
		square := it()
		if square == types.SQUARE_NONE {
			break
		}
		numberOfRooks[types.BLACK]++
		fileMask := bitboard.FileMaskOfSquare(square)
		if pos.PiecesBitboard[types.WHITE][types.PAWN]&fileMask > 0 {
			if pos.PiecesBitboard[types.BLACK][types.PAWN]&fileMask > 0 {
				e.phaseScores[midgame] -= rookOpenFile
				continue
			}
			e.phaseScores[endgame] -= rookHalfOpenFile
		}
	}
}
