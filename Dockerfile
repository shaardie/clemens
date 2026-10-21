FROM golang:1.27.2-trixie

RUN apt-get update && apt-get install -y curl libqt5core5a unzip meson wget zlib1g-dev

# Install models
RUN mkdir /models
WORKDIR /models
RUN wget https://github.com/CSSLab/maia-chess/releases/download/v1.0/maia-1100.pb.gz
RUN wget https://github.com/CSSLab/maia-chess/releases/download/v1.0/maia-1200.pb.gz
RUN wget https://github.com/CSSLab/maia-chess/releases/download/v1.0/maia-1300.pb.gz
RUN wget https://github.com/CSSLab/maia-chess/releases/download/v1.0/maia-1400.pb.gz
RUN wget https://github.com/CSSLab/maia-chess/releases/download/v1.0/maia-1500.pb.gz
RUN wget https://github.com/CSSLab/maia-chess/releases/download/v1.0/maia-1600.pb.gz
RUN wget https://github.com/CSSLab/maia-chess/releases/download/v1.0/maia-1700.pb.gz
RUN wget https://github.com/CSSLab/maia-chess/releases/download/v1.0/maia-1800.pb.gz
RUN wget https://github.com/CSSLab/maia-chess/releases/download/v1.0/maia-1900.pb.gz

# Install software
WORKDIR /tmp
RUN git clone -b release/0.30 --recurse-submodules https://github.com/LeelaChessZero/lc0.git
RUN cd lc0 && ./build.sh && install /tmp/lc0/build/release/lc0 /usr/bin
RUN curl -L https://github.com/cutechess/cutechess/releases/download/v1.3.1/cutechess-cli-1.3.1-linux64.tar.gz | tar xz && \
  install cutechess-cli/cutechess-cli /usr/bin && \
  cutechess-cli --version
RUN curl -LO https://github.com/algerbrex/blunder/releases/download/v8.5.5/blunder-8.5.5.zip && \
  unzip blunder-8.5.5.zip && \
  install -T blunder-8.5.5/linux/blunder-8.5.5-default /usr/bin/blunder && \
  which blunder
RUN curl -L https://github.com/official-stockfish/Stockfish/releases/download/sf_18/stockfish-ubuntu-x86-64.tar | tar x && \
  install -T stockfish/stockfish-ubuntu-x86-64 /usr/bin/stockfish
RUN curl -L https://github.com/Disservin/fastchess/releases/download/v1.8.2-alpha/fastchess-linux-x86-64.tar | tar x && \
  install fastchess-linux-x86-64/fastchess /usr/bin && rm -rf fastchess-linux-x86-64

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
