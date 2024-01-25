package uci

import (
	"fmt"
	"strings"
	"sync"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search"
)

type game struct {
	isReadyM *sync.Mutex
	position *position.Position
	maxDepth int

	result  search.SearchResult
	resultM *sync.Mutex
}

var s game

func (s *game) init() {
	s.isReadyM = &sync.Mutex{}
	s.resultM = &sync.Mutex{}
	s.maxDepth = 5
}

func (s *game) isReady() {
	s.isReadyM.Lock()
	defer s.isReadyM.Unlock()
	fmt.Println("readyok")
}

func (s *game) newGame() {
	s.isReadyM.Lock()
	defer s.isReadyM.Unlock()
	s.position = position.New()
}

func (s *game) newPosition(tokens []string) {
	s.isReadyM.Lock()
	defer s.isReadyM.Unlock()

	if len(tokens) == 0 {
		return
	}

	switch tokens[0] {
	case "startpos":
		s.position = position.New()
		tokens = tokens[1:]
	case "fen":
		if len(tokens) < 7 {
			return
		}
		fenPos, err := position.NewFromFen(strings.Join(tokens[1:7], " "))
		if err != nil {
			return
		}
		s.position = fenPos
		tokens = tokens[7:]
	}
	if len(tokens) <= 1 || tokens[0] != "moves" {
		return
	}
	tokens = tokens[1:]

	for _, token := range tokens {
		err := s.position.MakeMoveFromString(token)
		if err != nil {
			return
		}
	}
}

func (s *game) findBestMove(tokens []string) {
	s.isReadyM.Lock()
	defer s.isReadyM.Unlock()

	// Ignore Command if no position is set
	if s.position == nil {
		return
	}

	for depth := 1; depth <= s.maxDepth; depth++ {
		currPos := *s.position
		r := search.Search(&currPos, depth)
		fmt.Printf("info depth %d score cp %v nodes %v \n", depth, r.Score, r.Nodes)
		s.setResult(r)
	}

	fmt.Printf("bestmove %v\n", s.getResult().Move)
}

func (s *game) setResult(r search.SearchResult) {
	s.resultM.Lock()
	s.result = r
	s.resultM.Unlock()
}

func (s *game) getResult() search.SearchResult {
	s.resultM.Lock()
	defer s.resultM.Unlock()
	return s.result

}
