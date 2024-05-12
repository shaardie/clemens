package evaluation

import (
	"testing"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_eval_evalPieceSquareTables(t *testing.T) {
	tests := []struct {
		name  string
		fen   string
		wants [2]int16
	}{
		{
			name:  "white king on the first rank",
			fen:   "8/8/8/8/8/8/8/4K3 w - - 0 1",
			wants: [2]int16{0, -30},
		},
		{
			name:  "white king in the middle",
			fen:   "8/8/8/8/3K4/8/8/8 w - - 0 1",
			wants: [2]int16{-40, 40},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			require.NoError(t, err)
			e := &eval{}
			e.evalPieceSquareTables(pos)
			assert.Equal(t, tt.wants, e.phaseScores)
		})
	}
}
