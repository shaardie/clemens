package bitboard

import (
	"fmt"

	. "github.com/shaardie/clemens/pkg/types"
)

type Bitboard uint64

func (b Bitboard) SimpleString() string {
	return fmt.Sprintf("%064b", b)
}

func (b Bitboard) String() string {
	s := "\n+---+---+---+---+---+---+---+---+\n"
	for rank := RANK_8; rank >= RANK_1; rank-- {
		for file := FILE_A; file <= FILE_H; file++ {
			if (b & BitboardFromRankAndFile(rank, file)) != Empty {
				s += "| X "
			} else {
				s += "|   "
			}
		}
		s += fmt.Sprintf("| %d\n+---+---+---+---+---+---+---+---+\n", rank+1)
	}
	s += "  a   b   c   d   e   f   g   h\n"
	return s
}
