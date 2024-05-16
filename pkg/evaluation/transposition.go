package evaluation

import (
	"unsafe"
)

type transpositionEntry struct {
	zobristHash uint64
	score       int16
}

// 16MB
const transpositionTableSizeinMB = 1024 * 1024 * 16

var transpositionTableSize uint64
var tTable transpositionTable

type transpositionTable []transpositionEntry

func init() {
	reset()
}

func reset() {
	transpositionTableSize = uint64(transpositionTableSizeinMB / unsafe.Sizeof(transpositionEntry{}))
	tTable = make([]transpositionEntry, transpositionTableSize)
}

func (tt transpositionTable) get(zobristHash uint64) (int16, bool) {
	key := zobristHash % transpositionTableSize
	te := tt[key]
	if tt[key].zobristHash != zobristHash {
		return tt[key].score, false
	}
	return te.score, true
}

// save save the new transposition entry.
func (tt transpositionTable) save(zobristHash uint64, score int16) {
	key := zobristHash % transpositionTableSize
	tt[key].zobristHash = zobristHash
	tt[key].score = score
}
