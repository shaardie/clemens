package bitboard

import "fmt"

type Bitboard uint64

func (b Bitboard) String() string {
	return fmt.Sprintf("%064b", b)
}

func (b Bitboard) PrettyString() string {
	s := "\n+----+----+----+----+----+----+----+----+\n"
	for rank := rank8; rank >= rank1; rank-- {
		for file := file1; file <= file8; file++ {
			if (b & BitboardFromRankAndFile(rank, file)) != Empty {
				s += "| wX "
			} else {
				s += "|    "
			}
			fmt.Println()
		}
		s += fmt.Sprintf("| %d\n+----+----+----+----+----+----+----+----+\n", rank+1)
	}
	s += "  a   b   c   d   e   f   g   h\n"
	return s
}
