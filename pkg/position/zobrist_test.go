package position

import (
	"testing"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestPosition_ZobristHash(t *testing.T) {
	pos := New()
	assert.NotZero(t, pos.ZobristHash)

	var m move.Move
	m.SetSourceSquare(types.SQUARE_E2)
	m.SetTargetSquare(types.SQUARE_E4)
	pos.MakeMove(m)
	afterMoveHash := pos.ZobristHash
	pos.initZobristHash()
	reinitHash := pos.ZobristHash
	assert.Equal(t, reinitHash, afterMoveHash)
}
