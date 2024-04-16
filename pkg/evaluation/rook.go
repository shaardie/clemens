package evaluation

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

var (
	rookSeventhMidgame int16 = 20
	rookSeventhEndgame int16 = 30
	rookOpenFile       int16 = 10
	rookHalfOpenFile   int16 = 5
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
	rooks := pos.PiecesBitboard[types.WHITE][types.ROOK]
	for rooks != 0 {
		numberOfRooks[types.WHITE]++
		fileMask := bitboard.FileMaskOfSquare(bitboard.SquareIndexSerializationNextSquare(&rooks))
		if pos.PiecesBitboard[types.BLACK][types.PAWN]&fileMask > 0 {
			if pos.PiecesBitboard[types.WHITE][types.PAWN]&fileMask > 0 {
				e.phaseScores[midgame] += rookOpenFile
				continue
			}
			e.phaseScores[endgame] += rookHalfOpenFile
		}
	}
	rooks = pos.PiecesBitboard[types.BLACK][types.ROOK]
	for rooks != 0 {
		numberOfRooks[types.BLACK]++
		fileMask := bitboard.FileMaskOfSquare(bitboard.SquareIndexSerializationNextSquare(&rooks))
		if pos.PiecesBitboard[types.WHITE][types.PAWN]&fileMask > 0 {
			if pos.PiecesBitboard[types.BLACK][types.PAWN]&fileMask > 0 {
				e.phaseScores[midgame] -= rookOpenFile
				continue
			}
			e.phaseScores[endgame] -= rookHalfOpenFile
		}
	}
}
