package position

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/pieces/bishop"
	"github.com/shaardie/clemens/pkg/pieces/king"
	"github.com/shaardie/clemens/pkg/pieces/knight"
	"github.com/shaardie/clemens/pkg/pieces/pawn"
	"github.com/shaardie/clemens/pkg/pieces/rook"
	"github.com/shaardie/clemens/pkg/types"
)

func (pos *Position) GeneratePseudoLegalMoves() []move.Move {
	moves := []move.Move{}

	occupied := pos.AllPieces()
	destinations := ^pos.AllPiecesByColor(pos.sideToMove)

	// Sliding Pieces
	moves = append(moves,
		generateMoves(
			pos.piecesBitboard[pos.sideToMove][types.ROOK],
			occupied,
			destinations,
			rook.AttacksBySquare,
		)...,
	)
	moves = append(moves,
		generateMoves(
			pos.piecesBitboard[pos.sideToMove][types.BISHOP],
			occupied,
			destinations,
			bishop.AttacksBySquare,
		)...,
	)
	moves = append(moves,
		generateMoves(
			pos.piecesBitboard[pos.sideToMove][types.QUEEN],
			occupied,
			destinations,
			bishop.AttacksBySquare,
		)...,
	)

	// Pieces ignoring occupation
	moves = append(moves,
		generateMoves(
			pos.piecesBitboard[pos.sideToMove][types.KNIGHT],
			bitboard.Empty,
			destinations,
			func(square int, _ bitboard.Bitboard) bitboard.Bitboard {
				return knight.AttacksBySquare(square)
			},
		)...,
	)
	moves = append(moves,
		generateMoves(
			pos.piecesBitboard[pos.sideToMove][types.KING],
			bitboard.Empty,
			destinations,
			func(square int, _ bitboard.Bitboard) bitboard.Bitboard {
				return king.AttacksBySquare(square)
			},
		)...,
	)

	// Pawns
	for _, sourceSquare := range bitboard.SquareIndexSerialization(pos.piecesBitboard[pos.sideToMove][types.PAWN]) {
		// Pushes
		for _, targetSquare := range bitboard.SquareIndexSerialization(pawn.PushesBySquare(pos.sideToMove, sourceSquare, occupied)) {
			var m move.Move
			m.SetSourceSquare(sourceSquare)
			m.SetDestinationSquare(targetSquare)
			moves = append(moves, m)
		}
		// Attacks
		for _, targetSquare := range bitboard.SquareIndexSerialization(pawn.AttacksBySquare(pos.sideToMove, sourceSquare) & pos.AllPiecesByColor(types.SwitchColor(pos.sideToMove))) {
			var m move.Move
			m.SetSourceSquare(sourceSquare)
			m.SetDestinationSquare(targetSquare)
			moves = append(moves, m)
		}
		// En Passant
		if pos.enPassant != types.SQUARE_NONE {
			// PseudoPiece on square behind the en passant pawn
			var pseudoPieceSquare int
			switch pos.sideToMove {
			case types.WHITE:
				pseudoPieceSquare = pos.enPassant + types.FILE_NUMBER
			case types.BLACK:
				pseudoPieceSquare = pos.enPassant - types.FILE_NUMBER
			}

			for _, targetSquare := range bitboard.SquareIndexSerialization(pawn.AttacksBySquare(pos.sideToMove, pseudoPieceSquare)) {
				var m move.Move
				m.SetSourceSquare(sourceSquare)
				m.SetDestinationSquare(targetSquare)
				m.SetMoveType(move.EN_PASSANT)
				moves = append(moves, m)
			}
		}
	}

	// TODO castling
	return moves
}

func generateMoves(sources, occupied, destinations bitboard.Bitboard, attacks func(square int, occupied bitboard.Bitboard) bitboard.Bitboard) []move.Move {
	moves := []move.Move{}

	for _, sourceSquare := range bitboard.SquareIndexSerialization(sources) {
		for _, targetSquare := range bitboard.SquareIndexSerialization(attacks(sourceSquare, occupied) & destinations) {
			var m move.Move
			m.SetSourceSquare(sourceSquare)
			m.SetDestinationSquare(targetSquare)
			moves = append(moves, m)
		}
	}
	return moves
}
