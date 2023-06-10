package chess

import "chess-backend/comm/chess"

func DoUpgrade(table *chess.ChessTable, targetPieceType chess.ChessPieceType) {
	for _, v := range table {
		if v.GameSide == chess.SideWhite && v.Y == 8 {
			v.PieceType = targetPieceType
			return
		}

		if v.GameSide == chess.SideBlack && v.Y == 1 {
			v.PieceType = targetPieceType
			return
		}
	}
}
