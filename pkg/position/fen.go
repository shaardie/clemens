package position

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/shaardie/clemens/pkg/bitboard"
	"github.com/shaardie/clemens/pkg/types"
)

const (
	pieceToChar string = " PNBRQK  pnbrqk"
	fileToChar  string = "ABCDEFGH"
)

func NewFromFen(fen string) (Position, error) {
	reader := strings.NewReader(fen)
	pos := Position{}
	// Pieces
	for {
		square := types.SQUARE_A1
		r, _, err := reader.ReadRune()
		if err != nil {
			return pos, fmt.Errorf("failed to read piece placement, %w", err)
		}
		if unicode.IsSpace(r) {
			break
		}
		if unicode.IsDigit(r) {
			square += int(r - '0')
		} else if r == '/' {
			square += 8
		} else if unicode.IsSpace(r) {
			break
		}
		p, err := types.PieceFromChar(r)
		if err != nil {
			return pos, fmt.Errorf("failed to get piece from rune, %w", err)
		}
		pos.SetPiece(p, square)
	}

	// Side to move
	r, _, err := reader.ReadRune()
	if err != nil {
		return pos, fmt.Errorf("failed to read side to move, %w", err)
	}
	switch r {
	case 'w':
		pos.SideToMove = types.WHITE
	case 'b':
		pos.SideToMove = types.BLACK
	default:
		return pos, fmt.Errorf("unknown color for side to move")
	}

	return pos, nil
}

func (pos Position) ToFen() string {
	sb := strings.Builder{}
	for rank := types.RANK_8; rank >= types.RANK_1; rank = rank - 1 {
		for file := types.FILE_A; file <= types.FILE_H; file++ {
			emptyCount := 0
			for file <= types.FILE_H && pos.Empty(bitboard.SquareFromRankAndFile(rank, file)) {
				emptyCount++
				file++
			}
			if emptyCount > 0 {
				sb.WriteString(strconv.Itoa(emptyCount))
			}

			if file <= types.FILE_H {
				sb.WriteRune(
					pos.GetPieceFromSquare(bitboard.SquareFromRankAndFile(rank, file)).ToChar(),
				)
			}
		}
		if rank >= types.RANK_1 {
			sb.WriteRune('/')
		}
	}
	if pos.SideToMove == types.WHITE {
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

	if pos.enPassante == types.SQUARE_NONE {
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
