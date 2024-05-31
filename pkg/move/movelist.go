package move

import (
	"fmt"
	"slices"
	"strings"
)

const (
	moveListSize = 255
)

type MoveList struct {
	moves [moveListSize]Move
	size  uint8
}

func NewMoveList() *MoveList {
	return &MoveList{}
}

func (ml *MoveList) Reset() {
	ml.size = 0
}

func (ml *MoveList) Get(idx uint8) *Move {
	return &ml.moves[idx]
}

func (ml *MoveList) Length() uint8 {
	return ml.size
}

func (ml *MoveList) Append(m Move) {
	ml.moves[ml.size] = m
	ml.size++
}

func (ml *MoveList) Set(idx uint8, m Move) {
	ml.moves[ml.size] = m
	if ml.size <= idx {
		ml.size = idx + 1
	}
}

func (ml *MoveList) String() string {
	ss := make([]string, ml.Length())
	for i := range ss {
		ss[i] = fmt.Sprint(ml.moves[i])
	}
	return strings.Join(ss, " ")
}

func (ml *MoveList) SortIndex(currIdx uint8) {
	// Get current Move
	currMove := ml.moves[currIdx]

	// Initialize with baseline
	score := currMove.GetScore()
	idx := currIdx

	// Search all candidates for a better move
	for candidateIdx := currIdx + 1; candidateIdx < ml.Length(); candidateIdx++ {
		candidateScore := ml.Get(candidateIdx).GetScore()
		if candidateScore > score {
			idx = candidateIdx
			score = candidateScore
		}
	}

	// Swap moves
	ml.moves[currIdx] = ml.moves[idx]
	ml.moves[idx] = currMove
}

// Sort moves, highest score first
func (ml *MoveList) Sort() {
	slices.SortFunc(ml.
		moves[:ml.size],
		func(a, b Move) int {
			return int(int(b.GetScore()) - int(a.GetScore()))
		},
	)
}
