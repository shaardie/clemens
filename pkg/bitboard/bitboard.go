package bitboard

import "fmt"

type Bitboard uint64

func (b Bitboard) String() string {
	return fmt.Sprintf("%064b", b)
}

func (b Bitboard) PrettyString() string {
	s := "\n+---+---+---+---+---+---+---+---+\n"
	for rank := Rank8; rank >= Rank1; rank-- {
		for file := FileA; file <= FileH; file++ {
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
