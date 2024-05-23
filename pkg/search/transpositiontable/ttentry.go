package transpositiontable

import "github.com/shaardie/clemens/pkg/move"

type ttEntry struct {
	zobristHash uint64
	bestMove    move.Move
	score       int16
	depth       uint8

	// 0-1 for the Node Type
	// 2-7 for the Age
	ageAndNodeType uint8
}

func (te *ttEntry) getNodeType() nodeType {
	return nodeType(te.ageAndNodeType & 0b11)
}
func (te *ttEntry) setNodeType(nt nodeType) {
	te.ageAndNodeType = te.ageAndNodeType&0b11111100 | uint8(nt)
}
func (te *ttEntry) getAge() uint8 {
	return te.ageAndNodeType >> 2
}
func (te *ttEntry) setAge(age uint8) {
	te.ageAndNodeType = te.ageAndNodeType&0b11 | (age << 2)
}
