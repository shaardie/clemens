package uci

import (
	"bufio"
	"os"

	"github.com/shaardie/clemens/pkg/uci/game"
)

var g game.Game

func Run() error {
	g = game.New()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		handleInput(scanner.Text())
	}
	return nil
}
