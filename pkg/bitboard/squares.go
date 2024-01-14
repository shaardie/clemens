package bitboard

import "github.com/shaardie/clemens/pkg/types"

func BitBySquares(squares ...int) Bitboard {
	b := Empty
	for _, s := range squares {
		b |= One << s
	}
	return b
}

func RankMaskOfSquare(square int) Bitboard {
	return RankMask1 << (8 * types.RankOfSquare(square))
}

func BitboardFromRankAndFile(rank int, file int) Bitboard {
	return BitBySquares(types.SquareFromRankAndFile(rank, file))
}
