package evaluation

import (
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

// Calculates the Contempt Factor for drawish positions, see https://www.chessprogramming.org/Contempt_Factor
func Contempt(pos *position.Position) int {
	if IsEndgame(pos) {
		return 0
	}

	// We do not resign too early
	score := -400

	// Make the result side aware
	if pos.SideToMove == types.BLACK {
		score *= -1
	}
	return score
}
