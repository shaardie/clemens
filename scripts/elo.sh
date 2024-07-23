#!/bin/bash

set -eux
ELO="1100"

c-chess-cli \
    -engine cmd=clemens tc=40/1+0.05 \
    -engine name=maia-1100 cmd="lc0 -w /models/maia-1100.pb.gz" nodes=1 movetime=1 \
    -openings file=/openings/UHO_XXL_2022_+120_+149.epd order=random -repeat \
    -resign count=4 score=1000 \
    -draw number=40 count=8 score=10 \
    -pgn output_pgn_file.pgn \
    -games 1000 -concurrency 4 -log

ordo -D -W -a "$ELO" -A maia-1100 -p output_pgn_file.pgn
