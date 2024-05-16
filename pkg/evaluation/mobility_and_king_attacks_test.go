package evaluation

import (
	"testing"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_evalMobilityAndKingAttackValueByColor(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		we   types.Color
		want int16
	}{
		{
			name: "kiwipete, just for testing changes in the function",
			fen:  "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			we:   types.WHITE,
			want: 59,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := position.NewFromFen(tt.fen)
			require.NoError(t, err)
			assert.Equal(t, tt.want, evalMobilityAndKingAttackValueByColor(pos, tt.we))
		})
	}
}
