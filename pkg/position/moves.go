package position

import (
	"errors"
	"math"

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

// Static Values for MVV-LVA Ordering
// See https://www.chessprogramming.org/MVV-LVA
var MVV_LVA_SCORES [types.PIECE_TYPE_NUMBER - 1][types.PIECE_TYPE_NUMBER]uint16

func init() {
	// Init the MVV-LVA Values
	// For the values to be disjunct, the victim is multiplied by 10
	// To make a difference for PAWNs (value 0) victim a increased by 1
	for victim := types.QUEEN; victim >= types.PAWN; victim-- {
		for aggressor := types.PAWN; aggressor < types.PIECE_TYPE_NUMBER; aggressor++ {
			MVV_LVA_SCORES[victim][aggressor] = uint16(10*(victim+1) - (aggressor))
		}
	}
}

func (pos *Position) sortMoves(moves *move.MoveList, bestGuess move.Move) {
	for idx := uint8(0); idx < moves.Length(); idx++ {
		m := moves.Get(idx)
		if *m == bestGuess {
			m.SetScore(math.MaxUint16)
			continue
		}

		// Score captures
		target := pos.GetPiece(m.GetTargetSquare())
		// No capture no score
		if target == types.NO_PIECE {
			m.SetScore(0)
			continue
		}
		source := pos.GetPiece(m.GetSourceSquare())
		m.SetScore(MVV_LVA_SCORES[target.Type()][source.Type()])
	}

	moves.Sort()
}

// GeneratePseudoLegalCapturesOrdered generates all pseudo legal moves and directly order them
func (pos *Position) GeneratePseudoLegalCapturesOrdered(moves *move.MoveList, bestGuess move.Move) {
	pos.GeneratePseudoLegalCaptures(moves)
	pos.sortMoves(moves, bestGuess)
}

func (pos *Position) GeneratePseudoLegalCaptures(moves *move.MoveList) {
	occupied := pos.AllPieces
	destinations := pos.AllPiecesByColor[types.SwitchColor(pos.SideToMove)]
	// Sliding Pieces
	generateMovesHelper(
		moves,
		pos.PiecesBitboard[pos.SideToMove][types.ROOK],
		occupied,
		destinations,
		rook.AttacksBySquare,
	)
	generateMovesHelper(
		moves,
		pos.PiecesBitboard[pos.SideToMove][types.BISHOP],
		occupied,
		destinations,
		bishop.AttacksBySquare,
	)
	generateMovesHelper(
		moves,
		pos.PiecesBitboard[pos.SideToMove][types.QUEEN],
		occupied,
		destinations,
		queen.AttacksBySquare,
	)

	// Pieces ignoring occupation
	generateMovesHelper(
		moves,
		pos.PiecesBitboard[pos.SideToMove][types.KNIGHT],
		bitboard.Empty,
		destinations,
		func(square int, _ bitboard.Bitboard) bitboard.Bitboard {
			return knight.AttacksBySquare(square)
		},
	)
	generateMovesHelper(
		moves,
		pos.PiecesBitboard[pos.SideToMove][types.KING],
		bitboard.Empty,
		destinations,
		func(square int, _ bitboard.Bitboard) bitboard.Bitboard {
			return king.AttacksBySquare(square)
		},
	)

	// Pawns
	sourceSquareIt := bitboard.SquareIndexSerializationIterator(pos.PiecesBitboard[pos.SideToMove][types.PAWN])
	for {
		sourceSquare := sourceSquareIt()
		if sourceSquare == types.SQUARE_NONE {
			break
		}

		pawnAttacks := pawn.AttacksBySquare(pos.SideToMove, sourceSquare)

		// Attacks
		targetSquareIt := bitboard.SquareIndexSerializationIterator(pawnAttacks & pos.AllPiecesByColor[types.SwitchColor(pos.SideToMove)])
		for {
			targetSquare := targetSquareIt()
			if targetSquare == types.SQUARE_NONE {
				break
			}
			var m move.Move
			m.SetSourceSquare(sourceSquare)
			m.SetTargetSquare(targetSquare)
			// Pawn Moves with optional Promotion
			pawnMoveWithPromotion(moves, pos.SideToMove, sourceSquare, targetSquare)
		}

		// En Passant
		if pos.EnPassant != types.SQUARE_NONE {
			targetSquareIt := bitboard.SquareIndexSerializationIterator(pawnAttacks & bitboard.BitBySquares(pos.EnPassant))
			for {
				targetSquare := targetSquareIt()
				if targetSquare == types.SQUARE_NONE {
					break
				}
				var m move.Move
				m.SetSourceSquare(sourceSquare)
				m.SetTargetSquare(targetSquare)
				m.SetMoveType(move.EN_PASSANT)
				moves.Append(m)
			}
		}
	}
}

// GeneratePseudoLegalMovesOrdered generates all pseudo legal moves and directly order them
func (pos *Position) GeneratePseudoLegalMovesOrdered(moves *move.MoveList, bestGuess move.Move) {
	pos.GeneratePseudoLegalMoves(moves)
	pos.sortMoves(moves, bestGuess)
}

// GeneratePseudoLegalMoves generates all pseudo legal moves
func (pos *Position) GeneratePseudoLegalMoves(moves *move.MoveList) {
	occupied := pos.AllPieces
	destinations := ^pos.AllPiecesByColor[pos.SideToMove]

	// Sliding Pieces
	generateMovesHelper(
		moves,
		pos.PiecesBitboard[pos.SideToMove][types.ROOK],
		occupied,
		destinations,
		rook.AttacksBySquare,
	)
	generateMovesHelper(
		moves,
		pos.PiecesBitboard[pos.SideToMove][types.BISHOP],
		occupied,
		destinations,
		bishop.AttacksBySquare,
	)

	generateMovesHelper(
		moves,
		pos.PiecesBitboard[pos.SideToMove][types.QUEEN],
		occupied,
		destinations,
		queen.AttacksBySquare,
	)

	// Pieces ignoring occupation
	generateMovesHelper(
		moves,
		pos.PiecesBitboard[pos.SideToMove][types.KNIGHT],
		bitboard.Empty,
		destinations,
		func(square int, _ bitboard.Bitboard) bitboard.Bitboard {
			return knight.AttacksBySquare(square)
		},
	)
	generateMovesHelper(
		moves,
		pos.PiecesBitboard[pos.SideToMove][types.KING],
		bitboard.Empty,
		destinations,
		func(square int, _ bitboard.Bitboard) bitboard.Bitboard {
			return king.AttacksBySquare(square)
		},
	)

	// Pawns
	sourceSquareIt := bitboard.SquareIndexSerializationIterator(pos.PiecesBitboard[pos.SideToMove][types.PAWN])
	for {
		sourceSquare := sourceSquareIt()
		if sourceSquare == types.SQUARE_NONE {
			break
		}

		pawnAttacks := pawn.AttacksBySquare(pos.SideToMove, sourceSquare)

		// Pushes
		targetSquareIt := bitboard.SquareIndexSerializationIterator(pawn.PushesBySquare(pos.SideToMove, sourceSquare, occupied))
		for {
			targetSquare := targetSquareIt()
			if targetSquare == types.SQUARE_NONE {
				break
			}
			var m move.Move
			m.SetSourceSquare(sourceSquare)
			m.SetTargetSquare(targetSquare)
			// Pawn Moves with optional Promotion
			pawnMoveWithPromotion(moves, pos.SideToMove, sourceSquare, targetSquare)
		}

		// Attacks
		targetSquareIt = bitboard.SquareIndexSerializationIterator(pawnAttacks & pos.AllPiecesByColor[types.SwitchColor(pos.SideToMove)])
		for {
			targetSquare := targetSquareIt()
			if targetSquare == types.SQUARE_NONE {
				break
			}
			var m move.Move
			m.SetSourceSquare(sourceSquare)
			m.SetTargetSquare(targetSquare)
			// Pawn Moves with optional Promotion
			pawnMoveWithPromotion(moves, pos.SideToMove, sourceSquare, targetSquare)
		}

		// En Passant
		if pos.EnPassant != types.SQUARE_NONE {
			targetSquareIt := bitboard.SquareIndexSerializationIterator(pawnAttacks & bitboard.BitBySquares(pos.EnPassant))
			for {
				targetSquare := targetSquareIt()
				if targetSquare == types.SQUARE_NONE {
					break
				}
				var m move.Move
				m.SetSourceSquare(sourceSquare)
				m.SetTargetSquare(targetSquare)
				m.SetMoveType(move.EN_PASSANT)
				moves.Append(m)
			}
		}
	}

	// Castling
	for _, c := range []Castling{WHITE_CASTLING_KING, WHITE_CASTLING_QUEEN, BLACK_CASTLING_KING, BLACK_CASTLING_QUEEN} {
		if c.Color() != pos.SideToMove {
			continue
		}
		if !pos.CanCastleNow(c) {
			continue
		}
		var m move.Move
		m.SetMoveType(move.CASTLING)
		sourceSquare := bitboard.LeastSignificantOneBit(pos.PiecesBitboard[pos.SideToMove][types.KING])
		var targetSquare int
		switch c.Side() {
		case CASTLING_KING:
			targetSquare = sourceSquare + 2
		case CASTLING_QUEEN:
			targetSquare = sourceSquare - 2
		}
		m.SetSourceSquare(sourceSquare)
		m.SetTargetSquare(targetSquare)
		moves.Append(m)
	}
}

func (pos *Position) MakeMove(m move.Move) {
	resetHalfmoveClock := false

	if pos.EnPassant != types.SQUARE_NONE {
		pos.zobristUpdateEnPassant(pos.EnPassant)
		pos.EnPassant = types.SQUARE_NONE

	}

	sourceSquare := m.GetSourceSquare()
	targetSquare := m.GetTargetSquare()

	targetPiece := pos.GetPiece(targetSquare)
	if targetPiece != types.NO_PIECE {
		pos.DeletePiece(targetSquare)
		resetHalfmoveClock = true
	}

	for _, s := range []int{sourceSquare, targetSquare} {
		switch s {
		case types.SQUARE_A1:
			pos.Castling = pos.Castling &^ WHITE_CASTLING_QUEEN
			pos.zobristUpdateCastling(WHITE_CASTLING_QUEEN)
		case types.SQUARE_H1:
			pos.Castling = pos.Castling &^ WHITE_CASTLING_KING
			pos.zobristUpdateCastling(WHITE_CASTLING_KING)
		case types.SQUARE_A8:
			pos.Castling = pos.Castling &^ BLACK_CASTLING_QUEEN
			pos.zobristUpdateCastling(BLACK_CASTLING_KING)
		case types.SQUARE_H8:
			pos.Castling = pos.Castling &^ BLACK_CASTLING_KING
			pos.zobristUpdateCastling(BLACK_CASTLING_KING)
		case types.SQUARE_E1:
			pos.Castling = pos.Castling &^ (WHITE_CASTLING_QUEEN | WHITE_CASTLING_KING)
			pos.zobristUpdateCastling(WHITE_CASTLING_QUEEN)
			pos.zobristUpdateCastling(WHITE_CASTLING_KING)
		case types.SQUARE_E8:
			pos.Castling = pos.Castling &^ (BLACK_CASTLING_QUEEN | BLACK_CASTLING_KING)
			pos.zobristUpdateCastling(BLACK_CASTLING_QUEEN)
			pos.zobristUpdateCastling(BLACK_CASTLING_KING)
		}
	}

	piece := pos.MovePiece(sourceSquare, targetSquare)
	switch piece.Type() {
	// Set en passant
	case types.PAWN:
		resetHalfmoveClock = true
		if abs(sourceSquare-targetSquare) == 2*types.FILE_NUMBER {
			pos.EnPassant = targetSquare
			pos.zobristUpdateEnPassant(pos.EnPassant)
			switch pos.SideToMove {
			case types.BLACK:
				pos.EnPassant += types.FILE_NUMBER
			case types.WHITE:
				pos.EnPassant -= types.FILE_NUMBER
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
		switch pos.SideToMove {
		case types.WHITE:
			pawnToRemoveSquare = targetSquare - types.FILE_NUMBER
		case types.BLACK:
			pawnToRemoveSquare = targetSquare + types.FILE_NUMBER
		}
		pos.DeletePiece(pawnToRemoveSquare)
	case move.PROMOTION:
		// Promote piece
		pos.DeletePiece(targetSquare)
		pos.SetPiece(types.NewPiece(pos.SideToMove, m.GetPromitionPieceType()), targetSquare)
	}

	// Update Side to Move
	pos.Ply++
	pos.SideToMove = types.SwitchColor(pos.SideToMove)
	pos.zobristUpdateColor()

	if resetHalfmoveClock {
		pos.HalfMoveClock = 0
	} else {
		pos.HalfMoveClock++
	}

	// Generate Helper Bitboards
	pos.generateHelperBitboards()
}

func (pos *Position) MakeNullMove() int {
	pos.Ply++
	ep := pos.EnPassant
	if pos.EnPassant != types.SQUARE_NONE {
		pos.zobristUpdateEnPassant(pos.EnPassant)
		pos.EnPassant = types.SQUARE_NONE

	}

	// Update Side to Move
	pos.SideToMove = types.SwitchColor(pos.SideToMove)
	pos.zobristUpdateColor()
	return ep
}

func (pos *Position) UnMakeNullMove(enPassantSquare int) {
	pos.Ply--

	if enPassantSquare != types.SQUARE_NONE {
		pos.EnPassant = enPassantSquare
		pos.zobristUpdateEnPassant(pos.EnPassant)
	}

	// Update Side to Move
	pos.SideToMove = types.SwitchColor(pos.SideToMove)
	pos.zobristUpdateColor()
}

// generateMovesHelper generates a list of moves from a given list of paramters
func generateMovesHelper(moves *move.MoveList, sources, occupied, destinations bitboard.Bitboard, attacks func(square int, occupied bitboard.Bitboard) bitboard.Bitboard) {
	var targetSquareIt func() int
	var sourceSquare, targetSquare int
	sourceSquareIt := bitboard.SquareIndexSerializationIterator(sources)
	for {
		sourceSquare = sourceSquareIt()
		if sourceSquare == types.SQUARE_NONE {
			break
		}
		targetSquareIt = bitboard.SquareIndexSerializationIterator(attacks(sourceSquare, occupied) & destinations)
		for {
			targetSquare = targetSquareIt()
			if targetSquare == types.SQUARE_NONE {
				break
			}
			var m move.Move
			m.SetSourceSquare(sourceSquare)
			m.SetTargetSquare(targetSquare)
			moves.Append(m)
		}
	}
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
	m.SetTargetSquare(destinationSquare)

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

func pawnMoveWithPromotion(moves *move.MoveList, sideToMove types.Color, sourceSquare, targetSquare int) {
	var m move.Move
	m.SetSourceSquare(sourceSquare)
	m.SetTargetSquare(targetSquare)
	// No promotion
	if sideToMove == types.WHITE && types.RankOfSquare(targetSquare) != types.RANK_8 {
		moves.Append(m)
		return
	}
	if sideToMove == types.BLACK && types.RankOfSquare(targetSquare) != types.RANK_1 {
		moves.Append(m)
		return
	}

	// Promotion
	for _, pt := range []types.PieceType{types.KNIGHT, types.BISHOP, types.ROOK, types.QUEEN} {
		// Copy the move, since promotion can only be set once.
		pm := m
		pm.SetMoveType(move.PROMOTION)
		pm.SetPromitionPieceType(pt)
		moves.Append(pm)
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
