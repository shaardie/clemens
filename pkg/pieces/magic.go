package pieces

import "github.com/shaardie/clemens/pkg/bitboard"

type magic struct {
	mask    bitboard.Bitboard
	magic   bitboard.Bitboard
	attacks []bitboard.Bitboard
	shift   uint
}

// func initMagic(t pieceType, magics []magic) {
// 	seeds := [][bitboard.RankNumber]int{
// 		{8977, 44560, 54343, 38998, 5731, 95205, 104912, 17020},
// 		{728, 10316, 55013, 32803, 12281, 15100, 16645, 255},
// 	}
// 	var (
// 		edges, b             bitboard.Bitboard
// 		occupancy, reference [4096]bitboard.Bitboard
// 		epoch                [4096]int
// 		cnt, size            int
// 	)

// 	for square := 0; square < 64; square++ {
// 		edges = (bitboard.RankMask1 | bitboard.RankMask8) & ^bitboard.RankMaskOfSquare(square)

// 		magics[s].mask =
// 	}

// }
