#!/bin/bash

set -euo pipefail

ENGINE="${ENGINE:-clemens-avx2}"
NAME="$("$ENGINE" <<<"uci" | awk '/^id name/ {sub(/^id name /, ""); print; exit}')"

cp $1 /usr/bin/compare-to

fastchess \
  -each tc=40/1+0.05 \
  -engine cmd="$ENGINE" name="$NAME" \
  -engine cmd=compare-to name=compare-to \
  -openings file=/openings/UHO_XXL_2022_+120_+149.epd order=random \
  -draw movenumber=40 movecount=8 score=10 \
  -resign movecount=4 score=1000 \
  -sprt elo0=0 elo1=5 alpha=0.05 beta=0.05 model=logistic \
  -rounds 5000 \
  -concurrency 3 2>&1 | grep --line-buffered -E '^(LLR|Finished game|Results of|Elo:)'
