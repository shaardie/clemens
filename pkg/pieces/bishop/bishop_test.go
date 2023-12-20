package bishop

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
	. "github.com/shaardie/clemens/pkg/types"
)

func TestAttacksBySquare(t *testing.T) {
	type args struct {
		square   int
		occupied bitboard.Bitboard
	}
	tests := []struct {
		name string
		args args
		want bitboard.Bitboard
	}{
		{
			name: "middle",
			args: args{
				square:   SQUARE_D4,
				occupied: bitboard.BitBySquares(SQUARE_E5, SQUARE_C3),
			},
			want: bitboard.BitBySquares(
				SQUARE_A7,
				SQUARE_B6,
				SQUARE_C3,
				SQUARE_C5,
				SQUARE_E3,
				SQUARE_E5,
				SQUARE_F2,
				SQUARE_G1,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AttacksBySquare(tt.args.square, tt.args.occupied); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AttacksBySquare() = %v, want %v", got, tt.want)
			}
		})
	}
}
