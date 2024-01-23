package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_handlePositionInput(t *testing.T) {
	tests := []struct {
		name   string
		tokens []string
		want   string
	}{
		{
			name:   "startpos",
			tokens: []string{"startpos"},
			want:   "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			name:   "startpos with moves",
			tokens: []string{"startpos", "moves", "e2e4", "e7e5"},
			want:   "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
		},
		{
			name:   "fen string",
			tokens: []string{"fen", "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR", "w", "KQkq", "e6", "0", "2"},
			want:   "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
		},
		{
			name:   "fen string with moves",
			tokens: []string{"fen", "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR", "w", "KQkq", "e6", "0", "2", "moves", "d2d3", "d7d6"},
			want:   "rnbqkbnr/ppp2ppp/3p4/4p3/4P3/3P4/PPP2PPP/RNBQKBNR w KQkq - 0 3",
		},
		{
			name:   "broken on command 1",
			tokens: strings.Split("startpos moves e2e4 b8c6 d2d3 a8b8 g1h3 b8a8 f2f3 a8b8 c2c4 b8a8 d1a4 a8b8 a4b4 c6b4 b2b3 b4c2 e1d2 c2a1 b1c3 d7d6 e4e5 c8h3 g2h3 d6e5 c3e4", " "),
			want:   "1r1qkbnr/ppp1pppp/8/4p3/2P1N3/1P1P1P1P/P2K3P/n1B2B1R b k - 0 13",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := handlePositionInput(tt.tokens)
			s := pos.ToFen()
			assert.Equal(t, s, tt.want)
		})
	}
}
