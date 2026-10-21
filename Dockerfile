FROM golang:1.27.3-trixie

# Use Bash
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]

# Install stockfish and fastchess
WORKDIR /tmp
RUN <<EOF
curl -fsSL https://github.com/official-stockfish/Stockfish/releases/download/sf_19/stockfish-linux-x86-64-universal.tar.gz | tar xz
install -T stockfish/stockfish-linux-x86-64-universal /usr/bin/stockfish
rm -rf stockfish
EOF

RUN <<EOF
curl -fsSL https://github.com/Disservin/fastchess/releases/download/v1.8.2-alpha/fastchess-linux-x86-64.tar | tar x
install fastchess-linux-x86-64/fastchess /usr/bin
rm -rf fastchess-linux-x86-64
EOF

# Install current clemens engine
RUN mkdir /go/clemens
WORKDIR /go/clemens
COPY go.mod go.sum ./
RUN go mod download
COPY pkg pkg
COPY cmd cmd
COPY Makefile ./
ARG VERSION=dev
RUN VERSION=${VERSION} make clemens elo && install bin/clemens-* /usr/bin && install bin/elo /usr/bin

# Copy Openings
COPY openings /openings

# COPY Scripts
COPY scripts /scripts

RUN mkdir -p /save
WORKDIR /save
