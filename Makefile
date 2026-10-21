# Settings, can be overridden on the command line, e.g. `make FEATURES=baseline GOOS=windows`.
#
# VERSION   version shown in the UCI engine name, e.g. "0587bb1" or "0587bb1-dirty"
# GOOS      target operating system: linux or windows
# FEATURES  CPU features to build for: baseline, popcnt, avx2 or avx512.
#           The default avx2 needs a CPU with AVX2; use baseline for older CPUs.
VERSION ?= $(shell git describe --always --dirty 2>/dev/null || echo dev)
GOOS ?= linux
FEATURES ?= avx2

# Output paths. Windows binaries get the .exe ending.
EXECUTABLE_ENDING_linux =
EXECUTABLE_ENDING_windows = .exe
CLEMENS_BIN = bin/clemens-$(FEATURES)$(EXECUTABLE_ENDING_$(GOOS))
PERFT_BIN = bin/perft-$(FEATURES)$(EXECUTABLE_ENDING_$(GOOS))
ELO_BIN = bin/elo$(EXECUTABLE_ENDING_$(GOOS))

# Go's GOAMD64 level for each value of FEATURES.
GOAMD64_baseline = v1
GOAMD64_popcnt   = v2
GOAMD64_avx2     = v3
GOAMD64_avx512   = v4

# Embed version and features into pkg/metadata.
BUILD_LD_FLAGS = -ldflags="\
  -X 'github.com/shaardie/clemens/pkg/metadata.Version=$(VERSION)' \
  -X 'github.com/shaardie/clemens/pkg/metadata.Features=$(FEATURES)'"

# User and group the measurement containers run as, so that files in save/ belong to the
# caller and not to root.
IMAGE_USER := $(shell id -u)
IMAGE_GROUP := $(shell id -g)

.PHONY: build clemens perft elo release benchmark benchmark_perft benchmark_search image test-elo compare-to test clean check-compare-to check-features

# --- Build -----------------------------------------------------------------------------------

# Build the engine, perft and the Elo tool into bin/ for one FEATURES/GOOS combination.
build: clemens perft elo

# The chess engine, e.g. bin/clemens-avx2.
clemens: check-features
	GOOS=$(GOOS) GOAMD64=$(GOAMD64_$(FEATURES)) go build $(BUILD_LD_FLAGS) -o $(CLEMENS_BIN) ./cmd/uci

# Perft tool to verify the move generation, e.g. bin/perft-avx2.
perft: check-features
	GOOS=$(GOOS) GOAMD64=$(GOAMD64_$(FEATURES)) go build $(BUILD_LD_FLAGS) -o $(PERFT_BIN) ./cmd/perft

# Elo measurement tool (cmd/elo). It does not depend on FEATURES and runs in the container.
elo:
	GOOS=$(GOOS) go build $(BUILD_LD_FLAGS) -o $(ELO_BIN) ./cmd/elo

# Run the tests and build the engine for all FEATURES, for Linux and Windows.
release: test
	$(MAKE) clemens FEATURES=baseline GOOS=linux
	$(MAKE) clemens FEATURES=popcnt GOOS=linux
	$(MAKE) clemens FEATURES=avx2 GOOS=linux
	$(MAKE) clemens FEATURES=avx512 GOOS=linux
	$(MAKE) clemens FEATURES=baseline GOOS=windows
	$(MAKE) clemens FEATURES=popcnt GOOS=windows
	$(MAKE) clemens FEATURES=avx2 GOOS=windows
	$(MAKE) clemens FEATURES=avx512 GOOS=windows

# --- Tests and benchmarks --------------------------------------------------------------------

# Run all tests.
test: check-features
	GOAMD64=$(GOAMD64_$(FEATURES)) go test ./... -cover

# Run the benchmarks. CPU profiles and test binaries are written to profile/, e.g. for
# `go tool pprof profile/search.test profile/profile_search.out`.
benchmark: benchmark_perft benchmark_search

benchmark_search: check-features profile
	GOAMD64=$(GOAMD64_$(FEATURES)) go test ./pkg/search -run=^$$ -bench ^BenchmarkSearch -cpuprofile profile/profile_search.out -o profile/search.test

benchmark_perft: check-features profile
	GOAMD64=$(GOAMD64_$(FEATURES)) go test ./cmd/perft -run=^$$ -bench ^BenchmarkPerft -cpuprofile profile/profile_perft.out -o profile/perft.test

# --- Measurements in Docker ------------------------------------------------------------------

# Docker image with the current engine, fastchess, Stockfish and the openings.
image:
	docker build . -t clemens --build-arg VERSION=$(VERSION)

# Measure the Elo of the current engine against Stockfish. PGNs are written to save/ and the
# result is printed as a row for the changelog in README.md.
test-elo: image save
	docker run --rm \
		--user=$(IMAGE_USER):$(IMAGE_GROUP) \
		--mount type=bind,source=$(PWD)/save,target=/save \
		clemens:latest elo

# SPRT of the current engine against another engine binary, e.g.
# `make compare-to COMPARE_TO=engines/clemens-0587bb1`. The binary has to run on Linux.
compare-to: check-compare-to image save
	docker run --rm \
		--user=$(IMAGE_USER):$(IMAGE_GROUP) \
		--mount type=bind,source=$(abspath $(COMPARE_TO)),target=/usr/bin/compare-to,readonly \
		--mount type=bind,source=$(PWD)/save,target=/save \
		clemens:latest /scripts/compare-to.sh

# --- Helpers ---------------------------------------------------------------------------------

# Stop make with a message if a setting is missing or invalid.
check-features:
	$(if $(GOAMD64_$(FEATURES)),,$(error unknown FEATURES '$(FEATURES)'))

check-compare-to:
	$(if $(COMPARE_TO),,$(error COMPARE_TO is missing))

# Output directories.
profile:
	mkdir -p profile

save:
	mkdir -p save

# Remove build output and profiles. Measurement results in save/ are kept.
clean:
	rm -rf bin profile
