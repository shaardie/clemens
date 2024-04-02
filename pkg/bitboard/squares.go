package bitboard

import "github.com/shaardie/clemens/pkg/types"

func BitBySquares(squares ...uint8) Bitboard {
	b := Empty
	for _, s := range squares {
		b |= One << s
	}
	return b
}

func RankMaskOfSquare(square uint8) Bitboard {
	return RankMask1 << (8 * types.RankOfSquare(square))
}

func FileMaskOfSquare(square uint8) Bitboard {
	return FileMaskA << types.FileOfSquare(square)
}

func BitboardFromRankAndFile(rank uint8, file uint8) Bitboard {
	return BitBySquares(types.SquareFromRankAndFile(rank, file))
}
