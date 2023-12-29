package position

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosition_ZobristHash(t *testing.T) {
	pos := New()
	assert.Equal(t, uint64(0x9683bf3ac6ef9ea1), pos.ZobristHash())
}
