#!/bin/bash
set -euo pipefail

ENGINE="${ENGINE:-clemens-avx2}"
ANCHOR="${ANCHOR:-2300}"
BATCH="${BATCH:-200}"
TARGET="${TARGET:-15}"
MIN="${MIN:-400}"
MAX="${MAX:-6000}"

NAME="$("$ENGINE" <<<"uci" | awk '/^id name/ {sub(/^id name /, ""); print; exit}')"
STAMP="$(date +%Y-%m-%dT%H%M)"
PGN="/save/elo-${STAMP}.pgn"

games=0
while ((games < MAX)); do
  c-chess-cli \
    -each tc=40/1+0.05 \
    -engine cmd="$ENGINE" name="$NAME" \
    -engine name=stockfish cmd=stockfish option.UCI_LimitStrength=true "option.UCI_Elo=$ANCHOR" \
    -openings file=/openings/UHO_XXL_2022_+120_+149.epd order=random -repeat \
    -resign count=4 score=1000 \
    -draw number=40 count=8 score=10 \
    -games "$BATCH" -concurrency 4 \
    -pgn "$PGN"
  games=$((games + BATCH))

  ordo -W -D -s 1000 -a "$ANCHOR" -A stockfish -p "$PGN" -c "/save/elo-${STAMP}.csv" >/dev/null
  # rating and error of Clemens from ordo's CSV (column layout to be checked)
  read -r rating error < <(awk -F, -v n="$NAME" '$2 ~ n {gsub(/"/, ""); print $3, $4}' "/save/elo-${STAMP}.csv")
  echo "after ${games} games: ${rating} ± ${error}"

  ((games >= MIN)) && awk -v e="$error" -v t="$TARGET" 'BEGIN {exit !(e <= t)}' && break
done

printf '| %s | %s | %.0f ± %.0f | n=%d | SF18 @ %s | 40/1+0.05 | UHO_XXL_2022_+120_+149 |\n' \
  "$(date +%F)" "$NAME" "$rating" "$error" "$games" "$ANCHOR" | tee -a /save/elo-changelog.md
