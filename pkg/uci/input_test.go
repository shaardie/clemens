package uci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_prepareInput(t *testing.T) {
	tests := []struct {
		s    string
		want []string
	}{
		{
			s:    "debug on",
			want: []string{"debug", "on"},
		},
		{
			s:    "   debug     on  ",
			want: []string{"debug", "on"},
		},
		{
			s:    "\t  debug \t  \t\ton\t  ",
			want: []string{"debug", "on"},
		},
		{
			s:    "joho debug on",
			want: []string{"debug", "on"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			assert.Equal(t, tt.want, prepareInput(tt.s))
		})
	}
}
