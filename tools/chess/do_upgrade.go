package chess

import "chess-backend/comm/chess"

func DoUpgrade(table *chess.ChessTable, targetPieceType chess.ChessPieceType) {
	for _, v := range table {
		if v != nil && v.GameSide == chess.SideWhite && v.Y == 8 && v.PieceType == chess.ChessPieceTypePawn {
			v.PieceType = targetPieceType
			return
		}

		if v != nil && v.GameSide == chess.SideBlack && v.Y == 1 && v.PieceType == chess.ChessPieceTypePawn {
			v.PieceType = targetPieceType
			return
		}
	}
}
