package evaluation

import "math"

const (
	INF      int16 = math.MaxInt16
	maxPlies       = 100
)

func IsCheckmateValue(value int16) bool {
	if value < -INF+maxPlies || value > INF-maxPlies {
		return true
	}
	return false
}
