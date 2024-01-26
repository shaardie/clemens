package uci

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_game_newPosition(t *testing.T) {
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
			want:   "1r1qkbnr/ppp1pppp/8/4p3/2P1N3/1P1P1P1P/P2K3P/n1B2B1R b k - 1 13",
		},
		{
			name:   "broken on command 1",
			tokens: strings.Split("startpos moves e2e4 h7h6 g1f3 a7a6 f1c4 d7d5 c4b3 d5e4 f3e5 e7e6 d2d3 e4d3 d1d3 f8b4 c1d2 d8d3 c2d3 c7c5 d2b4 c5b4 e1g1 f7f6 e5g6 h8h7 b1d2 b7b6 a1e1 e8f7 g6f4 g7g5 f4h5 f7g6 g2g4 h7h8 b3e6 c8e6 e1e6 b8d7 d2e4 g6h7 e4f6 g8f6 h5f6 d7f6 e6f6 h8b8 f1e1 b6b5 e1e7 h7h8 f6h6 h8g8 h6g6 g8f8 e7e6 f8f7 e6f6 f7e7 f6e6 e7d7 e6f6 d7e7 d3d4 b8g8 f6e6 e7f8 g6g8 f8g8 d4d5 g8f8 d5d6", " "),
			want:   "r4k2/8/p2PR3/1p4p1/1p4P1/8/PP3P1P/6K1 b - - 0 36",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.newPosition(tt.tokens)
			assert.Equal(t, tt.want, s.position.ToFen())
		})
	}
}
