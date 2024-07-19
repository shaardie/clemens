import logging
import argparse

import torch

from torch.utils.data import Dataset, DataLoader

logger = logging.getLogger(__name__)

NUM_FEATURES = 64 * 64 * 5 * 2 * 2
M = 4
N = 8
K = 1


class LSB:
    def __iter__(self, bitboard: int):
        self.b = bitboard
        return self

    def __next__(self):
        x = self.bitboard & -self.bitboard
        self.bitboard &= self.bitboard - 1
        return x


PAWN = 0
KNIGHT = 1
BISHOP = 2
ROOK = 3
QUEEN = 4
KING = 5

ExtPieceToClemens = [
    KNIGHT,
    BISHOP,
    ROOK,
    QUEEN,
    KING,
    PAWN,
    ROOK,
    PAWN,
]


WHITE = 0
BLACK = 1


class PositionsDataset(Dataset):
    def __init__(self, filename):
        self.data = []
        with open(filename, "rb") as f:
            while True:
                occ_bytes = f.read(8)
                if len(occ_bytes) == 0:
                    break
                occ = int.from_bytes(occ_bytes, signed=False)
                number_of_pieces = occ.bit_count()
                assert number_of_pieces <= 32
                turn_and_rules50 = int.from_bytes(f.read(1))
                turn = turn_and_rules50 & 1
                assert turn <= 1
                rules50 = turn_and_rules50 >> 1
                assert rules50 <= 100
                packed_pieces_size = (number_of_pieces + 1) // 2
                assert packed_pieces_size <= 16
                packed_pieces = f.read(packed_pieces_size)
                pieces = []
                kings = [None, None]
                for b in packed_pieces:
                    for piece in b & 0x0F, b >> 4:
                        square = (occ & -occ).bit_length() - 1
                        occ &= occ - 1
                        piece_type = ExtPieceToClemens[(piece & 0xFE) // 2]
                        piece_color = piece & 1
                        if piece_type == KING:
                            kings[piece_color] = square
                            continue
                        pieces.append((square, piece_type, piece_color))

                score = int.from_bytes(f.read(4), signed=True, byteorder="little")
                result = int.from_bytes(f.read(4), signed=False, byteorder="little")

                white_features = torch.zeros(NUM_FEATURES)
                black_features = torch.zeros(NUM_FEATURES)
                for piece in pieces:
                    white_features[self.__calc_index__(piece, kings[WHITE])] = 1
                    black_features[self.__calc_index__(piece, kings[BLACK])] = 1
                white_features = white_features.to_sparse()
                black_features = black_features.to_sparse()
                self.data.append((white_features, black_features, turn, score, result))

                if len(self.data) % 1000 == 0:
                    logger.debug(f"{len(self.data)} positions read")

    def __len__(self):
        return len(self.data)

    def __getitem__(self, idx):
        return self.data[idx]

    def __calc_index__(self, piece, king):
        piece_index = piece[1] * 2 + piece[2]
        return piece[0] + (piece_index + king)


class NNUE(torch.nn.Module):
    def __init__(self):
        super(NNUE, self).__init__()

        self.ft = torch.nn.Linear(NUM_FEATURES, M)
        self.l1 = torch.nn.Linear(2 * M, N)
        self.l2 = torch.nn.Linear(N, K)

    # The inputs are a whole batch!
    # `turn` indicates whether white is the side to move. 1 = true, 0 = false.
    def forward(self, white_features, black_features, turn, score, result):
        w = self.ft(white_features)  # white's perspective
        b = self.ft(black_features)  # black's perspective

        # Remember that we order the accumulators for 2 perspectives based on who is to move.
        # So we blend two possible orderings by interpolating between `stm` and `1-stm` tensors.
        accumulator = (turn * torch.cat([w, b], dim=1)) + (
            (1 - turn) * torch.cat([b, w], dim=1)
        )

        # Run the linear layers and use clamp_ as ClippedReLU
        l1_x = torch.clamp(accumulator, 0.0, 1.0)
        l2_x = torch.clamp(self.l1(l1_x), 0.0, 1.0)
        model_result = self.l2(l2_x)

        # Loss function
        scaling_factor = 400 # TODO better value
        lambda_ = 0.5 # TODO better value
        wdl_eval_model = torch.sigmoid(model_result / scaling_factor)
        wdl_eval_target = torch.sigmoid(score / scaling_factor)
        loss_eval   = (wdl_eval_model - wdl_eval_target)**2
        loss_result = (wdl_eval_model - result)**2
        loss = lambda_ * loss_eval + (1 - lambda_) * loss_result

        return loss


def init():
    # Command line arguments
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--verbose",
        "-v",
        action="store_true",
        default=False,
        help="be more verbose",
    )
    parser.add_argument(
        "--dataset",
        required=True,
        help="datasets used for training, can be given multiple times",
    )
    args = parser.parse_args()

    # Configure logging
    logging.basicConfig(level=logging.DEBUG if args.verbose else logging.INFO)

    return args


def collate(data):
    data = zip(*data)
    white_features = torch.stack(next(data))
    black_features = torch.stack(next(data))
    turn = torch.tensor(next(data)).reshape(-1,1)
    score = torch.tensor(next(data)).reshape(-1,1)
    result = torch.tensor(next(data)).reshape(-1,1)
    return (
        white_features,
        black_features,
        turn,
        score,
        result,
    )


def main():
    args = init()
    dataset = PositionsDataset(args.dataset)
    dataloader = DataLoader(dataset, batch_size=4, collate_fn=collate)
    nn = NNUE().to("cpu")

    for x in iter(dataloader):
        nn(*x)

    print(dir(nn))


if __name__ == "__main__":
    main()
