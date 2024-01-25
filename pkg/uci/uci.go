package uci

import (
	"bufio"
	"os"

	"github.com/shaardie/clemens/pkg/position"
)

var pos *position.Position

func init() {
	s.init()
}

func Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		handleInput(scanner.Text())
	}
	return nil
}
