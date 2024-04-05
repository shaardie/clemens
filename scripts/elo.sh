#!/bin/bash

set -eux
cutechess-cli -tournament gauntlet -concurrency 3 -pgnout output_pgn_file.pgn \
    -engine cmd=clemens st=0.1 \
    -engine name=maia-1900 cmd=lc0 arg=-w arg=/models/maia-1900.pb.gz nodes=1 st=1 \
    -each proto=uci -draw movenumber=40 movecount=4 score=8 \
    -resign movecount=4 score=500 \
    -openings file=/openings.pgn order=random -repeat -rounds 1200 -games 2
ordo -D -W -a 1900 -A "maia-1900" -p output_pgn_file.pgn
