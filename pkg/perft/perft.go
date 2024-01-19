package perft

import (
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
)

type Result struct {
	Move     move.Move
	Position position.Position
	Leafs    int
}

func Perft(pos *position.Position, depth int) int {
	if depth == 0 {
		return 1
	}
	var nodes int
	moves := pos.GeneratePseudoLegalMoves()
	var prevPos position.Position
	for _, m := range moves {
		prevPos = *pos
		pos.MakeMove(m)
		if pos.IsLegal() {
			nodes += Perft(pos, depth-1)
		}
		*pos = prevPos
	}
	return nodes
}

func Divided(pos *position.Position, depth int) []Result {
	moves := pos.GeneratePseudoLegalMoves()
	if depth == 0 {
		return nil
	}
	results := make([]Result, 0, len(moves))
	for _, m := range moves {
		prevPos := *pos
		pos.MakeMove(m)
		if pos.IsLegal() {
			results = append(results,
				Result{
					Move:     m,
					Position: *pos,
					Leafs:    Perft(pos, depth-1),
				},
			)
		}
		pos = &prevPos
	}
	return results
}
