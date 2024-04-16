#!/usr/bin/env python

import chess.pgn
import csv
import sys

MAXINT16 = 32767
MATE_VALUES = MAXINT16 - 100

def convert(input_filename: str, output_filename: str):
    x = 0
    fens = []
    results = []
    with open(input_filename) as input_file:

            while True:
                # Just to print something during the run
                x += 1
                if x % 1000 == 0:
                    print(x)

                # Get Chess Game from PGN
                game = chess.pgn.read_game(input_file)
                if not game:
                    break

                # Get result from header
                result_string = game.headers["Result"]
                if result_string == "0-1":
                    result = 0
                elif result_string == "1-0":
                    result = 1
                else:
                    result = 0.5

                # Iterate trough positions
                while True:
                    game = game.next()
                    if not game:
                        break

                    # Ignore Book Moves
                    if game.comment == "book":
                        continue

                    # Ignore Mate Scores
                    try:
                        value = int(float(game.comment.split("/")[0]) * 100)
                        if value > MATE_VALUES or value < -MATE_VALUES:
                            continue
                    except ValueError:
                        # probably a Mate Value
                        continue

                    # Get Fen String and skip existing ones
                    fen = game.board().fen()
                    # if fen in fens:
                    #     continue
                    fens.append(fen)
                    results.append(result)

    with open(output_filename, "w") as output_file:
        csv_writer = csv.writer(output_file)
        for idx in range(len(fens)):
            # Write directly to csv
            csv_writer.writerow([fens[idx], results[idx]])


def main():
    try:
        input_filename = sys.argv[1]
        output_filename = sys.argv[2]
    except KeyError:
        print(f"Usage: {sys.argv[0]} <input_filename> <output_filename>")
        sys.exit(1)

    convert(input_filename, output_filename)


if __name__ == "__main__":
    main()
