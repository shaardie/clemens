package utils

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

func TestSlidingAttacks(t *testing.T) {
	type args struct {
		square     int
		directions []func(bitboard.Bitboard) bitboard.Bitboard
		occupied   bitboard.Bitboard
	}
	tests := []struct {
		name string
		args args
		want bitboard.Bitboard
	}{
		{
			name: "north",
			args: args{
				square: types.SQUARE_B1,
				directions: []func(bitboard.Bitboard) bitboard.Bitboard{
					bitboard.NorthOne,
				},
				occupied: bitboard.BitBySquares(types.SQUARE_B3),
			},
			want: bitboard.BitBySquares(
				types.SQUARE_B2,
				types.SQUARE_B3,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SlidingAttacks(tt.args.square, tt.args.directions, tt.args.occupied); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlidingAttacks() = %v, want %v", got, tt.want)
			}
		})
	}
}
