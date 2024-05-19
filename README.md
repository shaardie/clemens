# Clemens

Clemens is an [UCI](https://www.shredderchess.com/de/schach-features/uci-universal-chess-interface.html) compatible Chess Engine written in [Go](https://go.dev/).


## Usage

The binaries for Linux and Windows can be downloaded from the [Releases](https://github.com/shaardie/clemens/releases) section.
Those binaries can be used in every UCI compatible Chess GUI, like [Cute Chess](https://cutechess.com/) or the [Arena Chess GUI](http://www.playwitharena.de/).
There is also a small web UI at [https://chess.haardiek.org/](https://chess.haardiek.org/) where you can play against the current version,
Of course, this version is very limited.

## Building from Source

You can also build the program from the source.
The project includes a [Makefile](./Makefile) for this purpose.
Executing `make` will build the binaries.
The binaries with the `.exe` extension are for Windows, while the others are for Linux.
The chess engine is the binary named `clemens`.

## Changelog

### Dev

* Evaluate [Isolanis](https://www.chessprogramming.org/Isolated_Pawn).
* Evaluate [Doubled Pawns](https://www.chessprogramming.org/Doubled_Pawn).
* Evaluate [Passed Pawns](https://www.chessprogramming.org/Passed_Pawn).
* [Transposition Table](https://www.chessprogramming.org/Transposition_Table) for Evaluations.
* [Repetition](https://www.chessprogramming.org/Repetitions)
* [Firty-Move Rule](https://www.chessprogramming.org/Fifty-move_Rule) in Evaluations.
* Bug Fix King Shield Evaluation.
* Bug Fix Piece Square Tables.
* Inplace Zobrist Hash Updates while making the move.
* Bug Fix Mobility and King Attacks.
* [Draw Evaluation](https://www.chessprogramming.org/Draw_Evaluation) and [Contempt Value](https://www.chessprogramming.org/Contempt_Factor).
* [Static Exchange Evaluation](https://www.chessprogramming.org/Static_Exchange_Evaluation) in Quiesence Search.
* Major Bug Fix in Transposition Tables.
* [Aging](https://www.chessprogramming.org/Transposition_Table#Aging) in Transposition Tables.
* [Killer Moves](https://www.chessprogramming.org/Killer_Move) in move order.
* Remove [Mate Distance Pruning](https://www.chessprogramming.org/Mate_Distance_Pruning).
* Remove [Null Move Pruning](https://www.chessprogramming.org/Null_Move_Pruning).
* [Futility Pruning](https://www.chessprogramming.org/Futility_Pruning)
* [Delta Pruning](https://www.chessprogramming.org/Delta_Pruning)
* Squares as uint8.
* [Principal Variation Search](https://www.chessprogramming.org/Principal_Variation_Search)
* Slightly better time control.
* [Static Null Move Pruning](https://www.chessprogramming.org/Reverse_Futility_Pruning).

### v0.3.0

* Move Ordering with [MVV-LVA](https://www.chessprogramming.org/MVV-LVA).
* Checkmate and Stalemate Detection.
* Improved [King Safery](https://www.chessprogramming.org/King_Pattern#King_Safety) by evaluating the King Shield.
* Pair Evaluation, e.g. [Bishop Pair](https://www.chessprogramming.org/Bishop_Pair).
* [Game Phases](https://www.chessprogramming.org/Game_Phases) based Evaluation.
* Piece Adjustments based on Pawns.
* [Rook Evaluation](https://www.chessprogramming.org/Evaluation_of_Pieces#Rook).
* [Mate Distance Pruning](https://www.chessprogramming.org/Mate_Distance_Pruning).
* ELO Rating Script.
* [Mobility](https://www.chessprogramming.org/Mobility) Evaluation.
* [King Zone Attacks](https://www.chessprogramming.org/King_Safety#Attacking_King_Zone).
* Move from Full Move to Ply.
* [Null Move Pruning](https://www.chessprogramming.org/Null_Move_Pruning).

### v0.2.0

* Better [UCI](https://www.shredderchess.com/de/schach-features/uci-universal-chess-interface.html) Handling.
* [Transposition Table](https://www.chessprogramming.org/Transposition_Table) for Search Results.
* Speed up using some Implementation Details, using less Memory on the Heap and less Memory Allocation in general.
* Better Mobility with [Piece-Square Tables](https://www.chessprogramming.org/Piece-Square_Tables) in Evaluation.
* [Quiescence Search](https://www.chessprogramming.org/Quiescence_Search).
* Add [Aspiration Windows](https://www.chessprogramming.org/Aspiration_Windows) to Search.
* Basic [Time Management](https://www.chessprogramming.org/Time_Management).
* Use best Move from last Search first.
* [Iterative Search](https://www.chessprogramming.org/Iterative_Search)
* [Principal Variation](https://www.chessprogramming.org/Principal_Variation) Line.

### v0.1.0

The first version that you could somehow play against.

* [Bitboard](https://www.chessprogramming.org/Bitboards) Implementation
* [Sliding Pieces](https://www.chessprogramming.org/Sliding_Pieces)
* [Magic Bitboards](https://www.chessprogramming.org/Magic_Bitboards)
* [Chess Position](https://www.chessprogramming.org/Chess_Position)
* Position from and to [FEN String](https://www.chessprogramming.org/Forsyth-Edwards_Notation)
* [Generate Pseudo Legal Moves](https://www.chessprogramming.org/Move_Generation)
* [Make Move](https://www.chessprogramming.org/Make_Move) Funktion
* [Perft](https://www.chessprogramming.org/Perft) Testing
* [Basic Evaluation](https://www.chessprogramming.org/Evaluation#Where_to_Start)
* Rudimentary[UCI](https://www.shredderchess.com/de/schach-features/uci-universal-chess-interface.html) Implementation
* [Alpha-Beta](https://www.chessprogramming.org/Alpha-Beta) Search Funktion

## Special Thanks

Since I had no idea how to write a chess engine at the beginning of the project, I did a lot of reading.
Therefore, my special thanks go to the [Chess Programming Wiki](https://www.chessprogramming.org/Main_Page), an excellent source of information, but also to the engines [Blunder](https://github.com/algerbrex/blunder), [CPW Engine](https://github.com/nescitus/cpw-engine) and [Stockfish](https://github.com/official-stockfish/Stockfish), from whose code I learned a lot.

Also, to my father, who reintroduced me to playing chess and thus was the inspiration for this project.
That's why the engine bears his name.
