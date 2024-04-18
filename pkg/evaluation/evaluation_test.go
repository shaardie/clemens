package evaluation

import (
	"testing"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/stretchr/testify/assert"
)

func TestPosition_Evaluation(t *testing.T) {
	tests := []struct {
		name                string
		fen                 string
		evalutationOverZero bool
		useExactValue       bool
		exactValue          int16
	}{
		{
			name:                "from white",
			fen:                 "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			evalutationOverZero: true,
			useExactValue:       true,
			exactValue:          153,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			assert.NoError(t, err)
			eval := Evaluation(pos)
			if tt.evalutationOverZero {
				assert.Positive(t, eval)
			} else {
				assert.Negative(t, eval)
			}
			if tt.useExactValue {
				assert.Equal(t, tt.exactValue, eval)
			}
		})
	}
}
