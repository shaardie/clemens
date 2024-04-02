package knight

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

func TestAttacksBySquare(t *testing.T) {
	type args struct {
		square uint8
	}
	tests := []struct {
		name string
		args args
		want bitboard.Bitboard
	}{
		{
			name: "middle",
			args: args{
				square: types.SQUARE_D4,
			},
			want: bitboard.BitBySquares(
				types.SQUARE_B3,
				types.SQUARE_B5,
				types.SQUARE_C2,
				types.SQUARE_C6,
				types.SQUARE_E2,
				types.SQUARE_E6,
				types.SQUARE_F3,
				types.SQUARE_F5,
			),
		},
		{
			name: "edge",
			args: args{
				square: types.SQUARE_B1,
			},
			want: bitboard.BitBySquares(
				types.SQUARE_A3,
				types.SQUARE_C3,
				types.SQUARE_D2,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AttacksBySquare(tt.args.square); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AttacksBySquare() = %v, want %v", got, tt.want)
			}
		})
	}
}
