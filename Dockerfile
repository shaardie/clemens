FROM golang:1.21.6-bookworm

# Install cutechess-cli and ordo and stockfish
WORKDIR /tmp
RUN apt-get update && apt-get install -y curl libqt5core5a unzip && \
    curl -L https://github.com/cutechess/cutechess/releases/download/v1.3.1/cutechess-cli-1.3.1-linux64.tar.gz | tar xz && \
    install cutechess-cli/cutechess-cli /usr/bin && \
    cutechess-cli --version && \
    curl -L https://github.com/michiguel/Ordo/releases/download/v1.2.6/ordo-1.2.6.tar.gz | tar xz && \
    install -T ordo-linux64 /usr/bin/ordo && \
    ordo --version && \
    curl -L https://github.com/official-stockfish/Stockfish/releases/download/sf_16/stockfish-ubuntu-x86-64.tar | tar x -C /tmp && \
    install -T stockfish/stockfish-ubuntu-x86-64 /usr/bin/stockfish && \
    stockfish --version && \
    curl -LO https://github.com/algerbrex/blunder/releases/download/v8.5.5/blunder-8.5.5.zip && \
    unzip blunder-8.5.5.zip && \
    install -T blunder-8.5.5/linux/blunder-8.5.5-default /usr/bin/blunder && \
    which blunder

# Install current clemens engine
RUN mkdir /go/clemens
WORKDIR /go/clemens
COPY go.mod go.sum Makefile ./
RUN go mod download
COPY pkg pkg
COPY cmd cmd
RUN make && install clemens /usr/bin

# Copy Openings
COPY openings.pgn scripts/elo.sh /

# Set tournement script
CMD [ "/elo.sh" ]
