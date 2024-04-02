package transpositiontable

import (
	"testing"

	"github.com/shaardie/clemens/pkg/evaluation"
	"github.com/shaardie/clemens/pkg/move"
	"github.com/stretchr/testify/assert"
)

func TestTT(t *testing.T) {

	PotentiallySave(1, 1, 2, 1, PVNode, 0)
	score, use, m := Get(1, -evaluation.INF, evaluation.INF, 1, 0)
	assert.Equal(t, int16(1), score)
	assert.Equal(t, true, use)
	assert.Equal(t, move.Move(1), m)

	// Another Hash
	_, use, m = Get(2, -evaluation.INF, evaluation.INF, 1, 0)
	assert.Equal(t, false, use)
	assert.Equal(t, move.NullMove, m)

	// Deeper
	PotentiallySave(1, 2, 3, 2, PVNode, 0)
	score, use, m = Get(1, -evaluation.INF, evaluation.INF, 1, 0)
	assert.Equal(t, int16(2), score)
	assert.Equal(t, true, use)
	assert.Equal(t, move.Move(2), m)

	// Not so deep, so ignored
	PotentiallySave(1, 3, 2, 3, PVNode, 0)
	score, use, m = Get(1, -evaluation.INF, evaluation.INF, 1, 0)
	assert.Equal(t, int16(2), score)
	assert.Equal(t, true, use)
	assert.Equal(t, move.Move(2), m)

	// Return alpha
	PotentiallySave(1, 2, 3, 2, AlphaNode, 0)
	score, use, m = Get(1, 10, evaluation.INF, 1, 0)
	assert.Equal(t, int16(10), score)
	assert.Equal(t, true, use)
	assert.Equal(t, move.Move(2), m)

	// Return beta
	PotentiallySave(1, 2, 3, 2, BetaNode, 0)
	score, use, m = Get(1, -evaluation.INF, -1, 1, 0)
	assert.Equal(t, int16(-1), score)
	assert.Equal(t, true, use)
	assert.Equal(t, move.Move(2), m)

	// Adjusted for Mate
	PotentiallySave(1, 2, 3, evaluation.INF-3, PVNode, 0)
	score, use, m = Get(1, -evaluation.INF, evaluation.INF, 1, 7)
	assert.Equal(t, evaluation.INF-10, score)
	assert.Equal(t, true, use)
	assert.Equal(t, move.Move(2), m)
	PotentiallySave(1, 2, 3, -evaluation.INF+3, PVNode, 0)
	score, use, m = Get(1, -evaluation.INF, evaluation.INF, 1, 7)
	assert.Equal(t, -evaluation.INF+10, score)
	assert.Equal(t, true, use)
	assert.Equal(t, move.Move(2), m)
}
