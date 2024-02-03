package game

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/shaardie/clemens/pkg/position"
	"github.com/shaardie/clemens/pkg/search"
	"github.com/shaardie/clemens/pkg/uci/state"
)

type Game interface {
	IsReady()
	NewPosition(tokens []string)
	StartSearch(tokens []string)
	StopSearch()
	private()
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
			fmt.Printf("info string error while making move %v, %v \n", token, err)
			return
		}
	}
}

func parseGo(tokens []string) (sp search.SearchParameter) {
	var err error

	for len(tokens) > 0 {
		t := tokens[0]
		tokens = tokens[1:]
		switch t {
		case "searchmoves":
			fmt.Println("info string searchmoves not implemented")
			return sp
		case "wtime":
			if len(tokens) == 0 {
				fmt.Println("info string white time missing")
				return
			}
			sp.WTime, err = strconv.Atoi(tokens[0])
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
			sp.BTime, err = strconv.Atoi(tokens[0])
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
			sp.WInc, err = strconv.Atoi(tokens[0])
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
			sp.BInc, err = strconv.Atoi(tokens[0])
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
			sp.MovesToGo, err = strconv.Atoi(tokens[0])
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
			d, err := strconv.Atoi(tokens[0])
			if err != nil {
				fmt.Printf("info string depth broken, %v\n", err)
				return
			}
			sp.Depth = uint8(d)
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
		case "infinite":
			sp.Infinite = true
			tokens = tokens[1:]
		default:
			fmt.Printf("info string unknown go command %v\n", t)
			return
		}
	}
	return
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
	ctx, cancel := context.WithCancel(context.Background())
	g.searchCancel = cancel
	go func() {
		defer cancel()
		g.search.Search(ctx, gp)
		fmt.Printf("bestmove %v\n", g.search.BestMove())
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
