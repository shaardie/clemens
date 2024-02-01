#!/bin/bash

set -eux

cutechess-cli -tournament gauntlet -concurrency 1 -pgnout output_pgn_file.pgn \
    -engine cmd=clemens tc=40/60+1 \
    -engine cmd=stockfish tc=40/60+1 \
    -engine cmd=blunder tc=40/60+1 \
    -each proto=uci \
    -draw movenumber=40 movecount=4 score=8 \
    -resign movecount=4 score=500 \
    -openings file=/openings.pgn policy=round -repeat -rounds 100 -games 2 -debug
ordo -Q -D -a 0 -A "Stockfish 16" -W -n8 -s1000 -U "0,1,2,3,4,5,6,7,8,9,10" -p output_pgn_file.pgn
