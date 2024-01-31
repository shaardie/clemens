package pvline

import (
	"fmt"
	"strings"

	"github.com/shaardie/clemens/pkg/move"
)

type PVLine struct {
	moves []move.Move
}

func (pvline *PVLine) GetBestMove() move.Move {
	if len(pvline.moves) == 0 {
		return move.NullMove
	}
	return pvline.moves[0]
}

func (pvline *PVLine) Update(bestMove move.Move, newLine *PVLine) {
	new := make([]move.Move, len(newLine.moves)+1)
	new[0] = bestMove
	copy(new[1:], newLine.moves)
	pvline.moves = new
}

func (pvline *PVLine) Reset() {
	pvline.moves = nil
}

func (pvline *PVLine) Copy() *PVLine {
	r := PVLine{
		moves: make([]move.Move, len(pvline.moves)),
	}
	copy(r.moves, pvline.moves)
	return &r
}

func (pvline *PVLine) String() string {
	pvStrings := make([]string, len(pvline.moves))
	for i := range pvline.moves {
		pvStrings[i] = fmt.Sprint(pvline.moves[i])
	}
	return strings.Join(pvStrings, " ")
}
