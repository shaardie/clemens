package position

import (
	"strconv"
	"strings"

	"github.com/shaardie/clemens/pkg/bitboard"
	. "github.com/shaardie/clemens/pkg/types"
)

const (
	pieceToChar string = " PNBRQK  pnbrqk"
	fileToChar  string = "ABCDEFGH"
)

// func NewFromFen(fen string) Position {
// 	reader := strings.NewReader(fen)
// 	pos := Position{}
// 	for {
// 		reader.ReadRune()
// 	}
// 	return pos
// }

func (pos Position) ToFen() string {
	sb := strings.Builder{}
	for rank := RANK_8; rank >= RANK_1; rank = rank - 1 {
		for file := FILE_A; file <= FILE_H; file++ {
			emptyCount := 0
			for file <= FILE_H && pos.Empty(bitboard.SquareFromRankAndFile(rank, file)) {
				emptyCount++
				file++
			}
			if emptyCount > 0 {
				sb.WriteString(strconv.Itoa(emptyCount))
			}

			if file <= FILE_H {
				sb.WriteByte(pieceToChar[pos.GetPieceFromSquare(bitboard.SquareFromRankAndFile(rank, file))])
			}
		}
		if rank >= RANK_1 {
			sb.WriteRune('/')
		}
	}
	if pos.SideToMove == WHITE {
		sb.WriteString(" w ")
	} else {
		sb.WriteString(" b ")
	}

	if !pos.CanCastle(ANY_CASTLING) {
		sb.WriteRune('-')
	} else {
		if pos.CanCastle(WHITE_CASTLING_KING) {
			sb.WriteRune('K')
		}
		if pos.CanCastle(WHITE_CASTLING_QUEEN) {
			sb.WriteRune('Q')
		}
		if pos.CanCastle(BLACK_CASTLING_KING) {
			sb.WriteRune('k')
		}
		if pos.CanCastle(BLACK_CASTLING_QUEEN) {
			sb.WriteRune('q')
		}
	}

	if pos.enPassante == SQUARE_NONE {
		sb.WriteString(" - ")
	} else {
		sb.WriteByte(fileToChar[bitboard.FileOfSquare(pos.enPassante)])
		sb.WriteString(strconv.Itoa(bitboard.RankOfSquare(pos.enPassante)))
	}

	sb.WriteString(strconv.Itoa(bitboard.RankOfSquare(pos.halfMoveClock)))
	sb.WriteRune(' ')
	sb.WriteString(strconv.Itoa(pos.numberOfFullMoves))
	return sb.String()
}
