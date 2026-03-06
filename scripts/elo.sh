#!/bin/bash

set -eux
ELO="1400"

c-chess-cli \
  -each tc=40/1+0.05 \
  -engine cmd=clemens \
  -engine name=maia-$ELO cmd="lc0 -w /models/maia-$ELO.pb.gz" nodes=1 movetime=1 \
  -openings file=/openings/UHO_XXL_2022_+120_+149.epd order=random -repeat \
  -resign count=4 score=1000 \
  -draw number=40 count=8 score=10 \
  -sprt -pgn output_pgn_file.pgn \
  -games 400 -concurrency 3

ordo -D -W -a "$ELO" -A maia-1600 -p output_pgn_file.pgn
