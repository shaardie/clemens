package evaluation

import (
	"testing"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestStaticExchangeEvaluation(t *testing.T) {
	tests := []struct {
		name         string
		fen          string
		sourceSquare int
		targetSquare int
		value        int16
	}{
		{
			name:         "position 1",
			fen:          "1k1r4/1pp4p/p7/4p3/8/P5P1/1PP4P/2K1R3 w - - 0 1",
			sourceSquare: types.SQUARE_E1,
			targetSquare: types.SQUARE_E5,
			value:        100,
		},
		{
			name:         "position 2",
			fen:          "1k1r3q/1ppn3p/p4b2/4p3/8/P2N2P1/1PP1R1BP/2K1Q3 w - - 0 1",
			sourceSquare: types.SQUARE_D3,
			targetSquare: types.SQUARE_E5,
			value:        -200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			assert.NoError(t, err)
			var m move.Move
			m.SetSourceSquare(tt.sourceSquare)
			m.SetTargetSquare(tt.targetSquare)
			eval := StaticExchangeEvaluation(pos, &m)
			assert.Equal(t, tt.value, eval)
		})
	}
}
