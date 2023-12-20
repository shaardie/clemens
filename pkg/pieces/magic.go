package pieces

import (
	"math/rand"

	"github.com/shaardie/clemens/pkg/bitboard"
)

var (
	RookTable    []bitboard.Bitboard
	RookMagics   [bitboard.SQUARE_NUMBER]magic
	BishopTable  []bitboard.Bitboard
	BishopMagics [bitboard.SQUARE_NUMBER]magic
)

type magic struct {
	Mask    bitboard.Bitboard
	Magic   bitboard.Bitboard
	Attacks []bitboard.Bitboard
	shift   uint
}

// Index computes attack Index
func (m magic) Index(occupied bitboard.Bitboard) uint {
	return uint(((occupied & m.Mask) * m.Magic) >> m.shift)
}

func initMagics() {
	BishopTable, BishopMagics = initMagic(BishopType)
	RookTable, RookMagics = initMagic(RookType)
}

func initMagic(t pieceType) (table []bitboard.Bitboard, magics [bitboard.SQUARE_NUMBER]magic) {
	for square := 0; square < 64; square++ {
		edges := (bitboard.RankMask1 | bitboard.RankMask8) & ^bitboard.RankMaskOfSquare(square)

		m := magic{}
		m.Mask = SlidingAttacks(t, square, 0) & ^edges
		m.shift = uint(64 - m.Mask.PopulationCount())

		b := bitboard.Empty
		occupancy := make([]bitboard.Bitboard, 0, 4096)
		reference := make([]bitboard.Bitboard, 0, 4096)
		for {
			occupancy = append(occupancy, b)
			reference = append(reference, SlidingAttacks(t, square, b))

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
