package position

import (
	"reflect"
	"testing"

	"github.com/shaardie/clemens/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewFromFen(t *testing.T) {
	pos := New()
	fenPos, err := NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	assert.NoError(t, err)
	if !reflect.DeepEqual(pos, fenPos) {
		t.Errorf("NewFromFen() = %v, want %v", fenPos, pos)
	}
}

func TestPosition_fenSetPieces(t *testing.T) {
	startPos := New()
	anotherPos := &Position{}
	anotherPos.PiecesBoard[types.SQUARE_H1] = types.WHITE_ROOK
	anotherPos.PiecesBoard[types.SQUARE_E2] = types.WHITE_KING
	anotherPos.PiecesBoard[types.SQUARE_E7] = types.BLACK_KNIGHT
	anotherPos.PiecesBoard[types.SQUARE_C8] = types.BLACK_KING

	anotherPos.boardToBitBoard()
	tests := []struct {
		name      string
		token     string
		wantedPos *Position
		wantErr   bool
	}{
		{
			name:  "beginning",
			token: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
			wantedPos: &Position{
				PiecesBoard:    startPos.PiecesBoard,
				PiecesBitboard: startPos.PiecesBitboard,
			},
			wantErr: false,
		},
		{
			name:      "another example",
			token:     "2k5/4n3/8/8/8/8/4K3/7R",
			wantedPos: anotherPos,
			wantErr:   false,
		},
		{
			name:      "broken token",
			token:     "banana",
			wantedPos: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := &Position{}
			err := pos.fenSetPieces(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			}
			if tt.wantedPos != nil {
				assert.Equal(t, pos, tt.wantedPos)
			}
		})
	}
}

func TestPosition_fenSetSideToMove(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		wantedPos *Position
		wantErr   bool
	}{
		{
			name:  "black",
			token: "b",
			wantedPos: &Position{
				SideToMove: types.BLACK,
			},
			wantErr: false,
		},
		{
			name:  "white",
			token: "w",
			wantedPos: &Position{
				SideToMove: types.WHITE,
			},
			wantErr: false,
		},
		{
			name:      "broken token",
			token:     "s",
			wantedPos: nil,
			wantErr:   true,
		},
		{
			name:      "token to long",
			token:     "banana",
			wantedPos: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := &Position{}
			err := pos.fenSetSideToMove(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			}
			if tt.wantedPos != nil {
				assert.Equal(t, pos, tt.wantedPos)
			}
		})
	}
}

func TestPosition_fenSetCastling(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		wantedPos *Position
		wantErr   bool
	}{
		{
			name:  "no castling",
			token: "-",
			wantedPos: &Position{
				castling: NO_CASTLING,
			},
			wantErr: false,
		},
		{
			name:  "some castling",
			token: "Kk",
			wantedPos: &Position{
				castling: WHITE_CASTLING_KING | BLACK_CASTLING_KING,
			},
			wantErr: false,
		},
		{
			name:  "all castling",
			token: "KQkq",
			wantedPos: &Position{
				castling: ANY_CASTLING,
			},
			wantErr: false,
		},
		{
			name:      "broken token",
			token:     "banana",
			wantedPos: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := &Position{}
			err := pos.fenSetCastling(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			}
			if tt.wantedPos != nil {
				assert.Equal(t, pos, tt.wantedPos)
			}
		})
	}
}

func TestPosition_fenSetEnPassante(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		wantedPos *Position
		wantErr   bool
	}{
		{
			name:  "no en passante",
			token: "-",
			wantedPos: &Position{
				enPassante: types.SQUARE_NONE,
			},
			wantErr: false,
		},
		{
			name:  "en passante d6",
			token: "d6",
			wantedPos: &Position{
				enPassante: types.SQUARE_D6,
			},
			wantErr: false,
		},
		{
			name:      "broken token",
			token:     "banana",
			wantedPos: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := &Position{}
			err := pos.fenSetEnPassante(tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			}
			if tt.wantedPos != nil {
				assert.Equal(t, pos, tt.wantedPos)
			}
		})
	}
}

func TestPosition_ToFen(t *testing.T) {
	assert.Equal(
		t,
		New().ToFen(),
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	)
}
