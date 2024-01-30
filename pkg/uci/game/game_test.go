package game

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
			tokens: strings.Split("startpos", " "),
			want:   "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			name:   "startpos with moves",
			tokens: strings.Split("startpos moves e2e4 e7e5", " "),
			want:   "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
		},
		{
			name:   "fen string",
			tokens: strings.Split("fen rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2", " "),
			want:   "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2",
		},
		{
			name:   "fen string with moves",
			tokens: strings.Split("fen rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2 moves d2d3 d7d6", " "),
			want:   "rnbqkbnr/ppp2ppp/3p4/4p3/4P3/3P4/PPP2PPP/RNBQKBNR w KQkq - 0 3",
		},
		{
			name:   "broken on command 1",
			tokens: strings.Split("startpos moves e2e4 b8c6 d2d3 a8b8 g1h3 b8a8 f2f3 a8b8 c2c4 b8a8 d1a4 a8b8 a4b4 c6b4 b2b3 b4c2 e1d2 c2a1 b1c3 d7d6 e4e5 c8h3 g2h3 d6e5 c3e4", " "),
			want:   "1r1qkbnr/ppp1pppp/8/4p3/2P1N3/1P1P1P1P/P2K3P/n1B2B1R b k - 1 13",
		},
		{
			name:   "broken on command 2",
			tokens: strings.Split("startpos moves e2e4 c7c5 g1f3 d7d6 d2d4 c5d4 f3d4 g8f6 b1c3 a7a6", " "),
			want:   "rnbqkb1r/1p2pppp/p2p1n2/8/3NP3/2N5/PPP2PPP/R1BQKB1R w KQkq - 0 6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newGameImpl()
			g.NewPosition(tt.tokens)
			assert.Equal(t, tt.want, g.position.ToFen())
		})
	}
}

func Test_parseGo(t *testing.T) {
	tests := []struct {
		name   string
		tokens []string
		want   goParameter
	}{
		{
			name:   "broken-1",
			tokens: strings.Split("wtime 60000 btime 60000 movestogo 35", " "),
			want: goParameter{
				wtime:     60000,
				btime:     60000,
				movesToGo: 35,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseGo(tt.tokens))
		})
	}
}
