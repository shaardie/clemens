package magic

import (
	"math/rand"

	"github.com/shaardie/clemens/pkg/bitboard"
	. "github.com/shaardie/clemens/pkg/types"
)

type Magic struct {
	Mask    bitboard.Bitboard
	Magic   bitboard.Bitboard
	Attacks []bitboard.Bitboard
	Shift   uint
}

// Index computes attack Index
func (m Magic) Index(occupied bitboard.Bitboard) uint {
	return uint(((occupied & m.Mask) * m.Magic) >> m.Shift)
}

func Init(attacks func(square int, occupied bitboard.Bitboard) bitboard.Bitboard) (table []bitboard.Bitboard, magics [SQUARE_NUMBER]Magic) {
	for square := SQUARE_A1; square < SQUARE_NUMBER; square++ {
		edges := (bitboard.RankMask1 | bitboard.RankMask8) & ^bitboard.RankMaskOfSquare(square)

		m := Magic{}
		m.Mask = attacks(square, 0) & ^edges
		m.Shift = uint(64 - m.Mask.PopulationCount())

		b := bitboard.Empty
		occupancy := make([]bitboard.Bitboard, 0, 4096)
		reference := make([]bitboard.Bitboard, 0, 4096)
		for {
			occupancy = append(occupancy, b)
			reference = append(reference, attacks(square, b))

			// https://www.chessprogramming.org/Traversing_Subsets_of_a_Set#All_Subsets_of_any_Set
			b = (b - m.Mask) & m.Mask

			if b == bitboard.Empty {
				break
			}
		}

		size := len(occupancy)

		oldTableSize := len(table)
		table = append(table, make([]bitboard.Bitboard, size)...)
		m.Attacks = table[oldTableSize:]

		epoch := make([]int, size)
		iteration := 0
		for i := 0; i < size; {
			// Find small magic
			for {
				m.Magic = sparseRand()
				if ((m.Magic * m.Mask) >> 56).PopulationCount() < 6 {
					break
				}
			}

			iteration++
			for i = 0; i < size; i++ {
				// iterate over the magic indices of the occupancies
				idx := m.Index(occupancy[i])

				// Index already used, see if attack matches, if yes -> failure
				if epoch[idx] == iteration && m.Attacks[idx] != reference[i] {
					break
				}

				epoch[idx] = iteration

				// set attack
				m.Attacks[idx] = reference[i]
				i++
			}
		}

		// Set magic
		magics[square] = m
	}
	return table, magics
}

func sparseRand() bitboard.Bitboard {
	return bitboard.Bitboard(rand.Uint64() & rand.Uint64() & rand.Uint64())

}
