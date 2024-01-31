package game

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search"
	"github.com/shaardie/clemens/pkg/types"
	"github.com/shaardie/clemens/pkg/uci/state"
)

type Game interface {
	IsReady()
	NewPosition(tokens []string)
	StartSearch(tokens []string)
	StopSearch()
	private()
}

type goParameter struct {
	wtime     int
	btime     int
	winc      int
	binc      int
	movesToGo int
	depth     int
	movetime  int
	infinite  bool
}

const (
	infoChannelSize = 16
)

type gameImpl struct {
	isWorking    *sync.Mutex
	position     *position.Position
	state        state.State
	maxTimeInMs  int
	maxDepth     uint8
	info         chan search.Info
	search       *search.Search
	searchCancel context.CancelFunc
}

func New() Game {
	return newGameImpl()
}
func newGameImpl() *gameImpl {
	g := &gameImpl{
		isWorking:   &sync.Mutex{},
		state:       state.New(),
		maxTimeInMs: 5000,
		maxDepth:    6,
		info:        make(chan search.Info, infoChannelSize),
	}
	go g.showInfo()
	return g
}

func (g *gameImpl) private() {}

func (g *gameImpl) IsReady() {
	g.isWorking.Lock()
	defer g.isWorking.Unlock()
	fmt.Println("readyok")
}

func (g *gameImpl) NewPosition(tokens []string) {
	g.isWorking.Lock()
	defer g.isWorking.Unlock()

	if g.state.Get() == state.RUNNING {
		fmt.Println("info string wrong idle state to set new position")
		return
	}

	if len(tokens) == 0 {
		fmt.Println("info string no new position set")
		return
	}

	switch tokens[0] {
	case "startpos":
		g.position = position.New()
		g.state.Set(state.POSITION_SET)
		tokens = tokens[1:]
	case "fen":
		if len(tokens) < 7 {
			fmt.Println("info string fen string to short, no new position set")
			return
		}
		fenPos, err := position.NewFromFen(strings.Join(tokens[1:7], " "))
		if err != nil {
			fmt.Printf("info string broken fen string, %v\n", err)
			return
		}
		g.position = fenPos
		g.state.Set(state.POSITION_SET)
		tokens = tokens[7:]
	}
	if len(tokens) <= 1 || tokens[0] != "moves" {
		return
	}
	tokens = tokens[1:]

	for _, token := range tokens {
		err := g.position.MakeMoveFromString(token)
		if err != nil {
			fmt.Printf("info string error while making move %v, %v ", token, err)
			return
		}
	}
}

func parseGo(tokens []string) (gp goParameter) {
	var err error

	for len(tokens) > 0 {
		t := tokens[0]
		tokens = tokens[1:]
		switch t {
		case "searchmoves":
			fmt.Println("info string searchmoves not implemented")
			return gp
		case "wtime":
			if len(tokens) == 0 {
				fmt.Println("info string white time missing")
				return
			}
			gp.wtime, err = strconv.Atoi(tokens[0])
			if err != nil {
				fmt.Printf("info string white time broken, %v\n", err)
				return
			}
			tokens = tokens[1:]
		case "btime":
			if len(tokens) == 0 {
				fmt.Println("info string black time missing")
				return
			}
			gp.btime, err = strconv.Atoi(tokens[0])
			if err != nil {
				fmt.Printf("info string black time broken, %v\n", err)
				return
			}
			tokens = tokens[1:]
		case "winc":
			if len(tokens) == 0 {
				fmt.Println("info string white increment time missing")
				return
			}
			gp.winc, err = strconv.Atoi(tokens[0])
			if err != nil {
				fmt.Printf("info string white increment time broken, %v\n", err)
				return
			}
			tokens = tokens[1:]
		case "binc":
			if len(tokens) == 0 {
				fmt.Println("info string black increment time missing")
				return
			}
			gp.binc, err = strconv.Atoi(tokens[0])
			if err != nil {
				fmt.Printf("info string black increment time broken, %v\n", err)
				return
			}
			tokens = tokens[1:]
		case "movestogo":
			if len(tokens) == 0 {
				fmt.Println("info string moves to go missing")
				return
			}
			gp.movesToGo, err = strconv.Atoi(tokens[0])
			if err != nil {
				fmt.Printf("info string moves to go broken, %v\n", err)
				return
			}
			tokens = tokens[1:]
		case "depth":
			if len(tokens) == 0 {
				fmt.Println("info string depth missing")
				return
			}
			gp.depth, err = strconv.Atoi(tokens[0])
			if err != nil {
				fmt.Printf("info string depth broken, %v\n", err)
				return
			}
			tokens = tokens[1:]
		case "nodes":
			if len(tokens) == 0 {
				fmt.Println("info string nodes missing")
				return
			}
			_, err = strconv.Atoi(tokens[0])
			if err != nil {
				fmt.Printf("info string nodes broken, %v\n", err)
				return
			}
			tokens = tokens[1:]
			fmt.Println("info string nodes limit not implemented")
		case "mate":
			if len(tokens) == 0 {
				fmt.Println("info string mate missing")
				return
			}
			_, err = strconv.Atoi(tokens[0])
			if err != nil {
				fmt.Printf("info string mate broken, %v\n", err)
				return
			}
			tokens = tokens[1:]
			fmt.Println("info string mate not implemented")
		default:
			fmt.Printf("info string unknown go command %v\n", t)
			return
		}
	}
	return
}

func (g *gameImpl) createSearchContext(gp goParameter) (context.Context, context.CancelFunc) {
	// No need for any timeout
	if gp.infinite {
		return context.WithCancel(context.Background())
	}

	var t, inc int
	if g.position.SideToMove == types.BLACK {
		t = gp.btime
		inc = gp.binc
	} else {
		t = gp.wtime
		inc = gp.winc
	}

	var movetime int
	if gp.movetime > 0 {
		movetime = gp.movetime

	} else if t > 0 && gp.movesToGo > 0 {
		// calculate reasonable time, there is possibly a better way
		movetime = (t + inc*gp.movesToGo) / gp.movesToGo
		// do not calculate too long
		if movetime > g.maxTimeInMs {
			movetime = g.maxTimeInMs
		}
	} else {
		movetime = g.maxTimeInMs
	}
	return context.WithTimeout(context.Background(), time.Duration(movetime)*time.Millisecond)
}

func (g *gameImpl) StartSearch(tokens []string) {
	g.isWorking.Lock()
	defer g.isWorking.Unlock()

	if g.state.Get() != state.POSITION_SET || g.position == nil {
		fmt.Println("info string no position is set")
		return
	}

	// Create Search with the correct properties
	gp := parseGo(tokens)
	g.search = search.NewSearch(*g.position)
	ctx, cancel := g.createSearchContext(gp)
	var depth uint8
	if gp.depth > 0 {
		depth = uint8(gp.depth)
	} else {
		depth = g.maxDepth
	}

	g.searchCancel = cancel
	go func() {
		defer cancel()
		g.search.Search(ctx, depth, g.info)
		<-ctx.Done()
		if g.state.Get() == state.RUNNING {
			g.printBestMove()
		}
		g.state.Set(state.IDLE)
	}()
	g.state.Set(state.RUNNING)
}
func (g *gameImpl) StopSearch() {
	g.isWorking.Lock()
	defer g.isWorking.Unlock()

	if g.state.Get() != state.RUNNING {
		return
	}
	g.searchCancel()
	g.printBestMove()
	g.state.Set(state.IDLE)
}

func (g *gameImpl) printBestMove() {
	fmt.Printf("bestmove %v\n", g.search.BestMove())
}

func (g *gameImpl) showInfo() {
	for info := range g.info {
		// Only print while running
		// if g.state.Get() != state.RUNNING {
		// 	return
		// }
		infoString := fmt.Sprintf("info depth %v score cp %v nodes %v time %v", info.Depth, info.Score, info.Nodes, info.Time)
		pvString := info.PV.String()
		if pvString != "" {
			infoString = fmt.Sprintf("%v pv %v", infoString, pvString)
		}
		fmt.Println(infoString)
	}
}
