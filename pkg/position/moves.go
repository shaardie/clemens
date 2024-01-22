package position

import (
	"errors"

	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/pieces/bishop"
	"github.com/shaardie/clemens/pkg/pieces/king"
	"github.com/shaardie/clemens/pkg/pieces/knight"
	"github.com/shaardie/clemens/pkg/pieces/pawn"
	"github.com/shaardie/clemens/pkg/pieces/queen"
	"github.com/shaardie/clemens/pkg/pieces/rook"
	"github.com/shaardie/clemens/pkg/types"
)

// GeneratePseudoLegalMoves generates all pseudo legal moves
func (pos *Position) GeneratePseudoLegalMoves() []move.Move {
	moves := [256]move.Move{}

	occupied := pos.AllPieces()
	destinations := ^pos.AllPiecesByColor(pos.sideToMove)
	idx := 0

	// Sliding Pieces
	idx = idx + generateMovesHelper(
		moves[idx:],
		pos.piecesBitboard[pos.sideToMove][types.ROOK],
		occupied,
		destinations,
		rook.AttacksBySquare,
	)
	idx = idx + generateMovesHelper(
		moves[idx:],
		pos.piecesBitboard[pos.sideToMove][types.BISHOP],
		occupied,
		destinations,
		bishop.AttacksBySquare,
	)

	idx = idx + generateMovesHelper(
		moves[idx:],
		pos.piecesBitboard[pos.sideToMove][types.QUEEN],
		occupied,
		destinations,
		queen.AttacksBySquare,
	)

	// Pieces ignoring occupation
	idx = idx + generateMovesHelper(
		moves[idx:],
		pos.piecesBitboard[pos.sideToMove][types.KNIGHT],
		bitboard.Empty,
		destinations,
		func(square int, _ bitboard.Bitboard) bitboard.Bitboard {
			return knight.AttacksBySquare(square)
		},
	)
	idx = idx + generateMovesHelper(
		moves[idx:],
		pos.piecesBitboard[pos.sideToMove][types.KING],
		bitboard.Empty,
		destinations,
		func(square int, _ bitboard.Bitboard) bitboard.Bitboard {
			return king.AttacksBySquare(square)
		},
	)

	// Pawns
	for _, sourceSquare := range bitboard.SquareIndexSerialization(pos.piecesBitboard[pos.sideToMove][types.PAWN]) {
		// Pushes
		for _, targetSquare := range bitboard.SquareIndexSerialization(pawn.PushesBySquare(pos.sideToMove, sourceSquare, occupied)) {
			var m move.Move
			m.SetSourceSquare(sourceSquare)
			m.SetDestinationSquare(targetSquare)
			// Pawn Moves with optional Promotion
			idx = idx + pawnMoveWithPromotion(moves[idx:], pos.sideToMove, sourceSquare, targetSquare)
		}
		// Attacks
		for _, targetSquare := range bitboard.SquareIndexSerialization(pawn.AttacksBySquare(pos.sideToMove, sourceSquare) & pos.AllPiecesByColor(types.SwitchColor(pos.sideToMove))) {
			var m move.Move
			m.SetSourceSquare(sourceSquare)
			m.SetDestinationSquare(targetSquare)
			// Pawn Moves with optional Promotion
			idx = idx + pawnMoveWithPromotion(moves[idx:], pos.sideToMove, sourceSquare, targetSquare)
		}
		// En Passant
		if pos.enPassant != types.SQUARE_NONE {
			attackingPawns := pawn.AttacksBySquare(pos.sideToMove, sourceSquare) & bitboard.BitBySquares(pos.enPassant)
			for _, targetSquare := range bitboard.SquareIndexSerialization(attackingPawns) {
				var m move.Move
				m.SetSourceSquare(sourceSquare)
				m.SetDestinationSquare(targetSquare)
				m.SetMoveType(move.EN_PASSANT)
				moves[idx] = m
				idx++
			}
		}
	}

	// Castling
	for _, c := range []Castling{WHITE_CASTLING_KING, WHITE_CASTLING_QUEEN, BLACK_CASTLING_KING, BLACK_CASTLING_QUEEN} {
		if c.Color() != pos.sideToMove {
			continue
		}
		if !pos.CanCastleNow(c) {
			continue
		}
		var m move.Move
		m.SetMoveType(move.CASTLING)
		sourceSquare := bitboard.SquareIndexSerialization(pos.piecesBitboard[pos.sideToMove][types.KING])[0]
		var targetSquare int
		switch c.Side() {
		case CASTLING_KING:
			targetSquare = sourceSquare + 2
		case CASTLING_QUEEN:
			targetSquare = sourceSquare - 2
		}
		m.SetSourceSquare(sourceSquare)
		m.SetDestinationSquare(targetSquare)
		moves[idx] = m
		idx++
	}
	return moves[:idx]
}

func (pos *Position) MakeMove(m move.Move) {
	previousSideToMove := pos.sideToMove
	pos.sideToMove = types.SwitchColor(pos.sideToMove)
	if previousSideToMove == types.BLACK {
		pos.numberOfFullMoves = pos.numberOfFullMoves + 1
	}
	pos.enPassant = types.SQUARE_NONE

	sourceSquare := m.GetSourceSquare()
	targetSquare := m.GetDestinationSquare()

	targetPiece := pos.GetPiece(targetSquare)
	if targetPiece != types.NO_PIECE {
		pos.DeletePiece(targetSquare)
	}

	for _, s := range []int{sourceSquare, targetSquare} {
		switch s {
		case types.SQUARE_A1:
			pos.castling = pos.castling &^ WHITE_CASTLING_QUEEN
		case types.SQUARE_H1:
			pos.castling = pos.castling &^ WHITE_CASTLING_KING
		case types.SQUARE_A8:
			pos.castling = pos.castling &^ BLACK_CASTLING_QUEEN
		case types.SQUARE_H8:
			pos.castling = pos.castling &^ BLACK_CASTLING_KING
		case types.SQUARE_E1:
			pos.castling = pos.castling &^ (WHITE_CASTLING_QUEEN | WHITE_CASTLING_KING)
		case types.SQUARE_E8:
			pos.castling = pos.castling &^ (BLACK_CASTLING_QUEEN | BLACK_CASTLING_KING)
		}
	}

	piece := pos.MovePiece(sourceSquare, targetSquare)
	switch piece.Type() {
	// Set en passant
	case types.PAWN:
		if abs(sourceSquare-targetSquare) == 2*types.FILE_NUMBER {
			pos.enPassant = targetSquare
			switch previousSideToMove {
			case types.BLACK:
				pos.enPassant += types.FILE_NUMBER
			case types.WHITE:
				pos.enPassant -= types.FILE_NUMBER
			default:
				panic("unknown color")
			}
		}
	}

	switch m.GetMoveType() {
	case move.CASTLING:
		switch targetSquare {
		case types.SQUARE_C1:
			pos.MovePiece(types.SQUARE_A1, types.SQUARE_D1)
		case types.SQUARE_G1:
			pos.MovePiece(types.SQUARE_H1, types.SQUARE_F1)
		case types.SQUARE_C8:
			pos.MovePiece(types.SQUARE_A8, types.SQUARE_D8)
		case types.SQUARE_G8:
			pos.MovePiece(types.SQUARE_H8, types.SQUARE_F8)
		default:
			panic("wrong source square for castling")
		}
	case move.EN_PASSANT:
		// Remove pawn behind moved pawn
		var pawnToRemoveSquare = 0
		switch previousSideToMove {
		case types.WHITE:
			pawnToRemoveSquare = targetSquare - types.FILE_NUMBER
		case types.BLACK:
			pawnToRemoveSquare = targetSquare + types.FILE_NUMBER
		}
		pos.DeletePiece(pawnToRemoveSquare)
	case move.PROMOTION:
		// Promote piece
		pos.DeletePiece(targetSquare)
		pos.SetPiece(types.NewPiece(previousSideToMove, m.GetPromitionPieceType()), targetSquare)
	}
}

// generateMovesHelper generates a list of moves from a given list of paramters
func generateMovesHelper(moves []move.Move, sources, occupied, destinations bitboard.Bitboard, attacks func(square int, occupied bitboard.Bitboard) bitboard.Bitboard) int {
	idx := 0
	for _, sourceSquare := range bitboard.SquareIndexSerialization(sources) {
		for _, targetSquare := range bitboard.SquareIndexSerialization(attacks(sourceSquare, occupied) & destinations) {
			var m move.Move
			m.SetSourceSquare(sourceSquare)
			m.SetDestinationSquare(targetSquare)
			moves[idx] = m
			idx++
		}
	}
	return idx
}

func (pos *Position) MakeMoveFromString(s string) error {
	var m move.Move
	if len(s) < 4 {
		return errors.New("input to small")
	}

	sourceSquare, err := types.SquareFromString(s[0:2])
	if err != nil {
		return err
	}
	m.SetSourceSquare(sourceSquare)

	destinationSquare, err := types.SquareFromString(s[2:4])
	if err != nil {
		return err
	}
	m.SetDestinationSquare(destinationSquare)

	pt := pos.GetPiece(sourceSquare).Type()
	if pt == types.KING && abs(sourceSquare-destinationSquare) == 2 {
		m.SetMoveType(move.CASTLING)
	} else if pt == types.PAWN {
		if types.FileOfSquare(sourceSquare) != types.FileOfSquare(destinationSquare) && pos.GetPiece(destinationSquare) == types.NO_PIECE {
			m.SetMoveType(move.EN_PASSANT)
		} else if len(s) == 5 {
			m.SetMoveType(move.PROMOTION)
			pt, err := types.PieceTypeFromString(string(s[4]))
			if err != nil {
				return err
			}
			m.SetPromitionPieceType(pt)
		}
	}

	pos.MakeMove(m)
	return nil
}

func pawnMoveWithPromotion(moves []move.Move, sideToMove types.Color, sourceSquare, targetSquare int) int {
	idx := 0
	var m move.Move
	m.SetSourceSquare(sourceSquare)
	m.SetDestinationSquare(targetSquare)
	// No promotion
	if sideToMove == types.WHITE && types.RankOfSquare(targetSquare) != types.RANK_8 {
		moves[idx] = m
		return 1
	}
	if sideToMove == types.BLACK && types.RankOfSquare(targetSquare) != types.RANK_1 {
		moves[idx] = m
		return 1
	}

	// Promotion
	for _, pt := range []types.PieceType{types.KNIGHT, types.BISHOP, types.ROOK, types.QUEEN} {
		// Copy the move, since promotion can only be set once.
		pm := m
		pm.SetMoveType(move.PROMOTION)
		pm.SetPromitionPieceType(pt)
		moves[idx] = pm
		idx++
	}

	return idx
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
