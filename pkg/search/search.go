package search

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/shaardie/clemens/pkg/evaluation"
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search/pvline"
	"github.com/shaardie/clemens/pkg/search/transpositiontable"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	inf                        = 100000
	widen_window               = 50
	max_depth            uint8 = math.MaxUint8
	quiescence_max_depth uint8 = 100
	maxTimeInMs                = 10000
)

type Search struct {
	Pos                    position.Position
	nodes                  uint64
	betaCutOffs            uint64
	alphaCutOffs           uint64
	quiescenceNodes        uint64
	transpositiontableHits uint64
	PV                     pvline.PVLine
	KillerMoves            [1024][2]move.Move
	SearchHistory          [1024]uint64
	SearchHistoryPly       int
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
	Time  int64
	PV    pvline.PVLine
	Score int
}

func (s *Search) bestMove() move.Move {
	return s.PV.GetBestMove()
}

func NewSearch(pos position.Position) *Search {
	s := &Search{
		Pos: pos,
	}
	s.SearchHistory[s.SearchHistoryPly] = pos.ZobristHash
	return s
}

func (s *Search) MakeMoveFromString(m string) error {
	err := s.Pos.MakeMoveFromString(m)
	if err != nil {
		return err
	}
	s.SearchHistoryPly++
	s.SearchHistory[s.SearchHistoryPly] = s.Pos.ZobristHash
	return nil
}

func (s *Search) Search(ctx context.Context, sp SearchParameter) move.Move {
	ctx, cancel := s.contextFromSearchParameter(ctx, sp)
	depth := max_depth
	if sp.Depth > 0 {
		depth = sp.Depth
	}
	s.SearchIterative(ctx, depth)
	cancel()

	// We need at least a valid move
	if s.bestMove() == move.NullMove {
		s.SearchIterative(context.TODO(), 1)
	}
	return s.bestMove()
}

func (s *Search) SearchIterative(ctx context.Context, maxDepth uint8) {
	alpha := -inf
	beta := inf
	var depth uint8 = 1
	goodGuess := move.NullMove
	for depth <= maxDepth {
		i, err := s.SearchRoot(ctx, depth, alpha, beta, goodGuess)
		// Timeout
		if err != nil {
			return
		}

		// Reduce the search space by using an aspiration window
		// See https://www.chessprogramming.org/Aspiration_Windows
		// If the score is not in the last windows,
		// re-run the search with the wider window, do not use the result and do not increase the depth.
		if i.Score <= alpha || i.Score >= beta {
			fmt.Printf("info string windows [%v,%v] too small. Re-run search.\n", alpha, beta)
			alpha = -inf
			beta = inf
			continue
		}

		s.PV = *i.PV.Copy()
		goodGuess = s.PV.GetBestMove()
		alpha = i.Score - widen_window
		beta = i.Score + widen_window
		depth++

		// Print info
		fmt.Printf("info depth %v score cp %v nodes %v time %v pv %v\n", i.Depth, i.Score, s.nodes, i.Time, i.PV)
		fmt.Printf("info string beta-cutoffs %v alpha-cutoffs %v quiescence-nodes %v transpositiontable-hits %v\n", s.betaCutOffs, s.alphaCutOffs, s.quiescenceNodes, s.transpositiontableHits)
	}
}

func (s *Search) SearchRoot(ctx context.Context, depth uint8, alpha, beta int, goodGuess move.Move) (Info, error) {
	start := time.Now()
	pos := s.Pos
	pvl := pvline.PVLine{}
	score, err := s.negamax(ctx, &pos, alpha, beta, depth, 0, true, &pvl, true, true)
	if err != nil {
		return Info{}, err
	}
	return Info{
		Depth: depth,
		Time:  time.Since(start).Milliseconds(),
		Score: score,
		PV:    *pvl.Copy(),
	}, nil
}

func (s *Search) negamax(ctx context.Context, pos *position.Position, alpha, beta int, maxDepth, ply uint8, pvNode bool, pvl *pvline.PVLine, isRoot bool, canNull bool) (int, error) {

	// value to info channel and check if we are done
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	depth := maxDepth - ply
	mateValue := -inf + int(ply)

	// Mate Distance Pruning
	// https://www.chessprogramming.org/Mate_Distance_Pruning
	if !isRoot {
		if alpha < mateValue {
			alpha = mateValue
		}
		if beta > -mateValue+1 {
			beta = -mateValue + 1
		}
		if alpha >= beta {
			return alpha, nil
		}
	}

	// Increase Depth, if in Check.
	// This also means that we do not enter quiescence, if in check.
	isInCheck := pos.IsInCheck(pos.SideToMove)
	if isInCheck {
		maxDepth++
	}

	// Evaluate the leaf node
	if ply == maxDepth {
		return s.quiescence(ctx, pos, alpha, beta, ply)
	}
	s.nodes++

	if !isRoot && s.isRepetition(pos.ZobristHash) {
		if evaluation.IsEndgame(pos) {
			return 0, nil
		}
		return -100, nil
	}
	s.SearchHistoryPly++
	s.SearchHistory[s.SearchHistoryPly] = pos.ZobristHash
	defer func() {
		s.SearchHistoryPly--
		s.SearchHistory[s.SearchHistoryPly] = 0
	}()

	pvMove := s.PV.GetBestMoveByPly(ply)
	ttMove := move.NullMove

	// Check if we can use the transition table but not on root
	if !isRoot {
		te, found, isGoodGuess := transpositiontable.TTable.Get(pos.ZobristHash, depth)
		if found {
			s.transpositiontableHits++
			switch te.NodeType {
			case transpositiontable.AlphaNode:
				if te.Score <= alpha {
					s.alphaCutOffs++
					return alpha, nil
				}
			case transpositiontable.BetaNode:
				if te.Score > beta {
					s.betaCutOffs++
					return beta, nil
				}
			case transpositiontable.PVNode:
				// In PV Nodes only return on exact hit, ignore check mates for now
				if !pvNode || (te.Score > alpha && te.Score < beta) || te.Score > -inf+100 || te.Score < inf-100 {
					return te.Score, nil
				}
			}
		}
		if isGoodGuess {
			ttMove = te.BestMove
		}
	}

	// Null Move Pruning
	// https://www.chessprogramming.org/Null_Move_Pruning
	if depth > 2 && canNull && evaluation.Evaluation(pos) > beta && !isInCheck && !pvNode {
		ep := pos.MakeNullMove()
		adaptiveDepth := uint8(2)
		if depth > 6 {
			adaptiveDepth = 3
		}
		v, err := s.negamax(ctx, pos, -beta, -beta+1, maxDepth-adaptiveDepth, ply+1, false, &pvline.PVLine{}, false, false)
		pos.UnMakeNullMove(ep)
		if err != nil {
			return 0, err
		}
		if v < beta {
			return beta, nil
		}
	}

	oldAlpha := alpha
	potentialPVLine := pvline.PVLine{}
	var prevPos position.Position
	var bestMove move.Move
	var legalMoves uint8

	// Generate all moves and order them
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalMoves(moves)
	s.orderMoves(pos, moves, pvMove, ttMove, ply)

	for i := uint8(0); i < moves.Length(); i++ {
		m := moves.Get(i)
		prevPos = *pos
		pos.MakeMove(*m)
		if !pos.IsLegal() {
			*pos = prevPos
			continue
		}

		legalMoves++
		score, err := s.negamax(ctx, pos, -beta, -alpha, maxDepth, ply+1, pvNode, &potentialPVLine, false, true)
		*pos = prevPos
		if err != nil {
			return 0, err
		}
		score = -score

		if score >= beta {
			transpositiontable.TTable.PotentiallySave(pos.ZobristHash, bestMove, depth, beta, transpositiontable.BetaNode)
			s.betaCutOffs++
			return beta, nil
		}

		if score > alpha {
			alpha = score
			bestMove = *m
			// Update Killer Move, if quiet move
			if pos.GetPiece(m.GetTargetSquare()) == types.NO_PIECE {
				if s.KillerMoves[ply][0] != bestMove {
					s.KillerMoves[ply][1] = s.KillerMoves[ply][0]
				}
				s.KillerMoves[ply][0] = bestMove
			}

			pvl.Update(bestMove, &potentialPVLine)
		}

		potentialPVLine.Reset()
	}

	// There are no legal moves, so it is either a checkmate or a stalemate
	if legalMoves == 0 {
		// Checkmate, set lowest possible value, but increase by the number of plys,
		// so the engine is looking for shorter mates.
		if pos.IsInCheck(pos.SideToMove) {
			return mateValue, nil
		}
		// stalemate
		return 0, nil
	}

	nt := transpositiontable.PVNode
	if oldAlpha == alpha {
		s.alphaCutOffs++
		nt = transpositiontable.AlphaNode
	}
	transpositiontable.TTable.PotentiallySave(pos.ZobristHash, bestMove, depth, alpha, nt)
	return alpha, nil
}

func (s *Search) quiescence(ctx context.Context, pos *position.Position, alpha, beta int, ply uint8) (int, error) {
	s.nodes++
	s.quiescenceNodes++
	// value to info channel and check if we are done
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	stand_pat := evaluation.Evaluation(pos)
	if stand_pat >= beta {
		s.betaCutOffs++
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
	s.orderMoves(pos, moves, move.NullMove, move.NullMove, ply)
	for i := uint8(0); i < moves.Length(); i++ {
		m := moves.Get(i)

		// Delta Pruning, https://www.chessprogramming.org/Delta_Pruning
		if m.GetMoveType() == move.NORMAL &&
			!evaluation.IsEndgame(pos) &&
			stand_pat+evaluation.PieceValue[pos.PiecesBoard[m.GetTargetSquare()].Type()]+200 < alpha {
			continue
		}

		// Static Exchange Evaluation, https://www.chessprogramming.org/Static_Exchange_Evaluation
		if m.GetMoveType() == move.NORMAL && evaluation.StaticExchangeEvaluation(pos, m) < 0 {
			continue
		}

		prevPos = *pos
		pos.MakeMove(*m)
		if !pos.IsLegal() {
			*pos = prevPos
			continue
		}
		score, err := s.quiescence(ctx, pos, -beta, -alpha, ply+1)
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

	// Current puffer of 10% for teardown process
	movetime = movetime - movetime/10

	fmt.Printf("info string calculated timeout %v\n", movetime)
	return context.WithTimeout(ctx, time.Duration(movetime)*time.Millisecond)
}

func (s *Search) isRepetition(hash uint64) bool {
	for i := 0; i < s.SearchHistoryPly; i++ {
		if s.SearchHistory[i] == hash {
			return true
		}
	}
	return false
}
