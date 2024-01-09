package position

import (
	"testing"

	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestPosition_IsCheck(t *testing.T) {

	tests := []struct {
		name  string
		fen   string
		color types.Color
		want  bool
	}{
		{
			name:  "is check",
			fen:   "3K4/8/8/2B5/3k4/8/8/8 w - - 0 1",
			color: types.BLACK,
			want:  true,
		},
		{
			name:  "is not check",
			fen:   "3K4/8/8/3B4/3k4/8/8/8 b - - 0 1",
			color: types.BLACK,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := NewFromFen(tt.fen)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, pos.IsCheck(tt.color))
		})
	}
}
