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

type Search struct {
	ctx              context.Context
	Pos              position.Position
	nodes            uint64
	PV               pvline.PVLine
	KillerMoves      [1024][2]move.Move
	searchHistory    [1024]uint64
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
	score, err := s.negamax(&pos, alpha, beta, depth, 0, &pvl, true)
	if err != nil {
		return Info{}, err
	}
	return Info{
		Depth: depth,
		Score: score,
		PV:    *pvl.Copy(),
	}, nil
}

func (s *Search) negamax(pos *position.Position, alpha, beta int16, maxDepth, ply uint8, pvl *pvline.PVLine, canNull bool) (int16, error) {
	// value to info channel and check if we are done
	select {
	case <-s.ctx.Done():
		return 0, s.ctx.Err()
	default:
	}

	isRoot := ply == 0
	depth := maxDepth - ply
	mateValue := -evaluation.INF + int16(ply)
	pvNode := beta-alpha != 1

	// Increase Depth, if in Check.
	// This also means that we do not enter quiescence, if in check.
	isInCheck := pos.IsInCheck(pos.SideToMove)
	if isInCheck {
		maxDepth++
	}

	// Evaluate the leaf node
	if ply == maxDepth {
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

	// Check if we can use Futility Pruning
	fPrune := !pvNode &&
		depth < futility_pruning_depth &&
		!isInCheck &&
		!evaluation.IsCheckmateValue(alpha) &&
		!evaluation.IsCheckmateValue(beta) &&
		evaluation.Evaluation(pos)+futility_pruning_margin[depth] <= alpha

	potentialPVLine := pvline.PVLine{}
	var bestMove move.Move
	var bestScore int16 = -evaluation.INF
	var legalMoves uint8
	nodeType := transpositiontable.AlphaNode
	var state position.State

	// Generate all moves and order them
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalMoves(moves)
	s.orderMoves(pos, moves, pvMove, ttMove, ply)

	for i := uint8(0); i < moves.Length(); i++ {
		m := moves.Get(i)

		isCapture := pos.IsCapture(*m)

		pos.MakeMove(*m, &state)
		if !pos.IsLegal() {
			pos.UnMakeMove(&state)
			continue
		}
		legalMoves++

		// Fulility Pruning
		if fPrune && legalMoves > 0 && !isCapture && !pos.IsInCheck(pos.SideToMove) {
			pos.UnMakeMove(&state)
			continue
		}

		score, err := s.PrincipalVariationSearch(pos, alpha, beta, maxDepth, ply, &potentialPVLine, canNull, nodeType == transpositiontable.PVNode)
		if err != nil {
			return 0, err
		}
		pos.UnMakeMove(&state)

		if score > bestScore {
			bestScore = score
			bestMove = *m
		}

		if score >= beta {
			nodeType = transpositiontable.BetaNode
			// Update Killer Move, if quiet move
			if pos.GetPiece(m.GetTargetSquare()) == types.NO_PIECE {
				if s.KillerMoves[ply][0] != bestMove {
					s.KillerMoves[ply][1] = s.KillerMoves[ply][0]
				}
				s.KillerMoves[ply][0] = bestMove
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

	// var prev position.Position
	var state position.State

	// Generate all captures and order them
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalCaptures(moves)
	s.orderMoves(pos, moves, move.NullMove, move.NullMove, ply)
	for i := uint8(0); i < moves.Length(); i++ {
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

		pos.MakeMove(*m, &state)
		if !pos.IsLegal() {
			pos.UnMakeMove(&state)
			continue
		}
		score, err := s.quiescence(pos, -beta, -alpha, ply+1)
		pos.UnMakeMove(&state)
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

	var t, inc int
	if s.Pos.SideToMove == types.BLACK {
		t = sp.BTime
		inc = sp.BInc
	} else {
		t = sp.WTime
		inc = sp.WInc
	}

	var movetime int
	if sp.MoveTime > 0 {
		movetime = sp.MoveTime
	} else if t > 0 && sp.MovesToGo > 0 {
		// calculate reasonable time, there is possibly a better way
		movetime = (t + inc*sp.MovesToGo) / sp.MovesToGo
		// do not calculate too long
		movetime = min(movetime, maxTimeInMs)
	} else {
		movetime = maxTimeInMs
	}

	// Current puffer of 10% for teardown process, but at least 50milliseconds
	movetime = movetime - max(movetime/10, 50)

	fmt.Printf("info string calculated timeout %v\n", movetime)
	return context.WithTimeout(ctx, time.Duration(movetime)*time.Millisecond)
}
