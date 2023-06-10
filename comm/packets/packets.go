package packets

import (
	"chess-backend/comm/chess"
	"encoding/json"
)

type PacketTypeServerMoveRespType int

const (
	PacketTypeServerMoveRespTypeOK PacketTypeServerMoveRespType = iota
	PacketTypeServerMoveRespTypeFailed
	PacketTypeServerMoveRespTypePawnUpgrade
)

type PacketType int

const (
	// 心跳包
	PacketTypeHeartbeat PacketType = iota

	// 客户端要求开始匹配
	PacketTypeClientStartMatch

	// 服务端表示已经开始匹配
	PacketTypeServerMatching

	// 匹配完毕, 即将开始游戏
	PacketTypeServerMatchedOK

	// 客户端发送下棋的消息
	PacketTypeClientMove

	// 服务端告知用户下棋结果, 可能用户的输入不合法, 这里提示, 可能成功, 可能发生兵的升变, 要求用户继续输入
	PacketTypeServerMoveResp

	// 告知服务端兵升变成什么
	PacketTypeClientSendPawnUpgrade

	// 通知游戏结束
	PacketTypeServerGameOver

	// 通知对方掉线
	PacketTypeServerRemoteLoseConnection

	// 对方下棋下好了
	PacketTypeServerNotifyRemoteMove
)

type PacketHeader struct {
	Type *PacketType `json:"type"`
}

type PacketHeartbeat struct {
	PacketHeader
}

func (p *PacketHeartbeat) MustMarshalToBytes() []byte {
	i := PacketTypeHeartbeat
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketClientStartMatch struct {
	PacketHeader
}

func (p *PacketClientStartMatch) MustMarshalToBytes() []byte {
	i := PacketTypeClientStartMatch
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerMatching struct {
	PacketHeader
}

func (p *PacketServerMatching) MustMarshalToBytes() []byte {
	i := PacketTypeServerMatching
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerMatchedOK struct {
	PacketHeader
	Side  chess.Side        `json:"game_side"`
	Table *chess.ChessTable `json:"game_table"`
}

func (p *PacketServerMatchedOK) MustMarshalToBytes() []byte {
	i := PacketTypeServerMatchedOK
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketClientMove struct {
	PacketHeader
	FromX rune `json:"from_x"`
	FromY int  `json:"from_y"`
	ToX   rune `json:"to_x"`
	ToY   int  `json:"to_y"`
}

func (p *PacketClientMove) MustMarshalToBytes() []byte {
	i := PacketTypeClientMove
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketClientSendPawnUpgrade struct {
	PacketHeader
	ChessPieceType chess.ChessPieceType `json:"piece_type"`
}

func (p *PacketClientSendPawnUpgrade) MustMarshalToBytes() []byte {
	i := PacketTypeClientSendPawnUpgrade
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerMoveResp struct {
	PacketHeader
	MoveRespType PacketTypeServerMoveRespType `json:"resp_type"`
	// 下面的字段只有在状态OK的时候出现
	TableOnOK *chess.ChessTable `json:"table,omitempty"`
}

func (p *PacketServerMoveResp) MustMarshalToBytes() []byte {
	i := PacketTypeServerMoveResp
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerGameOver struct {
	PacketHeader
	Table      *chess.ChessTable `json:"final_table"`
	WinnerSide chess.Side        `json:"winner_side"`
}

func (p *PacketServerGameOver) MustMarshalToBytes() []byte {
	i := PacketTypeServerGameOver
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerRemoteLoseConnection struct {
	PacketHeader
}

func (p *PacketServerRemoteLoseConnection) MustMarshalToBytes() []byte {
	i := PacketTypeServerRemoteLoseConnection
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}

type PacketServerNotifyRemoteMove struct {
	PacketHeader
	Table             *chess.ChessTable `json:"table"`
	RemotePawnUpgrade bool              `json:"remote_pawn_upgrade"`
}

func (p *PacketServerNotifyRemoteMove) MustMarshalToBytes() []byte {
	i := PacketTypeServerNotifyRemoteMove
	p.Type = &i
	bs, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return bs
}
