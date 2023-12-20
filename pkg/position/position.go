package position

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	. "github.com/shaardie/clemens/pkg/types"
)

type Position struct {
	Pieces      [COLOR_NUMBER][PIECE_NUMBER]bitboard.Bitboard
	blackByType [COLOR_NUMBER][PIECE_NUMBER]bitboard.Bitboard
}
