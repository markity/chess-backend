package chess

import "chess-backend/comm/chess"

type MoveResult struct {
	OK            bool
	GameOver      bool
	BlackScoreAdd int
	WhiteScoreAdd int
	PawnUpgrade   bool
}

// 输入规则: 不同且合法的坐标
func DoMove(table *chess.ChessTable, side chess.Side, fromX rune, fromY int, toX rune, toY int) (result MoveResult) {
	// 首先判断from处是否有棋子
	if table.GetPosition(fromX, fromY) == nil {
		result.OK = false
		return
	}

	// 棋子必须是自己的
	if table.GetPosition(fromX, fromY).GameSide != side {
		result.OK = false
		return
	}

	return
}
