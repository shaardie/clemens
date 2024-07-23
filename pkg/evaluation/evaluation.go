package evaluation

import (
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

// Evaluation evaluates the position.
func Evaluation(pos *position.Position) int16 {
	if isDraw(pos) {
		return Contempt(pos)
	}
	v := pos.Accumulator.Evaluate(pos.SideToMove)
	if pos.SideToMove == types.BLACK {
		return -v
	}
	return v
}
