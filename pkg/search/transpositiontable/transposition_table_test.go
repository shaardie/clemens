package transpositiontable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranspositionTable_Get(t *testing.T) {
	tests := []struct {
		name  string
		saves []TranspositionEntry
		arg   uint64
		arg1  uint8
		want  TranspositionEntry
		found bool
	}{
		{
			name: "regular",
			saves: []TranspositionEntry{
				{
					ZobristHash: 1,
					BestMove:    1,
					Depth:       2,
					Score:       1,
					NodeType:    PVNode,
				},
			},
			arg:  1,
			arg1: 1,
			want: TranspositionEntry{
				ZobristHash: 1,
				BestMove:    1,
				Depth:       2,
				Score:       1,
				NodeType:    PVNode,
			},
			found: true,
		},
		{
			name: "wrong depth",
			saves: []TranspositionEntry{
				{
					ZobristHash: 1,
					BestMove:    1,
					Depth:       2,
					Score:       1,
					NodeType:    PVNode,
				},
			},
			arg:  1,
			arg1: 3,
		},
		{
			name: "override",
			saves: []TranspositionEntry{
				{
					ZobristHash: 1,
					BestMove:    1,
					Depth:       2,
					Score:       1,
					NodeType:    PVNode,
				},
				{
					ZobristHash: 1,
					BestMove:    1,
					Depth:       3,
					Score:       2,
					NodeType:    PVNode,
				},
			},
			arg:  1,
			arg1: 1,
			want: TranspositionEntry{
				ZobristHash: 1,
				BestMove:    1,
				Depth:       3,
				Score:       2,
				NodeType:    PVNode,
			},
			found: true,
		},
		{
			name: "no override",
			saves: []TranspositionEntry{
				{
					ZobristHash: 1,
					BestMove:    1,
					Depth:       2,
					Score:       1,
					NodeType:    PVNode,
				},
				{
					ZobristHash: 1,
					BestMove:    1,
					Depth:       1,
					Score:       2,
					NodeType:    PVNode,
				},
			},
			arg:  1,
			arg1: 1,
			want: TranspositionEntry{
				ZobristHash: 1,
				BestMove:    1,
				Depth:       2,
				Score:       1,
				NodeType:    PVNode,
			},
			found: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reset()
			for _, s := range tt.saves {
				TTable.PotentiallySave(s.ZobristHash, s.BestMove, s.Depth, s.Score, s.NodeType)
			}
			entry, found := TTable.Get(tt.arg, tt.arg1)
			assert.Equal(t, tt.found, found)
			if found {
				assert.Equal(t, entry, tt.want)
			}
		})
	}
}
