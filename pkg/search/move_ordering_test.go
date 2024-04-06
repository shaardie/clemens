package search

// func TestSearch_orderMoves(t *testing.T) {
// 	pos, err := position.NewFromFen("r3k2r/p1ppqpb1/Bn4p1/3pN3/4nP2/2B5/PPP3QP/R3K2R w KQkq - 0 5")
// 	require.NoError(t, err)

// 	// Generate all moves and order them
// 	moves := move.NewMoveList()
// 	pos.GeneratePseudoLegalMoves(moves)
// 	s := NewSearch(*pos)
// 	s.Search(context.TODO(), SearchParameter{Depth: 8})
// 	s.orderMoves(pos, moves, s.PV.GetBestMove(), move.NullMove, 0)
// 	t.Log(moves)
// 	t.Fail()
// }
