package chess

import (
	"chess-backend/comm/chess"
	"math"
)

type MoveResult struct {
	OK       bool
	GameOver bool
	// 可能出现逼将和, 这时是平局
	GameWinner chess.Side
	// 兵升变
	PawnUpgrade bool
	// 将军
	KingThreat bool
}

func isIndexMatch(x1 int, y1 int, x2 int, y2 int) bool {
	return x1 == x2 && y1 == y2
}

func isOnSameLine(aX rune, aY int, bX rune, bY int) bool {
	return aX == bX || aY == bY
}

// 判断两个坐标是不是倾斜的
func isTwoIndexIncline(ax int, ay int, bx int, by int) bool {
	// 防止除0异常
	return math.Abs(float64(ax-bx)) == math.Abs(float64(ay-by))
}

// 判定两个点中间是否有直线
// 要求在同一条线上
func hasChessBetweenTwoPointsInLine(table *chess.ChessTable, aX rune, aY int, bX rune, bY int) bool {
	x1, y1 := chess.MustPositionToIndex(aX, aY)
	x2, y2 := chess.MustPositionToIndex(bX, bY)
	if x1 == x2 {
		var yMin int
		var yMax int
		if y1 > y2 {
			yMin = y2
			yMax = y1
		} else {
			yMin = y1
			yMax = y2
		}

		for y0 := yMin + 1; y0 < yMax; y0++ {
			if table[y0*8+x1] != nil {
				return true
			}
		}
	} else {
		var xMin int
		var xMax int
		if x1 > x2 {
			xMin = x2
			xMax = x1
		} else {
			xMin = x1
			xMax = x2
		}

		for x0 := xMin + 1; x0 < xMax; x0++ {
			if table[y1*8+x0] != nil {
				return true
			}
		}
	}

	return false
}

// 这个函数要求两个点是倾斜排布的
// 判定两个倾斜的点中间是否有棋挡住
func hasChessBetweenTwoInclinedPoints(table *chess.ChessTable, ax, bx, ay, by int) bool {
	var diffX int
	var diffY int

	if by > ay {
		diffY = 1
	} else {
		diffY = -1
	}

	if bx > ax {
		diffX = 1
	} else {
		diffX = -1
	}

	for x0, y0 := ax+diffX, bx+diffX; CheckChessIndexValid(x0, y0) && x0 != bx; x0, y0 = x0+diffX, y0+diffY {
		if table.GetIndex(x0, y0) != nil {
			return true
		}
	}

	return false
}

func checkIndexThreat(table *chess.ChessTable, side chess.Side, x int, y int) bool {
	// 对方的side
	var selfSide = side
	var remoteSide chess.Side
	if selfSide == chess.SideBlack {
		remoteSide = chess.SideWhite
	} else {
		remoteSide = chess.SideBlack
	}

	// 0. 检查周围9个格子内有没有对方的king
	if x0, y0 := x+1, y+1; CheckChessIndexValid(x0, y0) {
		if table.GetIndex(x0, y0) != nil &&
			table.GetIndex(x0, y0).PieceType == chess.ChessPieceTypeKing && table.GetIndex(x0, y0).GameSide != side {
			return true
		}
	}
	if x0, y0 := x+1, y-1; CheckChessIndexValid(x0, y0) {
		if table.GetIndex(x0, y0) != nil &&
			table.GetIndex(x0, y0).PieceType == chess.ChessPieceTypeKing && table.GetIndex(x0, y0).GameSide != side {
			return true
		}
	}
	if x0, y0 := x-1, y+1; CheckChessIndexValid(x0, y0) {
		if table.GetIndex(x0, y0) != nil &&
			table.GetIndex(x0, y0).PieceType == chess.ChessPieceTypeKing && table.GetIndex(x0, y0).GameSide != side {
			return true
		}
	}
	if x0, y0 := x-1, y-1; CheckChessIndexValid(x0, y0) {
		if table.GetIndex(x0, y0) != nil &&
			table.GetIndex(x0, y0).PieceType == chess.ChessPieceTypeKing && table.GetIndex(x0, y0).GameSide != side {
			return true
		}
	}
	if x0, y0 := x, y+1; CheckChessIndexValid(x0, y0) {
		if table.GetIndex(x0, y0) != nil &&
			table.GetIndex(x0, y0).PieceType == chess.ChessPieceTypeKing && table.GetIndex(x0, y0).GameSide != side {
			return true
		}
	}
	if x0, y0 := x, y-1; CheckChessIndexValid(x0, y0) {
		if table.GetIndex(x0, y0) != nil &&
			table.GetIndex(x0, y0).PieceType == chess.ChessPieceTypeKing && table.GetIndex(x0, y0).GameSide != side {
			return true
		}
	}
	if x0, y0 := x+1, y; CheckChessIndexValid(x0, y0) {
		if table.GetIndex(x0, y0) != nil &&
			table.GetIndex(x0, y0).PieceType == chess.ChessPieceTypeKing && table.GetIndex(x0, y0).GameSide != side {
			return true
		}
	}
	if x0, y0 := x-1, y; CheckChessIndexValid(x0, y0) {
		if table.GetIndex(x0, y0) != nil &&
			table.GetIndex(x0, y0).PieceType == chess.ChessPieceTypeKing && table.GetIndex(x0, y0).GameSide != side {
			return true
		}
	}

	// 1. 检查pawn的threat, 即为对角位置是兵
	if selfSide == chess.SideWhite {
		if x0, y0 := x+1, y+1; CheckChessIndexValid(x0, y0) {
			p := table.GetIndex(x0, y0)
			if p != nil && p.PieceType == chess.ChessPieceTypePawn && p.GameSide == remoteSide {
				return true
			}
		}

		if x0, y0 := x-1, y+1; CheckChessIndexValid(x0, y0) {
			p := table.GetIndex(x0, y0)
			if p != nil && p.PieceType == chess.ChessPieceTypePawn && p.GameSide == remoteSide {
				return true
			}
		}
	} else {
		if x0, y0 := x+1, y-1; CheckChessIndexValid(x0, y0) {
			p := table.GetIndex(x0, y0)
			if p != nil && p.PieceType == chess.ChessPieceTypePawn && p.GameSide == remoteSide {
				return true
			}
		}

		if x0, y0 := x-1, y-1; CheckChessIndexValid(x0, y0) {
			p := table.GetIndex(x0, y0)
			if p != nil && p.PieceType == chess.ChessPieceTypePawn && p.GameSide == remoteSide {
				return true
			}
		}
	}

	// 2. 检查bishop和queen threat, 遍历所有的斜对角位置
	// 右上角
	for x0, y0 := x+1, y+1; CheckChessIndexValid(x0, y0); x0, y0 = x0+1, y0+1 {
		p := table.GetIndex(x0, y0)
		if p != nil {
			// 被自家棋子挡住了, 这是安全的
			if p.GameSide == selfSide {
				break
			}

			if p.GameSide == remoteSide && (p.PieceType == chess.ChessPieceTypeBishop ||
				p.PieceType == chess.ChessPieceTypeQueen) {
				return true
			}
		}
	}
	// 右下角
	for x0, y0 := x+1, y-1; CheckChessIndexValid(x0, y0); x0, y0 = x0+1, y0-1 {
		p := table.GetIndex(x0, y0)
		if p != nil {
			// 被自家棋子挡住了, 这是安全的
			if p.GameSide == selfSide {
				break
			}

			if p.GameSide == remoteSide && (p.PieceType == chess.ChessPieceTypeBishop ||
				p.PieceType == chess.ChessPieceTypeQueen) {
				return true
			}
		}
	}
	// 左下角
	for x0, y0 := x-1, y-1; CheckChessIndexValid(x0, y0); x0, y0 = x0-1, y0-1 {
		p := table.GetIndex(x0, y0)
		if p != nil {
			// 被自家棋子挡住了, 这是安全的
			if p.GameSide == selfSide {
				break
			}

			if p.GameSide == remoteSide && (p.PieceType == chess.ChessPieceTypeBishop ||
				p.PieceType == chess.ChessPieceTypeQueen) {
				return true
			}
		}
	}
	// 左上角
	for x0, y0 := x-1, y+1; CheckChessIndexValid(x0, y0); x0, y0 = x0-1, y0+1 {
		p := table.GetIndex(x0, y0)
		if p != nil {
			// 被自家棋子挡住了, 这是安全的
			if p.GameSide == selfSide {
				break
			}

			if p.GameSide == remoteSide && (p.PieceType == chess.ChessPieceTypeBishop ||
				p.PieceType == chess.ChessPieceTypeQueen) {
				return true
			}
		}
	}

	// 3. 检查rook和queen threat
	for x0, y0 := x+1, y; CheckChessIndexValid(x0, y0); x0 = x0 + 1 {
		p := table.GetIndex(x0, y0)
		if p != nil {
			// 被自家棋子挡住了, 这是安全的
			if p.GameSide == selfSide {
				break
			}

			if p.GameSide == remoteSide && (p.PieceType == chess.ChessPieceTypeRook ||
				p.PieceType == chess.ChessPieceTypeQueen) {
				return true
			}
		}
	}
	for x0, y0 := x-1, y; CheckChessIndexValid(x0, y0); x0 = x0 - 1 {
		p := table.GetIndex(x0, y0)
		if p != nil {
			// 被自家棋子挡住了, 这是安全的
			if p.GameSide == selfSide {
				break
			}

			if p.GameSide == remoteSide && (p.PieceType == chess.ChessPieceTypeRook ||
				p.PieceType == chess.ChessPieceTypeQueen) {
				return true
			}
		}
	}
	for x0, y0 := x, y+1; CheckChessIndexValid(x0, y0); y0 = y0 + 1 {
		println(x0, y0)
		p := table.GetIndex(x0, y0)
		if p != nil {
			// 被自家棋子挡住了, 这是安全的
			if p.GameSide == selfSide {
				break
			}

			if p.GameSide == remoteSide && (p.PieceType == chess.ChessPieceTypeRook ||
				p.PieceType == chess.ChessPieceTypeQueen) {
				return true
			}
		}
	}
	for x0, y0 := x, y-1; CheckChessIndexValid(x0, y0); y0 = y0 - 1 {
		p := table.GetIndex(x0, y0)
		if p != nil {
			// 被自家棋子挡住了, 这是安全的
			if p.GameSide == selfSide {
				break
			}

			if p.GameSide == remoteSide && (p.PieceType == chess.ChessPieceTypeRook ||
				p.PieceType == chess.ChessPieceTypeQueen) {
				return true
			}
		}
	}

	// 检查knight threat, 8个可能的位置
	if x0, y0 := x+2, y+1; CheckChessIndexValid(x0, y0) {
		p := table.GetIndex(x0, y0)
		if p != nil && p.PieceType == chess.ChessPieceTypeKnight && p.GameSide == remoteSide {
			return true
		}
	}
	if x0, y0 := x+2, y-1; CheckChessIndexValid(x0, y0) {
		p := table.GetIndex(x0, y0)
		if p != nil && p.PieceType == chess.ChessPieceTypeKnight && p.GameSide == remoteSide {
			return true
		}
	}
	if x0, y0 := x-2, y+1; CheckChessIndexValid(x0, y0) {
		p := table.GetIndex(x0, y0)
		if p != nil && p.PieceType == chess.ChessPieceTypeKnight && p.GameSide == remoteSide {
			return true
		}
	}
	if x0, y0 := x-2, y-1; CheckChessIndexValid(x0, y0) {
		p := table.GetIndex(x0, y0)
		if p != nil && p.PieceType == chess.ChessPieceTypeKnight && p.GameSide == remoteSide {
			return true
		}
	}
	if x0, y0 := x+1, y+2; CheckChessIndexValid(x0, y0) {
		p := table.GetIndex(x0, y0)
		if p != nil && p.PieceType == chess.ChessPieceTypeKnight && p.GameSide == remoteSide {
			return true
		}
	}
	if x0, y0 := x+1, y-2; CheckChessIndexValid(x0, y0) {
		p := table.GetIndex(x0, y0)
		if p != nil && p.PieceType == chess.ChessPieceTypeKnight && p.GameSide == remoteSide {
			return true
		}
	}
	if x0, y0 := x-1, y+2; CheckChessIndexValid(x0, y0) {
		p := table.GetIndex(x0, y0)
		if p != nil && p.PieceType == chess.ChessPieceTypeKnight && p.GameSide == remoteSide {
			return true
		}
	}
	if x0, y0 := x-1, y-2; CheckChessIndexValid(x0, y0) {
		p := table.GetIndex(x0, y0)
		if p != nil && p.PieceType == chess.ChessPieceTypeKnight && p.GameSide == remoteSide {
			return true
		}
	}

	return false
}

// 检查某一个单元格是否受对方威胁
func checkPositionThreat(table *chess.ChessTable, side chess.Side, X rune, Y int) bool {
	x, y := chess.MustPositionToIndex(X, Y)
	return checkIndexThreat(table, side, x, y)
}

// 检查一个单元格的周围8个是否都受威胁
func checkAround8Threat(table *chess.ChessTable, side chess.Side, X rune, Y int) bool {
	x, y := chess.MustPositionToIndex(X, Y)

	// 1
	if x0, y0 := x+1, y+1; !CheckChessIndexValid(x0, y0) {
		if checkIndexThreat(table, side, x0, y0) {
			return false
		}
	}

	// 2
	if x0, y0 := x+1, y-1; !CheckChessIndexValid(x0, y0) {
		if !checkIndexThreat(table, side, x0, y0) {
			return false
		}
	}

	// 3
	if x0, y0 := x-1, y-1; !CheckChessIndexValid(x0, y0) {
		if !checkIndexThreat(table, side, x0, y0) {
			return false
		}
	}

	// 4
	if x0, y0 := x-1, y+1; !CheckChessIndexValid(x0, y0) {
		if !checkIndexThreat(table, side, x0, y0) {
			return false
		}
	}

	// 5
	if x0, y0 := x+1, y; !CheckChessIndexValid(x0, y0) {
		if !checkIndexThreat(table, side, x0, y0) {
			return false
		}
	}

	// 6
	if x0, y0 := x-1, y; !CheckChessIndexValid(x0, y0) {
		if !checkIndexThreat(table, side, x0, y0) {
			return false
		}
	}

	// 7
	if x0, y0 := x, y+1; !CheckChessIndexValid(x0, y0) {
		if !checkIndexThreat(table, side, x0, y0) {
			return false
		}
	}

	// 8
	if x0, y0 := x, y-1; !CheckChessIndexValid(x0, y0) {
		if !checkIndexThreat(table, side, x0, y0) {
			return false
		}
	}

	return true
}

func findKing(table *chess.ChessTable, side chess.Side) *chess.ChessPiece {
	for _, v := range table {
		if v != nil && v.GameSide == side && v.PieceType == chess.ChessPieceTypeKing {
			return v
		}
	}
	panic("unreachable")
}

func findJustMoved2Pawn(table *chess.ChessTable) *chess.ChessPiece {
	for _, v := range table {
		if v != nil && v.PieceType == chess.ChessPieceTypePawn && v.PawnMovedTwoLastTime {
			return v
		}
	}

	return nil
}

func findAllJustMoved2Pawn(table *chess.ChessTable) []*chess.ChessPiece {
	result := make([]*chess.ChessPiece, 0)
	for _, v := range table {
		if v != nil && v.PieceType == chess.ChessPieceTypePawn && v.PawnMovedTwoLastTime {
			result = append(result, v)
		}
	}

	return result
}

// 输入规则: 不同且合法的坐标
func DoMove(table *chess.ChessTable, side chess.Side, fromX rune, fromY int, toX rune, toY int) (result MoveResult) {
	// 一些要用到的基本数据
	fromx, fromy := chess.MustPositionToIndex(fromX, fromY)
	tox, toy := chess.MustPositionToIndex(toX, toY)
	// 对方的side
	var remoteSide chess.Side
	if side == chess.SideBlack {
		remoteSide = chess.SideWhite
	} else {
		remoteSide = chess.SideBlack
	}

	fromPiece := table.GetPosition(fromX, fromY)
	toPiece := table.GetPosition(toX, toY)

	// 下面是一些基本判断

	// 首先判断from处是否有棋子
	if fromPiece == nil {
		result.OK = false
		return
	}

	// 棋子必须是自己的
	if fromPiece.GameSide != side {
		result.OK = false
		return
	}

	// 目的地的棋子不能是自己的棋子
	if toPiece != nil && toPiece.GameSide == side {
		result.OK = false
		return
	}

	// 判断棋子的类型, 4种基本棋子做一套逻辑
	if fromPiece.PieceType != chess.ChessPieceTypeKing &&
		fromPiece.PieceType != chess.ChessPieceTypePawn {
		switch fromPiece.PieceType {
		case chess.ChessPieceTypeRook:
			// 不在同一条直线
			if !isOnSameLine(fromX, fromY, toX, toY) {
				result.OK = false
				return
			}

			// 中间有别的棋子
			if hasChessBetweenTwoPointsInLine(table, fromX, fromY, toX, toY) {
				result.OK = false
				return
			}
		case chess.ChessPieceTypeKnight:
			// 列举8个可能的位置
			if !isIndexMatch(tox, toy, fromx+2, fromy+1) && !isIndexMatch(tox, toy, fromx+2, fromy-1) &&
				!isIndexMatch(tox, toy, fromx-2, fromy+1) && !isIndexMatch(tox, toy, fromx-2, fromy-1) &&
				!isIndexMatch(tox, toy, fromx+1, fromy+2) && !isIndexMatch(tox, toy, fromx+1, fromy-2) &&
				!isIndexMatch(tox, toy, fromx-1, fromy+2) && !isIndexMatch(tox, toy, fromx-1, fromy-2) {
				result.OK = false
				return
			}

			if toPiece != nil && toPiece.GameSide == side {
				result.OK = false
				return
			}
		case chess.ChessPieceTypeBishop:
			// 非倾斜
			if !isTwoIndexIncline(fromx, fromy, tox, toy) {
				result.OK = false
				return
			}

			// 有格挡
			if hasChessBetweenTwoInclinedPoints(table, fromx, fromy, tox, toy) {
				result.OK = false
				return
			}

			// 目的地是自己的棋子, 前面已经判断了
			// if toPiece != nil && toPiece.GameSide == side {
			// 	result.OK = false
			// 	return
			// }
		case chess.ChessPieceTypeQueen:
			// 不合法的位移
			if !isOnSameLine(fromX, fromY, toX, toY) && !isTwoIndexIncline(fromx, fromy, tox, toy) {
				result.OK = false
				return
			}

			// 有格挡
			if isOnSameLine(fromX, fromY, toX, toY) && hasChessBetweenTwoPointsInLine(table, fromX, fromY, toX, toY) {
				result.OK = false
				return
			}

			// 有格挡
			if isTwoIndexIncline(fromx, fromy, tox, toy) && hasChessBetweenTwoInclinedPoints(table, fromx, fromy, toy, toy) {
				result.OK = false
				return
			}
		}

		// 拷贝一个测试table, 用来测试移动后自己不会暴露于威胁
		testTable := table.Copy()
		// 移动棋盘的子
		testFromPiece := testTable.GetIndex(fromx, fromy)
		testTable.ClearPosition(fromX, fromY)
		testFromPiece.Moved = true
		testFromPiece.X = toX
		testFromPiece.Y = toY
		testTable.SetPosition(testFromPiece)

		selfKing := findKing(testTable, side)
		if checkPositionThreat(testTable, side, selfKing.X, selfKing.Y) {
			result.OK = false
			return
		}

		// 正式设置表
		table.ClearIndex(fromx, fromy)
		fromPiece.Moved = true
		fromPiece.X = toX
		fromPiece.Y = toY
		table.SetPosition(fromPiece)

		// 设置一下justMoved
		justMovedPawn := findJustMoved2Pawn(table)
		if justMovedPawn != nil {
			justMovedPawn.PawnMovedTwoLastTime = false
		}

		// 找到对方的king
		remoteKing := findKing(table, remoteSide)

		// 是否将军
		kingThreat := checkPositionThreat(table, remoteSide, remoteKing.X, remoteKing.Y)

		// 王的8个单元格是否都受威胁
		kingAroundAllThreat := checkAround8Threat(table, remoteSide, remoteKing.X, remoteKing.Y)

		// 赢
		if kingThreat && kingAroundAllThreat {
			result.OK = true
			result.GameOver = true
			result.GameWinner = side
			return
		}

		// 将军但游戏没有结束
		if kingThreat {
			result.OK = true
			result.GameOver = false
			result.KingThreat = true
			return
		}

		result.OK = true
		result.GameOver = false
		result.KingThreat = false
		return
	}

	didKingRookSwitch := false
	pawnUpgrade := false

	// 特别处理两个特别的子
	switch fromPiece.PieceType {
	case chess.ChessPieceTypeKing:
		wantRookKingSwitch := false

		// 需要判断是否只在九宫格里面移动, 或者是否想要发生车王易位
		if !isIndexMatch(fromx+1, fromy+1, tox, toy) && !isIndexMatch(fromx+1, fromy-1, tox, toy) &&
			!isIndexMatch(fromx-1, fromy+1, tox, toy) && !isIndexMatch(fromx-1, fromy-1, tox, toy) &&
			!isIndexMatch(fromx+1, fromy, tox, toy) && !isIndexMatch(fromx-1, fromy, tox, toy) &&
			!isIndexMatch(fromx, fromy+1, tox, toy) && !isIndexMatch(fromx, fromy-1, tox, toy) {
			// 简单判断下是否想要王车易位
			if side == chess.SideWhite && fromX == 'e' && fromY == 1 &&
				((toX == 'g' && toY == 1) || (toX == 'c' && toY == 1)) {
				wantRookKingSwitch = true
			} else if side == chess.SideBlack && fromX == 'e' && fromY == 8 &&
				((toX == 'g' && toY == 8) || (toX == 'c' && toY == 8)) {
				wantRookKingSwitch = true
			} else {
				// 不想王车易位
				result.OK = false
				return
			}
		}

		// 不试图王车易位, 那么直接判定to
		if !wantRookKingSwitch {
			// to的地方不能是自己人
			if toPiece != nil && toPiece.GameSide == side {
				result.OK = false
				return
			}

			testTable := table.Copy()
			testFromPiece := testTable.GetPosition(fromX, fromY)
			testTable.ClearPosition(fromX, fromY)
			testFromPiece.X = toX
			testFromPiece.Y = toY
			testTable.SetPosition(testFromPiece)

			// to的地方不能有威胁
			if checkIndexThreat(testTable, side, tox, toy) {
				result.OK = false
				return
			}
			// 特别逻辑: 王车易位
		} else {
			// 想要王车易位, 上面已经判断了from, to的坐标了

			if side == chess.SideWhite {
				// 短
				if toX == 'g' {
					rookPiece := table.GetPosition('h', 1)

					// 挡住
					if hasChessBetweenTwoPointsInLine(table, 'e', 1, 'h', 1) {
						result.OK = false
						return
					}

					// rook为nil或不为rook
					if rookPiece == nil || rookPiece.PieceType != chess.ChessPieceTypeRook {
						result.OK = false
						return
					}

					// 移动过
					if fromPiece.Moved || rookPiece.Moved {
						result.OK = false
						return
					}

					// 检查路过的威胁
					testTable1 := table.Copy()
					testFromPiece := testTable1.GetPosition(fromX, fromY)
					testRookPiece := testTable1.GetPosition('h', 1)
					testTable1.ClearPosition(fromX, fromY)
					testTable1.ClearPosition(testFromPiece.X, testRookPiece.Y)
					testFromPiece.X = 'f'
					testFromPiece.Y = 1
					testRookPiece.X = 'e'
					testRookPiece.Y = 1
					testTable1.SetPosition(testFromPiece)
					testTable1.SetPosition(testRookPiece)

					testTable2 := table.Copy()
					testFromPiece2 := testTable2.GetPosition(fromX, fromY)
					testRookPiece2 := testTable2.GetPosition('h', 1)
					testTable1.ClearPosition(fromX, fromY)
					testTable2.ClearPosition(testRookPiece2.X, testFromPiece2.Y)
					testFromPiece2.X = 'g'
					testFromPiece2.Y = 1
					testRookPiece2.X = 'e'
					testRookPiece2.Y = 1
					testTable1.SetPosition(testFromPiece2)
					testTable1.SetPosition(testRookPiece2)
					if checkPositionThreat(testTable1, side, 'f', 1) || checkPositionThreat(testTable2, side, 'g', 1) {
						result.OK = false
						return
					}

					// ok, 可以易位
					table.ClearPosition('e', 1)
					table.ClearPosition('h', 1)
					fromPiece.X = 'g'
					fromPiece.Y = 1
					fromPiece.Moved = true
					table.SetPosition(fromPiece)
					rookPiece.X = 'f'
					rookPiece.Y = 1
					fromPiece.Moved = true
					table.SetPosition(rookPiece)
					// 长 to = c
				} else {
					rookPiece := table.GetPosition('a', 1)

					// 挡住
					if hasChessBetweenTwoPointsInLine(table, 'e', 1, 'a', 1) {
						result.OK = false
						return
					}

					// rook为nil或不为rook
					if rookPiece == nil || rookPiece.PieceType != chess.ChessPieceTypeRook {
						result.OK = false
						return
					}

					// 移动过
					if fromPiece.Moved || rookPiece.Moved {
						result.OK = false
						return
					}

					// 检查路过的威胁
					testTable1 := table.Copy()
					testFromPiece := testTable1.GetPosition(fromX, fromY)
					testRookPiece := testTable1.GetPosition('a', 1)
					testTable1.ClearPosition(fromX, fromY)
					testTable1.ClearPosition(testRookPiece.X, testRookPiece.Y)
					testFromPiece.X = 'd'
					testFromPiece.Y = 1
					// 让车帮忙档一下, 比较易于算威胁
					testRookPiece.X = 'e'
					testRookPiece.Y = 1
					testTable1.SetPosition(testFromPiece)

					testTable2 := table.Copy()
					testFromPiece2 := testTable2.GetPosition(fromX, fromY)
					testRookPiece2 := testTable2.GetPosition('a', 1)
					testTable2.ClearPosition(fromX, fromY)
					testTable2.ClearPosition(testRookPiece2.X, testRookPiece2.Y)
					testFromPiece2.X = 'c'
					testFromPiece2.Y = 1
					testRookPiece2.X = 'd'
					testRookPiece2.Y = 1
					testTable2.SetPosition(testFromPiece2)
					testTable2.SetPosition(testRookPiece2)
					if checkPositionThreat(testTable1, side, 'd', 1) || checkPositionThreat(testTable2, side, 'c', 1) {
						result.OK = false
						return
					}

					// ok, 可以易位
					table.ClearPosition('e', 1)
					table.ClearPosition('a', 1)
					fromPiece.X = 'c'
					fromPiece.Y = 1
					fromPiece.Moved = true
					table.SetPosition(fromPiece)
					rookPiece.X = 'd'
					rookPiece.Y = 1
					fromPiece.Moved = true
					table.SetPosition(rookPiece)
				}
				// 黑方想要王车易位
			} else {
				// 短
				if toX == 'g' {
					rookPiece := table.GetPosition('h', 8)
					// 挡住
					if hasChessBetweenTwoPointsInLine(table, 'e', 8, 'h', 8) {
						result.OK = false
						return
					}

					// to为nil或不为rook
					if rookPiece == nil || rookPiece.PieceType != chess.ChessPieceTypeRook {
						result.OK = false
						return
					}

					// 移动过
					if fromPiece.Moved || toPiece.Moved {
						result.OK = false
						return
					}

					// 检查路过的威胁
					testTable1 := table.Copy()
					testFromPiece := testTable1.GetPosition(fromX, fromY)
					testRookPiece := testTable1.GetPosition('h', 8)
					testTable1.ClearPosition(fromX, fromY)
					testTable1.ClearPosition(testFromPiece.X, testRookPiece.Y)
					testFromPiece.X = 'f'
					testFromPiece.Y = 8
					testRookPiece.X = 'e'
					testRookPiece.Y = 8
					testTable1.SetPosition(testFromPiece)
					testTable1.SetPosition(testRookPiece)

					testTable2 := table.Copy()
					testFromPiece2 := testTable2.GetPosition(fromX, fromY)
					testRookPiece2 := testTable2.GetPosition('h', 8)
					testTable1.ClearPosition(fromX, fromY)
					testTable2.ClearPosition(testRookPiece2.X, testFromPiece2.Y)
					testFromPiece2.X = 'g'
					testFromPiece2.Y = 8
					testRookPiece2.X = 'e'
					testRookPiece2.Y = 8
					testTable1.SetPosition(testFromPiece2)
					testTable1.SetPosition(testRookPiece2)
					if checkPositionThreat(testTable1, side, 'f', 8) || checkPositionThreat(testTable2, side, 'g', 8) {
						result.OK = false
						return
					}

					// ok, 可以易位
					table.ClearPosition('e', 8)
					table.ClearPosition('h', 8)
					fromPiece.X = 'g'
					fromPiece.Y = 8
					fromPiece.Moved = true
					table.SetPosition(fromPiece)
					rookPiece.X = 'f'
					rookPiece.Y = 8
					fromPiece.Moved = true
					table.SetPosition(rookPiece)
					// 长 to = c
				} else {
					rookPiece := table.GetPosition('a', 8)

					// 挡住
					if hasChessBetweenTwoPointsInLine(table, 'e', 8, 'a', 8) {
						result.OK = false
						return
					}

					// to为nil或不为rook
					if toPiece == nil || toPiece.PieceType != chess.ChessPieceTypeRook {
						result.OK = false
						return
					}

					// 移动过
					if fromPiece.Moved || toPiece.Moved {
						result.OK = false
						return
					}

					// 检查路过的威胁
					testTable1 := table.Copy()
					testFromPiece := testTable1.GetPosition(fromX, fromY)
					testRookPiece := testTable1.GetPosition('a', 1)
					testTable1.ClearPosition(fromX, fromY)
					testTable1.ClearPosition(testRookPiece.X, testRookPiece.Y)
					testFromPiece.X = 'd'
					testFromPiece.Y = 8
					// 让车帮忙档一下, 比较易于算威胁
					testRookPiece.X = 'e'
					testRookPiece.Y = 8
					testTable1.SetPosition(testFromPiece)

					testTable2 := table.Copy()
					testFromPiece2 := testTable2.GetPosition(fromX, fromY)
					testRookPiece2 := testTable2.GetPosition('a', 8)
					testTable2.ClearPosition(fromX, fromY)
					testTable2.ClearPosition(testRookPiece2.X, testRookPiece2.Y)
					testFromPiece2.X = 'c'
					testFromPiece2.Y = 8
					testRookPiece2.X = 'd'
					testRookPiece2.Y = 8
					testTable2.SetPosition(testFromPiece2)
					testTable2.SetPosition(testRookPiece2)
					if checkPositionThreat(testTable1, side, 'd', 1) || checkPositionThreat(testTable2, side, 'c', 1) {
						result.OK = false
						return
					}

					// ok, 可以易位
					table.ClearPosition('e', 8)
					table.ClearPosition('a', 8)
					fromPiece.X = 'c'
					fromPiece.Y = 8
					fromPiece.Moved = true
					table.SetPosition(fromPiece)
					rookPiece.X = 'd'
					rookPiece.Y = 8
					fromPiece.Moved = true
					table.SetPosition(rookPiece)
				}
			}
			didKingRookSwitch = true
		}
	case chess.ChessPieceTypePawn:
		if side == chess.SideWhite {
			// 兵至少不能后退
			if toy <= fromy {
				result.OK = false
				return
			}

			diffY := toy - fromy
			diffX := tox - fromx

			// 共6中情况
			if diffX != 1 && diffX != -1 && diffX != 0 {
				result.OK = false
				return
			}
			if diffY != 2 && diffY != 1 {
				result.OK = false
				return
			}

			// 2种
			if diffY == 2 && diffX != 0 {
				result.OK = false
				return
				// 1种
			} else if diffY == 2 && diffX == 0 {
				if fromPiece.Moved {
					result.OK = false
					return
				}

				// 再判断是否有挡住的
				if hasChessBetweenTwoPointsInLine(table, fromX, fromY, toX, toY) {
					result.OK = false
					return
				}

				// 判断to是否有子
				if table.GetPosition(toX, toY) != nil {
					result.OK = false
					return
				}

				fromPiece.PawnMovedTwoLastTime = true

				// 2
			} else if diffY == 1 && diffX != 0 {
				// 必须斜着吃, 这里需要判断一下吃过路兵
				if toPiece == nil {
					if justMoveTwoPawn := findJustMoved2Pawn(table); justMoveTwoPawn == nil {
						result.OK = false
						return
					} else {
						justMoveTwoPawnX, justMoveTwoPawnY := chess.MustPositionToIndex(justMoveTwoPawn.X, justMoveTwoPawn.Y)
						if toy-justMoveTwoPawnY != 1 || tox != justMoveTwoPawnX {
							result.OK = false
							return
						}
						// 这里要吃掉过路兵
						table.ClearPosition(justMoveTwoPawn.X, justMoveTwoPawn.Y)
					}
				}
				// 1 diffY == 1 && diffX == 0
			} else {
				if toPiece != nil {
					result.OK = false
					return
				}
			}

			pawnUpgrade = toy == 7
		} else {
			if toy >= fromy {
				result.OK = false
				return
			}

			diffY := fromy - toy
			diffX := fromx - tox

			// 共6中情况
			if diffX != 1 && diffX != -1 && diffX != 0 {
				result.OK = false
				return
			}
			if diffY != 2 && diffY != 1 {
				result.OK = false
				return
			}

			// 2种
			if diffY == 2 && diffX != 0 {
				result.OK = false
				return
				// 1种
			} else if diffY == 2 && diffX == 0 {
				if fromPiece.Moved {
					result.OK = false
					return
				}

				// 再判断是否有挡住的
				if hasChessBetweenTwoPointsInLine(table, fromX, fromY, toX, toY) {
					result.OK = false
					return
				}

				// 判断to是否有子
				if table.GetPosition(toX, toY) != nil {
					result.OK = false
					return
				}

				fromPiece.PawnMovedTwoLastTime = true
				// 2
			} else if diffY == 1 && diffX != 0 {
				// 必须斜着吃, 这里需要判断一下吃过路兵
				if toPiece == nil {
					if justMoveTwoPawn := findJustMoved2Pawn(table); justMoveTwoPawn == nil {
						result.OK = false
						return
					} else {
						justMoveTwoPawnX, justMoveTwoPawnY := chess.MustPositionToIndex(justMoveTwoPawn.X, justMoveTwoPawn.Y)
						if justMoveTwoPawnY-toy != 1 || tox != justMoveTwoPawnX {
							result.OK = false
							return
						}
						// 这里要吃掉过路兵
						table.ClearPosition(justMoveTwoPawn.X, justMoveTwoPawnY)
					}
				}
				// 1 diffY == 1 && diffX == 0
			} else {
				if toPiece != nil {
					result.OK = false
					return
				}
			}

			pawnUpgrade = toy == 0
		}
	}

	if !didKingRookSwitch {
		// 移动棋盘的子
		table.ClearPosition(fromX, fromY)
		fromPiece.Moved = true
		fromPiece.X = toX
		fromPiece.Y = toY
		table.SetPosition(fromPiece)
	}

	// 设置一下justMoved
	for _, v := range findAllJustMoved2Pawn(table) {
		if v != fromPiece {
			v.PawnMovedTwoLastTime = false
		}
	}

	// 处理兵升变
	result.PawnUpgrade = pawnUpgrade

	// 找到对方的king
	remoteKing := findKing(table, remoteSide)

	// 是否将军
	kingThreat := checkPositionThreat(table, remoteSide, remoteKing.X, remoteKing.Y)

	// 王的8个单元格是否都受威胁
	kingAroundAllThreat := checkAround8Threat(table, remoteSide, remoteKing.X, remoteKing.Y)

	// 赢
	if kingThreat && kingAroundAllThreat {
		result.OK = true
		result.GameOver = true
		result.GameWinner = side
		return
	}

	// 判断和棋
	if !kingThreat && kingAroundAllThreat {
		result.OK = true
		result.GameOver = true
		result.GameWinner = chess.SideBoth
		return
	}

	result.OK = true
	result.GameOver = false
	result.KingThreat = kingThreat
	return
}
