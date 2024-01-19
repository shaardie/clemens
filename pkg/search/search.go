package search

import (
	"github.com/shaardie/clemens/pkg/move"
	"github.com/shaardie/clemens/pkg/position"
)

type node struct {
	pos        *position.Position
	evaluation float64
	leafs      map[move.Move]*node
}

func (n *node) generateLeafs() {
}
