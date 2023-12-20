package queen

import (
	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/pieces/bishop"
	"github.com/shaardie/clemens/pkg/pieces/rook"
)

// AttacksBySquare returns the attacks for a given square.
func AttacksBySquare(square int, occupied bitboard.Bitboard) bitboard.Bitboard {
	return rook.AttacksBySquare(square, occupied) | bishop.AttacksBySquare(square, occupied)
}
