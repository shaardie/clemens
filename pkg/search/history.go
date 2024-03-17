package search

import "github.com/shaardie/clemens/pkg/position"

func (s *Search) pushHistory(pos *position.Position) {
	s.searchHistory[s.searchHistoryPly] = pos.ZobristHash
	s.searchHistoryPly++
}

func (s *Search) popHistory() {
	s.searchHistoryPly--
	s.searchHistory[s.searchHistoryPly] = 0
}

func (s *Search) isRepetition(pos *position.Position) bool {
	for i := 0; i < s.searchHistoryPly; i++ {
		if s.searchHistory[i] == pos.ZobristHash {
			return true
		}
	}
	return false
}

func (s *Search) MakeMoveFromString(m string) error {
	err := s.Pos.MakeMoveFromString(m)
	if err != nil {
		return err
	}
	s.pushHistory(&s.Pos)
	return nil
}
