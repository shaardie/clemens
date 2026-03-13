VERSION ?= dev
LD_FLAGS = -ldflags="-X 'github.com/shaardie/clemens/pkg/metadata.Version=$(VERSION)'"
COMPARE_TO ?= $(PWD)/clemens

# Set SIMD=1 to enable AVX2/FMA SIMD support (requires Go 1.26+, amd64)
# Example: make clemens SIMD=1
ifdef SIMD
  GOEXPERIMENT_FLAG = GOEXPERIMENT=simd
  SIMD_SUFFIX = -simd
else
  GOEXPERIMENT_FLAG =
  SIMD_SUFFIX =
endif

.PHONY: clemens perft benchmark test clean

all: clemens perft

clemens: test
	GOOS=linux $(GOEXPERIMENT_FLAG) go build $(LD_FLAGS) -o clemens$(SIMD_SUFFIX) ./cmd/uci
	GOOS=windows $(GOEXPERIMENT_FLAG) go build $(LD_FLAGS) -o clemens$(SIMD_SUFFIX).exe ./cmd/uci

perft: test
	GOOS=linux $(GOEXPERIMENT_FLAG) go build $(LD_FLAGS) -o perft$(SIMD_SUFFIX) ./cmd/perft
	GOOS=windows $(GOEXPERIMENT_FLAG) go build $(LD_FLAGS) -o perft$(SIMD_SUFFIX).exe ./cmd/perft

benchmark: benchmark_perft benchmark_search

benchmark_search:
	$(GOEXPERIMENT_FLAG) go test ./pkg/search -run=^$$ -bench ^BenchmarkSearch -cpuprofile profile_search.out
benchmark_perft:
	$(GOEXPERIMENT_FLAG) go test ./cmd/perft -run=^$$ -bench ^BenchmarkPerft -cpuprofile profile_perft.out

elo:
	docker build . -t elo && mkdir -p $(PWD)/save && docker run --rm -v $(PWD)/save:/save elo:latest /scripts/elo.sh

compare-to:
	docker build . -t elo && docker run --rm -v $(COMPARE_TO):/compare-to elo:latest /scripts/compare-to.sh /compare-to

test:
	$(GOEXPERIMENT_FLAG) go test ./... -cover

clean:
	rm -rf clemens clemens-simd clemens.exe clemens-simd.exe perft perft-simd perft.exe perft-simd.exe profile.out search.test save
