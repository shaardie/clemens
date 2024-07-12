#!/bin/bash

set -eux
ELO="2300"

c-chess-cli \
    -each tc=40/1+0.05 \
    -engine cmd=clemens \
    -engine name=stockfish cmd=stockfish option.UCI_LimitStrength=true "option.UCI_Elo=$ELO" \
    -openings file=/openings/UHO_XXL_2022_+120_+149.epd order=random -repeat \
    -resign count=4 score=1000 \
    -draw number=40 count=8 score=10 \
    -sprt -pgn output_pgn_file.pgn \
    -games 400 -concurrency 3

ordo -D -W -a "$ELO" -A stockfish -p output_pgn_file.pgn
