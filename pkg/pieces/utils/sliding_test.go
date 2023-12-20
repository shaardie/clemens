package utils

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/bitboard"
)

// func TestSlidingAttxfsdfacks(t *testing.T) {
// 	type args struct {
// 		t        pieceType
// 		square   int
// 		occupied bitboard.Bitboard
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bitboard.Bitboard
// 	}{
// 		{
// 			name: "rook",
// 			args: args{t: RookType, square: bitboard.SQUARE_A1, occupied: bitboard.RankMask6 | bitboard.FileMaskE},
// 			want: bitboard.BitBySquares(
// 				bitboard.SQUARE_A2,
// 				bitboard.SQUARE_A3,
// 				bitboard.SQUARE_A4,
// 				bitboard.SQUARE_A5,
// 				bitboard.SQUARE_A6,
// 				bitboard.SQUARE_B1,
// 				bitboard.SQUARE_C1,
// 				bitboard.SQUARE_D1,
// 				bitboard.SQUARE_E1,
// 			),
// 		},
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := SlidingAttacks(tt.args.t, tt.args.square, tt.args.occupied); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("RookAttacks() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

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
				square: bitboard.SQUARE_B1,
				directions: []func(bitboard.Bitboard) bitboard.Bitboard{
					bitboard.NorthOne,
				},
				occupied: bitboard.BitBySquares(bitboard.SQUARE_B3),
			},
			want: bitboard.BitBySquares(
				bitboard.SQUARE_B2,
				bitboard.SQUARE_B3,
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
