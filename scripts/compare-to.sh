#!/bin/bash

set -eux

cp $1 /usr/bin/compare-to
cutechess-cli -tournament gauntlet -concurrency 3 -pgnout output_pgn_file.pgn \
    -engine cmd=clemens tc=40/1+0.05 \
    -engine cmd=compare-to name=compare-to tc=40/1+0.05 \
    -each proto=uci -draw movenumber=40 movecount=4 score=8 \
    -resign movecount=4 score=500 \
    -openings file=/openings.pgn order=random -repeat -rounds 500 -games 2
