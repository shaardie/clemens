package nnue

import "github.com/shaardie/clemens/pkg/types"

const squareNumber = int(types.SQUARE_NUMBER)

// FeatureIndex calculates the feature index for a piece on a square
func FeatureIndex(piece types.Piece, square uint8) (whiteFeatureIndex int, blackFeatureIndex int) {
	isWhitePiece := piece.Color() == types.WHITE
	pieceType := int(piece.Type())

	// white perspective: own pieces 0-5, opponent pieces 6-11
	wPt := int(pieceType)
	if !isWhitePiece {
		wPt += 6
	}
	wIdx := wPt*squareNumber + int(square)

	// white perspective: own pieces 0-5, opponent pieces 6-11
	// Schwarze Perspektive: eigene Figuren 0–5, gegnerische 6–11
	// Flipping square, see https://www.chessprogramming.org/Color_Flipping#Flipping_an_8x8_Board
	bPt := int(pieceType)
	if isWhitePiece {
		bPt += 6
	}
	bIdx := bPt*squareNumber + int(square^56)

	return wIdx, bIdx
}
