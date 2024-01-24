package uci

import (
	"bufio"
	"fmt"
	"os"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search"
)

var pos *position.Position

// This is not good
func run(currPos position.Position) {
	m := search.Search(&currPos, 1)
	fmt.Printf("bestmove %v\n", m)
}

func Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		handleInput(scanner.Text())
	}
	return nil
}
