#!/bin/bash

set -eux
ELO="2100"

cutechess-cli -tournament gauntlet -concurrency 3 -pgnout output_pgn_file.pgn \
    -engine cmd=clemens st=0.1 \
    -engine name=stockfish cmd=stockfish option.UCI_LimitStrength=true "option.UCI_Elo=$ELO" tc=40/1+0.05 \
    -each proto=uci -draw movenumber=40 movecount=4 score=8 \
    -resign movecount=4 score=500 \
    -openings file=/openings.pgn order=random -repeat -rounds 400 -games 2
ordo -D -W -a "$ELO" -A stockfish -p output_pgn_file.pgn
