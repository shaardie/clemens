package move

import (
	"testing"

	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestMove(t *testing.T) {
	var m Move
	m.SetSourceSquare(types.SQUARE_B8)
	m.SetDestinationSquare(types.SQUARE_H8)
	m.SetMoveType(PROMOTION)
	m.SetPromitionPieceType(types.QUEEN)
	assert.Equal(t, types.SQUARE_B8, m.GetSourceSquare())
	assert.Equal(t, types.SQUARE_H8, m.GetDestinationSquare())
	assert.Equal(t, PROMOTION, m.GetMoveType())
	assert.Equal(t, types.QUEEN, m.GetPromitionPieceType())
}
