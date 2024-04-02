package magic

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
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

func Init(attacksFunc func(square uint8, occupied bitboard.Bitboard) bitboard.Bitboard, rand func() uint64) (magics [types.SQUARE_NUMBER]Magic) {
	// array of all attacks, splittet in slices per square.
	// this is faster than having separate arrays.
	table := []bitboard.Bitboard{}

	for square := types.SQUARE_A1; square < types.SQUARE_NUMBER; square++ {
		m := Magic{}

		// First we calculate the mask and shift
		// The edges are not relevant for the occupancy,
		// because the squares can be accessed independent from the occupency.
		squareRankMask := bitboard.RankMask1 << bitboard.Bitboard(8*types.RankOfSquare(square))
		rankedges := (bitboard.RankMask1 | bitboard.RankMask8) & ^squareRankMask
		squareFileMask := bitboard.FileMaskA << bitboard.Bitboard(types.FileOfSquare(square))
		fileedges := (bitboard.FileMaskA | bitboard.FileMaskH) &^ squareFileMask
		edges := rankedges | fileedges
		m.Mask = attacksFunc(square, 0) & ^edges
		m.Shift = uint(64 - m.Mask.PopulationCount())

		occupancies := bitboard.AllSubnetsOf(m.Mask)
		size := len(occupancies)

		attacks := make([]bitboard.Bitboard, size)
		for i, occupancy := range occupancies {
			attacks[i] = attacksFunc(square, occupancy)
		}

		// Slice from global table
		oldTableSize := len(table)
		table = append(table, make([]bitboard.Bitboard, size)...)
		newTableSize := len(table)
		m.Attacks = table[oldTableSize:newTableSize]

		complete := false
		// TODO speed up this function
		for !complete {
			// Find small magic
			for {
				m.Magic = bitboard.Bitboard(rand() & rand() & rand())
				if ((m.Magic * m.Mask) >> 56).PopulationCount() < 6 {
					break
				}
			}
			setBitboardToEmpty(m.Attacks)
			for i, occupancy := range occupancies {
				idx := m.Index(occupancy)

				if m.Attacks[idx] != 0 {
					break
				}
				m.Attacks[idx] = attacks[i]
				if i == size-1 {
					complete = true
				}
			}
		}
		// Set magic
		magics[square] = m
	}
	return magics
}

// setBitboardToEmpty sets all entries to 0. It is faster than iterating.
func setBitboardToEmpty(a []bitboard.Bitboard) {
	if len(a) == 0 {
		return
	}
	a[0] = 0
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}
