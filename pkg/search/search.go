package search

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search/pvline"
	"github.com/shaardie/clemens/pkg/search/transpositiontable"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	inf                        = 100000
	widen_window               = 50
	max_depth            uint8 = 10
	quiescence_max_depth uint8 = 100
	maxTimeInMs                = 10000
)

type Search struct {
	pos                    position.Position
	nodes                  uint64
	betaCutOffs            uint64
	alphaCutOffs           uint64
	quiescenceNodes        uint64
	transpositiontableHits uint64
	m                      *sync.Mutex
	PV                     pvline.PVLine
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

func (s *Search) BestMove() move.Move {
	s.m.Lock()
	defer s.m.Unlock()
	return s.PV.GetBestMove()
}

func NewSearch(pos position.Position) *Search {
	return &Search{
		pos: pos,
		m:   &sync.Mutex{},
	}
}

func (s *Search) Search(ctx context.Context, sp SearchParameter) move.Move {
	ctx, cancel := s.contextFromSearchParameter(ctx, sp)
	depth := max_depth
	if sp.Depth > 0 {
		depth = sp.Depth
	}
	s.SearchIterative(ctx, depth)
	cancel()
	return s.BestMove()
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

		s.m.Lock()
		s.PV = *i.PV.Copy()
		goodGuess = s.PV.GetBestMove()
		s.m.Unlock()
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
	pos := s.pos
	pvl := pvline.PVLine{}
	score, err := s.negamax(ctx, &pos, alpha, beta, depth, 0, true, &pvl, true)
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

func (s *Search) negamax(ctx context.Context, pos *position.Position, alpha, beta int, maxDepth, ply uint8, pvNode bool, pvl *pvline.PVLine, isRoot bool) (int, error) {

	// value to info channel and check if we are done
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

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

	// Evaluate the leaf node
	if ply == maxDepth {
		return s.quiescence(ctx, pos, alpha, beta, ply)
	}

	s.nodes++

	goodGuess := move.NullMove

	// Check if we can use the transition table but not on root
	if isRoot {
		goodGuess = s.PV.GetBestMove()
	} else {
		te, found, isGoodGuess := transpositiontable.TTable.Get(pos.ZobristHash, maxDepth-ply)
		if found {
			s.transpositiontableHits++
			switch te.NodeType {
			case transpositiontable.AlphaNode:
				// return the smaller value of alpha and score
				s.alphaCutOffs++
				return max(te.Score, alpha), nil
			case transpositiontable.BetaNode:
				// return the smaller value of beta and score
				s.betaCutOffs++
				return min(te.Score, beta), nil
			case transpositiontable.PVNode:
				// In PV Nodes only return on exact hit, ignore check mates for now
				if !pvNode || (te.Score > alpha && te.Score < beta) || te.Score > -inf+100 || te.Score < inf-100 {
					return te.Score, nil
				}
			}
		}
		if isGoodGuess {
			goodGuess = te.BestMove
		}
	}
	isRoot = false

	oldAlpha := alpha
	potentialPVLine := pvline.PVLine{}
	var prevPos position.Position
	var bestMove move.Move
	var legalMoves uint8

	// Generate all moves
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalMovesOrdered(moves, goodGuess)
	for i := uint8(0); i < moves.Length(); i++ {
		m := moves.Get(i)
		prevPos = *pos
		pos.MakeMove(*m)
		if !pos.IsLegal() {
			*pos = prevPos
			continue
		}
		legalMoves++
		score, err := s.negamax(ctx, pos, -beta, -alpha, maxDepth, ply+1, pvNode, &potentialPVLine, isRoot)
		*pos = prevPos
		if err != nil {
			return 0, err
		}
		score = -score
		// value to info channel and check if we are done
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		if pvNode {
			pvNode = false
		}

		if score >= beta {
			transpositiontable.TTable.PotentiallySave(pos.ZobristHash, bestMove, maxDepth-ply, beta, transpositiontable.BetaNode)
			s.betaCutOffs++
			return beta, nil
		}

		if score > alpha {
			alpha = score
			bestMove = *m
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
	transpositiontable.TTable.PotentiallySave(pos.ZobristHash, bestMove, maxDepth-ply, alpha, nt)
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

	stand_pat := pos.Evaluation()
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

	// Generate all captures
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalCapturesOrdered(moves, move.NullMove)
	for i := uint8(0); i < moves.Length(); i++ {
		m := moves.Get(i)
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

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s *Search) contextFromSearchParameter(ctx context.Context, sp SearchParameter) (context.Context, context.CancelFunc) {
	// No need for any timeout
	if sp.Infinite {
		return ctx, func() {}
	}

	var t, inc int
	if s.pos.SideToMove == types.BLACK {
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

	// Current bad puffer
	movetime -= 100

	fmt.Printf("info string calculated timeout %v\n", movetime)
	return context.WithTimeout(ctx, time.Duration(movetime)*time.Millisecond)
}
