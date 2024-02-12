#!/bin/bash

set -eux
cutechess-cli -tournament gauntlet -concurrency 4 -pgnout /save/output_pgn_file.pgn \
    -engine cmd=clemens st=1 \
    -engine name=maia-1100 cmd=lc0 arg=-w arg=/models/maia-1100.pb.gz nodes=1 st=1 \
    -engine name=maia-1200 cmd=lc0 arg=-w arg=/models/maia-1200.pb.gz nodes=1 st=1 \
    -engine name=maia-1300 cmd=lc0 arg=-w arg=/models/maia-1300.pb.gz nodes=1 st=1 \
    -engine name=maia-1400 cmd=lc0 arg=-w arg=/models/maia-1400.pb.gz nodes=1 st=1 \
    -engine name=maia-1500 cmd=lc0 arg=-w arg=/models/maia-1500.pb.gz nodes=1 st=1 \
    -engine name=maia-1600 cmd=lc0 arg=-w arg=/models/maia-1600.pb.gz nodes=1 st=1 \
    -engine name=maia-1700 cmd=lc0 arg=-w arg=/models/maia-1700.pb.gz nodes=1 st=1 \
    -engine name=maia-1800 cmd=lc0 arg=-w arg=/models/maia-1800.pb.gz nodes=1 st=1 \
    -engine name=maia-1900 cmd=lc0 arg=-w arg=/models/maia-1900.pb.gz nodes=1 st=1 \
    -each proto=uci -draw movenumber=40 movecount=4 score=8 \
    -resign movecount=4 score=500 \
    -openings file=/openings.pgn order=random -repeat -rounds 200 -games 2
ordo -Q -D -a 0 -A "maia-1100" -W -n8 -s1000 -U "0,1,2,3,4,5,6,7,8,9,10" -p /save/output_pgn_file.pgn
