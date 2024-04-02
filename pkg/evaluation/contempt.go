package evaluation

import (
	"github.com/shaardie/clemens/pkg/position"
)

// Calculates the Contempt Factor for drawish positions, see https://www.chessprogramming.org/Contempt_Factor
func Contempt(pos *position.Position) int16 {
	if IsEndgame(pos) {
		return 0
	}

	// We do not resign too early
	return 400
}
