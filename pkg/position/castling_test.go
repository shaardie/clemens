package position

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosition_CanCastleNow(t *testing.T) {
	tests := []struct {
		name string
		fen  string
		c    Castling
		want bool
	}{
		{
			name: "default",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			c:    WHITE_CASTLING_KING,
			want: false,
		},
		{
			name: "king or rook moved",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1",
			c:    WHITE_CASTLING_KING,
			want: false,
		},
		{
			name: "white queen castling possible",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R3KBNR w KQkq - 0 1",
			c:    WHITE_CASTLING_QUEEN,
			want: true,
		},
		{
			name: "black king castling possible",
			fen:  "rnbqk2r/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			c:    BLACK_CASTLING_KING,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := NewFromFen(tt.fen)
			assert.NoError(t, err)
			if got := pos.CanCastleNow(tt.c); got != tt.want {
				t.Errorf("Position.CanCastleNow() = %v, want %v", got, tt.want)
			}
		})
	}
}
