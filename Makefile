VERSION ?= dev
LD_FLAGS = -ldflags="-X 'github.com/shaardie/clemens/pkg/metadata.Version=$(VERSION)'"
COMPARE_TO ?= $(PWD)/clemens

.PHONY: clemens perft benchmark test clean

all: clemens perft

clemens: test
	GOOS=linux go build $(LD_FLAGS) -o clemens ./cmd/uci
	GOOS=windows go build $(LD_FLAGS) -o clemens.exe ./cmd/uci

perft: test
	GOOS=linux go build $(LD_FLAGS) -o perft ./cmd/perft
	GOOS=windows go build $(LD_FLAGS) -o perft.exe ./cmd/perft

benchmark: benchmark_perft benchmark_search

benchmark_search:
	go test ./pkg/search -run=^$$ -bench ^BenchmarkSearch -cpuprofile profile_search.out
benchmark_perft:
	go test ./cmd/perft -run=^$$ -bench ^BenchmarkPerft -cpuprofile profile_perft.out

elo:
	docker build . -t elo && mkdir -p $(PWD)/save && docker run --rm -v $(PWD)/save:/save elo:latest

compare-to:
	docker build . -t elo && docker run --rm -v $(COMPARE_TO):/compare-to elo:latest /compare-to.sh /compare-to

test:
	go test ./... -cover

clean:
	rm -rf clemens clemens.exe perft perft.exe profile.out search.test save
