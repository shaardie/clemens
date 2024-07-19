#!/bin/bash

set -eux

c-chess-cli \
    -each depth=5 \
    -engine cmd=clemens \
    -engine cmd=clemens \
    -openings file=/openings/UHO_XXL_2022_+120_+149.epd order=random -repeat \
    -sample format=bin \
    -log \
    -games 400 -concurrency 3
