#!/bin/bash

set -eux

cutechess-cli -tournament gauntlet -concurrency 3 -pgnout output_pgn_file.pgn \
    -engine cmd=clemens name=clemens-1 st=0.08 \
    -engine cmd=clemens name=clemens-2 st=0.08 \
    -each proto=uci -draw movenumber=40 movecount=4 score=8 \
    -openings file=/openings.pgn order=random -repeat -rounds 10000 -games 2
