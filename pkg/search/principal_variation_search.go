package search

import (
	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search/pvline"
)

func (s *Search) PrincipalVariationSearch(pos *position.Position, alpha, beta int, maxDepth, ply uint8, pvl *pvline.PVLine, canNull bool, alphaWasUpdated bool) (int, error) {
	var score int
	var err error
	if !alphaWasUpdated {
		score, err = s.negamax(pos, -beta, -alpha, maxDepth, ply+1, pvl, true)
		if err != nil {
			return 0, err
		}
	} else {
		score, err = s.negamax(pos, -alpha-1, -alpha, maxDepth, ply+1, &pvline.PVLine{}, true)
		if err != nil {
			return 0, err
		}
		// Rerun search
		if -score > alpha {
			score, err = s.negamax(pos, -beta, -alpha, maxDepth, ply+1, pvl, true)
			if err != nil {
				return 0, err
			}
		}
	}
	score = -score
	return score, nil
}
