package move

import (
	"fmt"
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

func (ml *MoveList) Get(idx uint8) Move {
	return ml.moves[idx]
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
