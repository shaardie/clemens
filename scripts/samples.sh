#!/bin/bash

set -eux

c-chess-cli \
    -each depth=5 \
    -engine name=stockfish-1 cmd=stockfish \
    -engine name=stockfish-2 cmd=stockfish \
    -openings file=/openings/UHO_XXL_2022_+120_+149.epd order=random -repeat \
    -sample format=bin \
    -log \
    -games 100000 -concurrency 4
