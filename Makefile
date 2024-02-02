VERSION ?= dev
LD_FLAGS = -ldflags="-X 'github.com/shaardie/clemens/pkg/metadata.Version=$(VERSION)'"

.PHONY: clemens perft benchmark test clean

all: clemens perft

clemens: test
	GOOS=linux go build $(LD_FLAGS) -o clemens ./cmd/uci
	GOOS=windows go build $(LD_FLAGS) -o clemens.exe ./cmd/uci

perft: test
	GOOS=linux go build $(LD_FLAGS) -o perft ./cmd/perft
	GOOS=windows go build $(LD_FLAGS) -o perft.exe ./cmd/perft

benchmark:
	go test ./pkg/search -run=^$$ -bench ^BenchmarkSearch -cpuprofile profile.out

test:
	go test ./... -cover

clean:
	rm -f clemens clemens.exe perft perft.exe profile.out search.test
