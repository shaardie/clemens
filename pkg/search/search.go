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
	max_depth            uint8 = 5
	quiescence_max_depth uint8 = 10
	maxTimeInMs                = 10000
)

type Search struct {
	pos   position.Position
	nodes uint64
	m     *sync.Mutex
	PV    pvline.PVLine
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
	Nodes uint64
	PV    pvline.PVLine
	Score int
}

func (i Info) String() string {
	return fmt.Sprintf("depth=%v, time=%v, nodes=%v, pvline=%v, score=%v", i.Depth, i.Time, i.Nodes, i.PV, i.Score)
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
	for depth <= maxDepth {
		i, err := s.SearchRoot(ctx, depth, alpha, beta)
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
		s.m.Unlock()
		alpha = i.Score - widen_window
		beta = i.Score + widen_window
		depth++

		// Print info
		fmt.Printf("info depth %v score cp %v nodes %v time %v pv %v\n", i.Depth, i.Score, i.Nodes, i.Time, i.PV)

	}
}

func (s *Search) SearchRoot(ctx context.Context, depth uint8, alpha, beta int) (Info, error) {
	start := time.Now()
	pos := s.pos
	pvl := pvline.PVLine{}
	score, err := s.negamax(ctx, &pos, alpha, beta, depth, true, &pvl, true)
	if err != nil {
		return Info{}, err
	}
	return Info{
		Depth: depth,
		Time:  time.Since(start).Milliseconds(),
		Nodes: s.nodes,
		Score: score,
		PV:    *pvl.Copy(),
	}, nil
}

func (s *Search) negamax(ctx context.Context, pos *position.Position, alpha, beta int, depth uint8, pvNode bool, pvl *pvline.PVLine, isRoot bool) (int, error) {
	s.nodes++

	// value to info channel and check if we are done
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	// Evaluate the leaf node
	if depth == 0 {
		return s.quiescence(ctx, pos, alpha, beta)
		// return s.quiescence(ctx, pos, alpha, beta)
	}

	// Check if we can use the transition table but not on root
	if !isRoot {
		te, found := transpositiontable.TTable.Get(pos.ZobristHash, depth)
		if found {
			switch te.NodeType {
			case transpositiontable.AlphaNode:
				// return the smaller value of alpha and score
				return max(te.Score, alpha), nil
			case transpositiontable.BetaNode:
				// return the smaller value of beta and score
				return min(te.Score, beta), nil
			case transpositiontable.PVNode:
				// return exact value
				return te.Score, nil
			}
		}
		isRoot = false
	}

	oldAlpha := alpha
	potentialPVLine := pvline.PVLine{}
	var prevPos position.Position
	var bestMove move.Move

	// Generate all moves
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalMoves(moves)
	for i := uint8(0); i < moves.Length(); i++ {
		m := moves.Get(i)
		prevPos = *pos
		pos.MakeMove(m)
		if pos.IsLegal() {
			score, err := s.negamax(ctx, pos, -beta, -alpha, depth-1, pvNode, &potentialPVLine, isRoot)
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
				transpositiontable.TTable.PotentiallySave(prevPos.ZobristHash, bestMove, depth, beta, transpositiontable.BetaNode)
				return beta, nil
			}

			if score > alpha {
				alpha = score
				bestMove = m
				pvl.Update(bestMove, &potentialPVLine)
			}
		}
		*pos = prevPos
		potentialPVLine.Reset()
	}

	nt := transpositiontable.PVNode
	if oldAlpha == alpha {
		nt = transpositiontable.AlphaNode
	}
	transpositiontable.TTable.PotentiallySave(pos.ZobristHash, bestMove, depth, alpha, nt)
	return alpha, nil
}

func (s *Search) quiescence(ctx context.Context, pos *position.Position, alpha, beta int) (int, error) {
	return s.quiescenceWithMaxDepth(ctx, pos, alpha, beta, quiescence_max_depth)
}

func (s *Search) quiescenceWithMaxDepth(ctx context.Context, pos *position.Position, alpha, beta int, depth uint8) (int, error) {
	s.nodes++

	// value to info channel and check if we are done
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	stand_pat := pos.Evaluation()
	if stand_pat >= beta {
		return beta, nil
	}
	if alpha < stand_pat {
		alpha = stand_pat
	}

	if depth == 0 {
		return alpha, nil
	}
	depth--

	var prevPos position.Position

	// Generate all captures
	moves := move.NewMoveList()
	pos.GeneratePseudoLegalCaptures(moves)
	for i := uint8(0); i < moves.Length(); i++ {
		m := moves.Get(i)
		prevPos = *pos
		pos.MakeMove(m)
		if pos.IsLegal() {
			score, err := s.quiescence(ctx, pos, -beta, -alpha)
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
		*pos = prevPos
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
		fmt.Println("info string calculate until stopped")
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
	fmt.Printf("info string calculated timeout %v\n", movetime)
	return context.WithTimeout(ctx, time.Duration(movetime)*time.Millisecond)
}
