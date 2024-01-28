package position

import (
	"testing"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/stretchr/testify/assert"
)

func TestPosition_Evaluation(t *testing.T) {
	tests := []struct {
		name                string
		fen                 string
		evalutationOverZero bool
	}{
		{
			name:                "from white",
			fen:                 "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			evalutationOverZero: true,
		},
		{
			name:                "from black",
			fen:                 "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 b kq - 0 1",
			evalutationOverZero: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := NewFromFen(tt.fen)
			assert.NoError(t, err)
			moves := move.NewMoveList()
			eval := pos.Evaluation(moves)
			if tt.evalutationOverZero {
				assert.Positive(t, eval)
			} else {
				assert.Negative(t, eval)
			}
		})
	}
}
