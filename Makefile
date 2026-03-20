VERSION ?= dev
COMPARE_TO ?= $(PWD)/clemens-avx2

# GOAMD64 target level (v1=baseline, v2=SSE4.2, v3=AVX2, v4=AVX-512)
# Can be overridden: make clemens GOAMD64=v3
GOAMD64 ?= v1

# Human-readable feature names for each level
GOAMD64_FEATURES_v1 = baseline
GOAMD64_FEATURES_v2 = popcnt
GOAMD64_FEATURES_v3 = avx2
GOAMD64_FEATURES_v4 = avx512

FEATURES = $(GOAMD64_FEATURES_$(GOAMD64))

BUILD_LD_FLAGS = -ldflags="\
  -X 'github.com/shaardie/clemens/pkg/metadata.Version=$(VERSION)' \
  -X 'github.com/shaardie/clemens/pkg/metadata.Features=$(FEATURES)'"

.PHONY: clemens perft benchmark test clean all-variants

all: all-variants

clemens:
	GOOS=linux GOAMD64=$(GOAMD64) go build $(BUILD_LD_FLAGS) -o clemens-$(FEATURES) ./cmd/uci
	GOOS=windows GOAMD64=$(GOAMD64) go build $(BUILD_LD_FLAGS) -o clemens-$(FEATURES).exe ./cmd/uci

perft: test
	GOOS=linux GOAMD64=$(GOAMD64) go build $(BUILD_LD_FLAGS) -o perft-$(FEATURES) ./cmd/perft
	GOOS=windows GOAMD64=$(GOAMD64) go build $(BUILD_LD_FLAGS) -o perft-$(FEATURES).exe ./cmd/perft

all-variants:
	make clemens GOAMD64=v1
	make clemens GOAMD64=v2
	make clemens GOAMD64=v3
	make clemens GOAMD64=v4

benchmark: benchmark_perft benchmark_search

benchmark_search:
	GOAMD64=$(GOAMD64) go test ./pkg/search -run=^$$ -bench ^BenchmarkSearch -cpuprofile profile_search.out
benchmark_perft:
	GOAMD64=$(GOAMD64) go test ./cmd/perft -run=^$$ -bench ^BenchmarkPerft -cpuprofile profile_perft.out

elo:
	docker build . -t elo && mkdir -p $(PWD)/save && docker run --rm -v $(PWD)/save:/save elo:latest /scripts/elo.sh

compare-to:
	docker build . -t elo && docker run --rm -v $(COMPARE_TO):/compare-to elo:latest /scripts/compare-to.sh /compare-to

test:
	GOAMD64=$(GOAMD64) go test ./... -cover

clean:
	rm -rf clemens-* perft-* profile*.out search.test save
