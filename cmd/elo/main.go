package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	eloLinePrefix = "Elo: "
	minGames      = 400
)

var (
	engine       = flag.String("engine", "clemens-avx2", "engine to test")
	estimatedElo = flag.Float64("estimated-elo", 2300, "elo to test against. The nearer this is on the actual elo, the faster the tests will converge.")
	accuracy     = flag.Float64("accuracy", 15, "target error in Elo: the test stops once the half-width of the 95% confidence interval reported by fastchess (the \"+/-\" value) is at or below this value.")

	eloLineRegex = regexp.MustCompile(`([\-\.0-9]+) \+\/\- ([\.0-9]+)`)
)

type result struct {
	time     time.Time
	gameType string
	games    int
	eloDiff  float64
	accuracy float64
	engine   string
	opponent string
}

func (r result) toMarkdown() string {
	return fmt.Sprintf("| %v | %v | %0.f ± %0.f | %v | %v | %v |", r.time.Format(time.DateOnly), r.engine, *estimatedElo+r.eloDiff, r.accuracy, r.games, r.opponent, r.gameType)
}

func getName() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, *engine)
	cmd.Stdin = strings.NewReader("uci\n")
	out := strings.Builder{}
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run command, %w", err)
	}

	const prefix = "id name "
	var name string
	for line := range strings.Lines(out.String()) {
		if s, ok := strings.CutPrefix(line, prefix); ok {
			name = strings.TrimSpace(s)
			break
		}
	}
	if name == "" {
		return "", fmt.Errorf("name for engine %v is empty", *engine)
	}

	return name, nil
}

func runEloTest(name string) (result, error) {
	r := result{
		engine:   name,
		opponent: fmt.Sprintf("stockfish %v", *estimatedElo),
		time:     time.Now(),
		gameType: "40/1+0.05",
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(
		ctx, "fastchess",
		"-concurrency", "4",
		"-rounds", "6000",
		"-each", fmt.Sprintf("tc=%v", r.gameType),
		"-engine", fmt.Sprintf("cmd=%v", *engine), fmt.Sprintf("name=%v", r.engine),
		"-engine", "cmd=stockfish", fmt.Sprintf("name=%v", r.opponent), "option.UCI_LimitStrength=true", fmt.Sprintf("option.UCI_Elo=%v", *estimatedElo),
		"-draw", "movenumber=40", "movecount=8", "score=10",
		"-resign", "movecount=4", "score=1000",
		"-openings", "file=/openings/UHO_XXL_2022_+120_+149.epd", "order=random",
		"-pgnout", fmt.Sprintf("file=elo-%v.pgn", r.time.Unix()),
	)
	cmd.Stderr = os.Stderr
	cmd.WaitDelay = 10 * time.Second
	cmd.Cancel = func() error {
		return cmd.Process.Signal(syscall.SIGINT)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return r, fmt.Errorf("failed to get stdout, %w", err)
	}

	if err := cmd.Start(); err != nil {
		return r, fmt.Errorf("failed to start command, %w", err)
	}

	success := false
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, eloLinePrefix) {
			continue
		}

		r.games += 20 // Elo Line is shown every 20 games
		fmt.Println(line)

		// Parse the line
		matches := eloLineRegex.FindStringSubmatch(line)
		if len(matches) != 3 {
			log.Println("Warning! Failed to parse line. Continue and try again")
			continue
		}
		r.eloDiff, err = strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return r, fmt.Errorf("failed to parse elo difference, %w", err)
		}
		r.accuracy, err = strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return r, fmt.Errorf("failed to parse elo accuracy, %w", err)
		}
		if r.accuracy <= *accuracy && r.games > minGames {
			success = true
			cancel()
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return r, fmt.Errorf("reading fastchess output, %w", err)
	}

	if err := cmd.Wait(); err != nil {
		if errors.Is(err, exec.ErrWaitDelay) {
			log.Println("Warning! Failed to stop fastchess in time. Maybe the engines are still running.")
		}

		// If the calculation was successful and fastchess was stopped, an error is expected
		if !success {
			return r, fmt.Errorf("failed to wait for fastchess to terminate, %w", err)
		}
	}

	if !success {
		return r, errors.New("failed to calculate elo")
	}

	return r, nil
}

func main() {
	flag.Parse()

	name, err := getName()
	if err != nil {
		log.Fatalf("failed to get name of the engine %v, %v", *engine, err)
	}

	r, err := runEloTest(name)
	if err != nil {
		log.Fatalf("failed to run elo test, %v", err)
	}
	fmt.Println(r.toMarkdown())
}
