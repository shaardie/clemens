package evaluation

import (
	"testing"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/stretchr/testify/assert"
)

func TestPosition_EvalDraw(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		draw bool
	}{
		{
			name: "kiwipete",
			fen:  "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			draw: false,
		},
		{
			name: "only kings",
			fen:  "k7/8/8/8/8/8/8/K7 w - - 0 1",
			draw: true,
		},
		{
			name: "existing pawn",
			fen:  "k7/8/8/8/8/5P2/8/K7 w - - 0 1",
			draw: false,
		},
		{
			name: "only one minor each",
			fen:  "k7/8/8/8/2B5/6b1/8/K7 w - - 0 1",
			draw: true,
		},
		{
			name: "more than 2 minor",
			fen:  "k7/8/8/8/1BB5/6b1/1N6/K7 w - - 0 1",
			draw: false,
		},
		{
			name: "two bishops against one",
			fen:  "k7/8/8/8/1BB5/6b1/8/K7 w - - 0 1",
			draw: true,
		},
		{
			name: "two bishops against everything else",
			fen:  "k7/8/8/8/1BB5/8/8/K4n2 w - - 0 1",
			draw: false,
		},
		{
			name: "two knights against bare king",
			fen:  "k7/8/8/8/8/8/8/K3nn2 w - - 0 1",
			draw: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			assert.NoError(t, err)
			e := eval{}
			assert.Equal(t, tt.draw, e.isDraw(pos))
		})
	}
}
