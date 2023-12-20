package knight

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
	. "github.com/shaardie/clemens/pkg/types"
)

func TestAttacksBySquare(t *testing.T) {
	type args struct {
		square int
	}
	tests := []struct {
		name string
		args args
		want bitboard.Bitboard
	}{
		{
			name: "middle",
			args: args{
				square: SQUARE_D4,
			},
			want: bitboard.BitBySquares(
				SQUARE_B3,
				SQUARE_B5,
				SQUARE_C2,
				SQUARE_C6,
				SQUARE_E2,
				SQUARE_E6,
				SQUARE_F3,
				SQUARE_F5,
			),
		},
		{
			name: "edge",
			args: args{
				square: SQUARE_B1,
			},
			want: bitboard.BitBySquares(
				SQUARE_A3,
				SQUARE_C3,
				SQUARE_D2,
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
