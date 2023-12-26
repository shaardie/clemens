package types

import (
	"testing"
)

func TestPiece_Color(t *testing.T) {
	tests := []struct {
		name string
		p    Piece
		want Color
	}{
		{
			name: "black bishop",
			p:    BLACK_BISHOP,
			want: BLACK,
		},
		{
			name: "white king",
			p:    WHITE_KING,
			want: WHITE,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Color(); got != tt.want {
				t.Errorf("Piece.Color() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPiece_Type(t *testing.T) {
	tests := []struct {
		name string
		p    Piece
		want PieceType
	}{
		{
			name: "black bishop",
			p:    BLACK_BISHOP,
			want: BISHOP,
		},
		{
			name: "white king",
			p:    WHITE_KING,
			want: KING,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Type(); got != tt.want {
				t.Errorf("Piece.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}
