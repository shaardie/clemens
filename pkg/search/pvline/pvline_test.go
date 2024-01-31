package pvline

import (
	"testing"

	"github.com/shaardie/clemens/pkg/move"
	"github.com/stretchr/testify/assert"
)

func TestPVLine_Update(t *testing.T) {
	tests := []struct {
		name    string
		oldLine PVLine
		move    move.Move
		newLine PVLine
		wanted  PVLine
	}{
		{
			name:    "regular update",
			oldLine: PVLine{[]move.Move{1, 2, 3}},
			move:    4,
			newLine: PVLine{[]move.Move{5, 6, 7}},
			wanted:  PVLine{[]move.Move{4, 5, 6, 7}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.oldLine.Update(tt.move, &tt.newLine)
			assert.Equal(t, tt.wanted, tt.oldLine)
		})
	}
}
