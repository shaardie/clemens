package evaluation

import (
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
)

var tempo = [game_number]int16{20, 20}

func (e *eval) evalTempoBonus(pos *position.Position) {
	var s int16 = 1
	if pos.SideToMove == types.BLACK {
		s = -1
	}

	e.phaseScores[midgame] += s * tempo[midgame]
	e.phaseScores[endgame] += s * tempo[endgame]
}
