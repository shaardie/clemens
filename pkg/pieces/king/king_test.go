package king

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
				SQUARE_C3,
				SQUARE_C4,
				SQUARE_C5,
				SQUARE_D3,
				SQUARE_D5,
				SQUARE_E3,
				SQUARE_E4,
				SQUARE_E5,
			),
		},
		{
			name: "edge",
			args: args{
				square: SQUARE_A1,
			},
			want: bitboard.BitBySquares(
				SQUARE_A2,
				SQUARE_B1,
				SQUARE_B2,
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
