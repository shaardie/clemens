package bitboard

import (
	"reflect"
	"testing"
)

func TestIsolatingSubsets(t *testing.T) {
	type args struct {
		b Bitboard
	}
	tests := []struct {
		name string
		args args
		want []Bitboard
	}{
		{
			name: "default",
			args: args{
				0b1001101,
			},
			want: []Bitboard{
				0b0000001,
				0b0000100,
				0b0001000,
				0b1000000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsolatingSubsets(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsolatingSubsets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSquareIndexSerialization(t *testing.T) {
	type args struct {
		b Bitboard
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "default",
			args: args{
				0b1001101,
			},
			want: []int{
				0, 2, 3, 6,
			},
		},
		{
			name: "empty",
			args: args{
				Empty,
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SquareIndexSerialization(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SquareIndexSerialization() = %v, want %v", got, tt.want)
			}
		})
	}
}
