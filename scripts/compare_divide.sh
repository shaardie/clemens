#!/bin/bash

set -eu

FEN_STRING="$1"
DEPTH="$2"

CLEMENS_OUTPUT="$(go run cmd/perft/perft.go -depth "$DEPTH" -position "$FEN_STRING" -divide | grep ^[a-h] | sort)"

STOCKFISH_OUTPUT="$(echo "position fen $FEN_STRING
go perft $DEPTH
quit
" | stockfish | grep ^[a-h] | sort)"

colordiff -u <(echo "$CLEMENS_OUTPUT") <(echo "$STOCKFISH_OUTPUT")
