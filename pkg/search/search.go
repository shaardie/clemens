package search

import (
	"context"
	"fmt"
	"time"

	"github.com/shaardie/clemens/pkg/evaluation"
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search/pvline"
	"github.com/shaardie/clemens/pkg/search/transpositiontable"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	widen_window               = 50
	max_depth            uint8 = 100
	quiescence_max_depth uint8 = 100
	maxTimeInMs                = 1000000
)

const futility_pruning_depth uint8 = 5

var futility_pruning_margin = [futility_pruning_depth]int16{0, 100, 150, 200, 250}

const staticNullMovePruningMarging int16 = 75

type Search struct {
	ctx              context.Context
	Pos              position.Position
	nodes            uint64
	PV               pvline.PVLine
	KillerMoves      [1024][2]move.Move
	searchHistory    [1024]uint64
	history          [types.COLOR_NUMBER][types.SQUARE_NUMBER][types.SQUARE_NUMBER]uint16
	counter          [types.COLOR_NUMBER][types.SQUARE_NUMBER][types.SQUARE_NUMBER]move.Move
	searchHistoryPly int
}

type SearchParameter struct {
	WTime     int
	BTime     int
	WInc      int
	BInc      int
	MovesToGo int
	Depth     uint8
	MoveTime  int
	Infinite  bool
}

type Info struct {
	Depth uint8
	PV    pvline.PVLine
	Score int16
}

func (s *Search) bestMove() move.Move {
	return s.PV.GetBestMove()
}

func NewSearch(pos position.Position) *Search {
	s := &Search{
		Pos: pos,
	}
	return s
}

func (s *Search) Search(ctx context.Context, sp SearchParameter) move.Move {
	ctx, cancel := s.contextFromSearchParameter(ctx, sp)
	depth := max_depth
	if sp.Depth > 0 {
		depth = sp.Depth
	}
	s.ctx = ctx
	s.SearchIterative(depth)
	cancel()

	// We need at least a valid move
	if s.bestMove() == move.NullMove {
		s.ctx = context.TODO()
		s.SearchIterative(1)
	}
	return s.bestMove()
}

func (s *Search) SearchIterative(maxDepth uint8) {
	start := time.Now()
	alpha := -evaluation.INF
	beta := evaluation.INF
	var depth uint8 = 1
	for depth <= maxDepth {
		i, err := s.SearchRoot(depth, alpha, beta)
		// Timeout
		if err != nil {
			return
		}

		// Reduce the search space by using an aspiration window
		// See https://www.chessprogramming.org/Aspiration_Windows
		// If the score is not in the last windows,
		// re-run the search with the wider window, do not use the result and do not increase the depth.
		if i.Score <= alpha || i.Score >= beta {
			fmt.Printf("info string windows [%v,%v] too small for value %v. Re-run search.\n", alpha, beta, i.Score)
			alpha = -evaluation.INF
			beta = evaluation.INF
			continue
		}

		s.PV = *i.PV.Copy()
		alpha = i.Score - widen_window
		beta = i.Score + widen_window
		depth++

		// Print info
		t := max(time.Since(start).Milliseconds(), 1) // should never be zero
		fmt.Printf(
			"info depth %v score cp %v time %v nodes %v nps %v hashfull %v pv %v\n",
			i.Depth,
			i.Score,
			t,
			s.nodes,
			int64(s.nodes)*1000/t,
			transpositiontable.HashFull(),
			i.PV,
		)
	}
}

func (s *Search) SearchRoot(depth uint8, alpha, beta int16) (Info, error) {

	s.KillerMoves = [1024][2]move.Move{}
	pos := s.Pos
	pvl := pvline.PVLine{}
	score, err := s.negamax(&pos, alpha, beta, depth, 0, &pvl, true, move.NullMove)
	if err != nil {
		return Info{}, err
	}
	return Info{
		Depth: depth,
		Score: score,
		PV:    *pvl.Copy(),
	}, nil
}

func (s *Search) negamax(pos *position.Position, alpha, beta int16, depth, ply uint8, pvl *pvline.PVLine, canNull bool, previousMove move.Move) (int16, error) {
	// value to info channel and check if we are done
	select {
	case <-s.ctx.Done():
		return 0, s.ctx.Err()
	default:
	}

	isRoot := ply == 0
	mateValue := -evaluation.INF + int16(ply)
	pvNode := beta-alpha != 1

	// Increase Depth, if in Check.
	// This also means that we do not enter quiescence, if in check.
	isInCheck := pos.IsInCheck(pos.SideToMove)
	if isInCheck {
		depth++
	}

	// Evaluate the leaf node
	if depth <= 0 {
		return s.quiescence(pos, alpha, beta, ply)
	}
	s.nodes++

	// Check if the position is a repetition.
	// On the first repetitions we return our contempt value.
	// https://www.chessprogramming.org/Repetitions
	if !isRoot && !isInCheck && s.isRepetition(pos) {
		return evaluation.Contempt(pos), nil
	}
	s.pushHistory(pos)
	defer s.popHistory()

	pvMove := s.PV.GetBestMoveByPly(ply)

	// Check if we can use the transition table but not on root
	score, use, ttMove := transpositiontable.Get(pos.ZobristHash, alpha, beta, depth, ply)
	if !isRoot && !pvNode && use {
		return score, nil
	}

	// Static Null Move Pruning
	if !isInCheck && !pvNode && !evaluation.IsCheckmateValue(beta) {
		// score - margin as potential new beta
		b := evaluation.Evaluation(pos) - staticNullMovePruningMarging*int16(depth)
		if b >= beta {
			return b, nil
		}
	}

	potentialPVLine := pvline.PVLine{}

	// Null Move Pruning
	// https://www.chessprogramming.org/Null_Move_Pruning
	if depth > 2 && canNull && !isInCheck && !pvNode && !evaluation.IsPawnEndgame(pos) && evaluation.Evaluation(pos) > beta {
		ep := pos.MakeNullMove()
		var R uint8 = 2
		if depth > 6 {
			R = 3
		}
		score, err := s.negamax(pos, -beta, -beta+1, depth-R-1, ply+1, &potentialPVLine, false, move.NullMove)
		pos.UnMakeNullMove(ep)
		potentialPVLine.Reset()
		if err != nil {
			return 0, err
		}
		score = -score
		if score >= beta {
			return beta, nil
		}
	}

	// Check if we can use Futility Pruning
	fPrune := !pvNode &&
		depth < futility_pruning_depth &&
		!isInCheck &&
		!evaluation.IsCheckmateValue(alpha) &&
		!evaluation.IsCheckmateValue(beta) &&
		evaluation.Evaluation(pos)+futility_pruning_margin[depth] <= alpha

	var prevPos position.Position
	var bestMove move.Move
	var bestScore int16 = -evaluation.INF
	var legalMoves uint8
	var err error
	nodeType := transpositiontable.AlphaNode

	// Generate all moves and order them
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalMoves(moves)
	s.scoreMoves(pos, moves, pvMove, ttMove, ply)
	// s.orderMoves(pos, moves, pvMove, ttMove, ply)

	for i := range moves.Length() {
		moves.SortIndex(i)
		m := moves.Get(i)
		prevPos = *pos
		pos.MakeMove(*m)
		if !pos.IsLegal() {
			*pos = prevPos
			continue
		}
		legalMoves++

		// Fulility Pruning
		if fPrune && !prevPos.IsCapture(*m) && m.GetMoveType() != move.PROMOTION && !pos.IsInCheck(pos.SideToMove) {
			*pos = prevPos
			continue
		}

		// Principal Variation Search
		if nodeType != transpositiontable.PVNode {
			// Alpha was not updated, we do not have a better move yet, so we search full
			score, err = s.negamax(pos, -beta, -alpha, depth-1, ply+1, &potentialPVLine, true, previousMove)
			if err != nil {
				return 0, err
			}
			score = -score
		} else {
			// Alpha was already updated, so we found a better move already. So we are only doint a null windows search
			// Although, if this raises alpha, we have do to full search.
			score, err = s.negamax(pos, -alpha-1, -alpha, depth-1, ply+1, &pvline.PVLine{}, true, previousMove)
			if err != nil {
				return 0, err
			}
			score = -score
			// Rerun search
			if score > alpha {
				score, err = s.negamax(pos, -beta, -alpha, depth-1, ply+1, &potentialPVLine, true, previousMove)
				if err != nil {
					return 0, err
				}
				score = -score
			}
		}

		*pos = prevPos

		if score > bestScore {
			bestScore = score
			bestMove = *m
		}

		if score >= beta {
			nodeType = transpositiontable.BetaNode
			if pos.GetPiece(m.GetTargetSquare()) == types.NO_PIECE && m.GetMoveType() != move.EN_PASSANT {
				// Update Killer Move, if quiet move
				if s.KillerMoves[ply][0] != bestMove {
					s.KillerMoves[ply][1] = s.KillerMoves[ply][0]
				}
				s.KillerMoves[ply][0] = bestMove

				// Remember move for history heuristic
				sourceSquare := m.GetSourceSquare()
				targetSquare := m.GetTargetSquare()
				s.history[pos.SideToMove][sourceSquare][targetSquare] += uint16(depth) * uint16(depth)
				if s.history[pos.SideToMove][m.GetSourceSquare()][m.GetTargetSquare()] > killerMoveScore-2 {
					for i := range types.SQUARE_NUMBER {
						for j := range types.SQUARE_NUMBER {
							s.history[pos.SideToMove][i][j] /= 2
						}
					}
				}

				// Update counter moves
				if previousMove != move.NullMove {
					s.counter[pos.SideToMove][previousMove.GetSourceSquare()][previousMove.GetTargetSquare()] = *m
				}
			}
			break
		}

		if score > alpha {
			nodeType = transpositiontable.PVNode
			alpha = score

			pvl.Update(bestMove, &potentialPVLine)

		}
		potentialPVLine.Reset()
	}

	// There are no legal moves, so it is either a checkmate or a stalemate
	if legalMoves == 0 {
		// Checkmate, set lowest possible value, but increase by the number of plys,
		// so the engine is looking for shorter mates.
		if isInCheck {
			return mateValue, nil
		}
		// stalemate
		return evaluation.Contempt(pos), nil
	}

	select {
	case <-s.ctx.Done():
		return 0, s.ctx.Err()
	default:
		transpositiontable.PotentiallySave(pos.ZobristHash, bestMove, depth, bestScore, nodeType, s.Pos.HalfMoveClock)
	}
	return bestScore, nil
}

func (s *Search) quiescence(pos *position.Position, alpha, beta int16, ply uint8) (int16, error) {
	s.nodes++
	// value to info channel and check if we are done
	select {
	case <-s.ctx.Done():
		return 0, s.ctx.Err()
	default:
	}

	stand_pat := evaluation.Evaluation(pos)
	if stand_pat >= beta {
		return beta, nil
	}
	if alpha < stand_pat {
		alpha = stand_pat
	}

	// Hard limit
	if ply == quiescence_max_depth {
		return alpha, nil
	}

	var prevPos position.Position

	// Generate all captures and order them
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalCaptures(moves)
	// s.orderMoves(pos, moves, move.NullMove, move.NullMove, ply)
	s.scoreMoves(pos, moves, move.NullMove, move.NullMove, ply)
	for i := range moves.Length() {
		moves.SortIndex(i)
		m := moves.Get(i)

		// Delta Pruning, https://www.chessprogramming.org/Delta_Pruning
		// If the current capture plus some safety margin is not able to raise alpha, we can skip the move.
		if m.GetMoveType() != move.EN_PASSANT {
			// Safety margin of 2 centipawns
			margin := 2 * evaluation.PieceValue[types.PAWN]
			// Take promotion into account.
			if m.GetMoveType() == move.PROMOTION {
				margin = margin - evaluation.PieceValue[types.PAWN] + evaluation.PieceValue[m.GetPromitionPieceType()]
			}
			// Skip Delta Pruning in the endgame to not become blind against insufficient material.
			if stand_pat+evaluation.PieceValue[pos.PiecesBoard[m.GetTargetSquare()].Type()]+margin < alpha && !evaluation.IsEndgame(pos) {
				continue
			}
		}

		// If the static exchange of pieces on the target square does not gain any positive material value,
		// we can ignore this move completely.
		// En Passants are excluded because the target square of the pawn is not the square of the capture.
		if m.GetMoveType() != move.EN_PASSANT && evaluation.StaticExchangeEvaluation(pos, m) < 0 {
			continue
		}

		prevPos = *pos
		pos.MakeMove(*m)
		if !pos.IsLegal() {
			*pos = prevPos
			continue
		}
		score, err := s.quiescence(pos, -beta, -alpha, ply+1)
		*pos = prevPos
		if err != nil {
			return 0, err
		}
		score = -score
		if score >= beta {
			return beta, nil
		}

		if score > alpha {
			alpha = score
		}
	}

	return alpha, nil
}

func (s *Search) contextFromSearchParameter(ctx context.Context, sp SearchParameter) (context.Context, context.CancelFunc) {
	// No need for any timeout
	if sp.Infinite {
		return ctx, func() {}
	}
	movetime := calculateTime(s.Pos.SideToMove, s.searchHistoryPly, sp)
	fmt.Printf("info string calculated timeout %v\n", movetime)
	return context.WithTimeout(ctx, time.Duration(movetime)*time.Millisecond)
}

func calculateTime(sideToMove types.Color, plys int, sp SearchParameter) int {
	var t, inc int
	if sideToMove == types.BLACK {
		t = sp.BTime
		inc = sp.BInc
	} else {
		t = sp.WTime
		inc = sp.WInc
	}

	// calculate for 60 moves or at least 20 remaining
	var remainingMoves = max(60-plys/2, 20)

	var movetime int
	if sp.MoveTime > 0 {
		movetime = sp.MoveTime
	} else if t > 0 {
		// calculate reasonable time, there is possibly a better way
		movetime = (t + inc*remainingMoves) / remainingMoves
		// do not calculate too long
		movetime = min(movetime, maxTimeInMs)
	} else { // I do not know, calcuate a second
		movetime = 1000
	}

	// Current puffer of 10% for teardown process, but at least 50milliseconds
	return movetime - max(movetime/10, 50)
}
