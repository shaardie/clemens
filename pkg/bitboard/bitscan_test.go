package bitboard

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeastSignificantOneBit(t *testing.T) {
	assert.Panics(t, func() { LeastSignificantOneBit(Empty) })
	for i := 0; i < 64; i++ {
		t.Run(fmt.Sprintf("LSB %v", i), func(t *testing.T) {
			if got := LeastSignificantOneBit(1 << i); got != i {
				t.Errorf("LeastSignificantOneBit() = %v, want %v", got, i)
			}
		})
	}
}

func TestMostSignificantOneBit(t *testing.T) {
	assert.Panics(t, func() { MostSignificantOneBit(Empty) })
	for i := 0; i < 64; i++ {
		t.Run(fmt.Sprintf("LSB %v", i), func(t *testing.T) {
			if got := MostSignificantOneBit(1 << (63 - i)); got != i {
				t.Errorf("LeastSignificantOneBit() = %v, want %v", got, i)
			}
		})
	}
}
