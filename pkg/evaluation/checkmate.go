package evaluation

const (
	INF      = 100000
	maxPlies = 100
)

func IsCheckmateValue(value int) bool {
	if value < -INF+maxPlies || value > INF-maxPlies {
		return true
	}
	return false
}
