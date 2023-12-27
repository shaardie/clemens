package position

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFromFen(t *testing.T) {
	pos := New()
	fenPos, err := NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	assert.NoError(t, err)
	if !reflect.DeepEqual(pos, fenPos) {
		t.Errorf("NewFromFen() = %v, want %v", fenPos, pos)
	}
}
