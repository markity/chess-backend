package game

import (
	"chess-backend/comm/chess"
	"chess-backend/comm/packets"
	"chess-backend/comm/settings"

	chesstool "chess-backend/tools/chess"
	othertool "chess-backend/tools/other"
	packtool "chess-backend/tools/packet"

	"github.com/Allenxuxu/gev"
)

type ConnState int

const (
	// 连接上服务端的初始状态状态
	ConnStateNone ConnState = iota
	ConnStateMatching
	ConnStateGaming
)

type GameState int

const (
	GameStateWaitingWhitePut GameState = iota
	GameStateWaitingBlackPut
	// 等待黑方的兵升迁
	GameStateWaitingBlackUpgrade
	// 等待白方的兵升迁
	GameStateWaitingWhiteUpgrade
)

type ConnHandler struct{}

func (ch *ConnHandler) OnConnect(c *gev.Connection) {
	connID := int(AtomicIDIncrease.Add(1))
	connCtx := &ConnContext{ID: int(connID), LoseHertbeatCount: 0, Conn: c, ConnState: ConnStateNone, Gcontext: nil}

	ConnMapLock.Lock()
	ConnMap[connID] = connCtx
	ConnMapLock.Unlock()

	c.SetContext(connID)
}

func (ch *ConnHandler) OnClose(c *gev.Connection) {
	connID := c.Context().(int)
	ConnMapLock.Lock()
	if ConnMap[connID].ConnState == ConnStateGaming {
		var remoteConnContext *ConnContext
		if ConnMap[connID].Gcontext.BlackConnContext.ID == connID {
			remoteConnContext = ConnMap[connID].Gcontext.WhiteConnContext
		} else {
			remoteConnContext = ConnMap[connID].Gcontext.BlackConnContext
		}
		packet := packets.PacketServerRemoteLoseConnection{}
		packetBytesWithHeader := packtool.DoPackWith4BytesHeader(packet.MustMarshalToBytes())
		remoteConnContext.Conn.Send(packetBytesWithHeader)
		remoteConnContext.ConnState = ConnStateNone
	}
	delete(ConnMap, connID)
	ConnMapLock.Unlock()
}

func (ch *ConnHandler) OnMessage(c *gev.Connection, ctx interface{}, data []byte) interface{} {
	// 没有收到消息
	if data == nil {
		return nil
	}

	connID := c.Context().(int)

	packIface := packets.ServerParse(data)

	ConnMapLock.Lock()
	defer ConnMapLock.Unlock()
	switch packet := packIface.(type) {
	case *packets.PacketHeartbeat:
		ConnMap[connID].LoseHertbeatCount = 0
	case *packets.PacketClientStartMatch:
		// 协议错误
		if ConnMap[connID].ConnState != ConnStateNone {
			c.Close()
		}

		matched := false
		// 找一个正在match的连接
		for _, v := range ConnMap {
			if v.ID != connID && v.ConnState == ConnStateMatching {
				// 随机摇game side
				var whiteConnContext *ConnContext
				var blackConnContext *ConnContext
				if othertool.RandGetBool() {
					whiteConnContext = ConnMap[connID]
					blackConnContext = ConnMap[v.ID]
				} else {
					blackConnContext = ConnMap[connID]
					whiteConnContext = ConnMap[v.ID]
				}

				table := chess.NewChessTable()
				gameContext := GameContext{
					WhiteConnContext: whiteConnContext,
					BlackConnContext: blackConnContext,
					Gstate:           GameStateWaitingWhitePut,
					Table:            table,
				}
				gameContext.WhiteConnContext = whiteConnContext
				gameContext.BlackConnContext = blackConnContext

				matchingPacket := packets.PacketServerMatching{}
				matchingPacketWithHeader := packtool.DoPackWith4BytesHeader(matchingPacket.MustMarshalToBytes())

				packetForBlack := packets.PacketServerMatchedOK{Side: chess.SideBlack, Table: table}
				packetForBlackBytesWithHeader := packtool.DoPackWith4BytesHeader(packetForBlack.MustMarshalToBytes())
				v.ConnState = ConnStateGaming
				v.Gcontext = &gameContext

				v.Conn.Send(matchingPacketWithHeader)
				v.Conn.Send(packetForBlackBytesWithHeader)

				packetForWhite := packets.PacketServerMatchedOK{Side: chess.SideWhite, Table: table}
				packetForWhiteBytesWithHeader := packtool.DoPackWith4BytesHeader(packetForWhite.MustMarshalToBytes())
				ConnMap[connID].ConnState = ConnStateGaming
				ConnMap[connID].Gcontext = &gameContext

				ConnMap[connID].Conn.Send(matchingPacketWithHeader)
				ConnMap[connID].Conn.Send(packetForWhiteBytesWithHeader)
				matched = true
				break
			}
		}

		// 找不到一个匹配的, 那么标记为正在匹配
		if !matched {
			ConnMap[connID].ConnState = ConnStateMatching
			retPacket := packets.PacketServerMatching{}
			retPacketBytesWithHeader := packtool.DoPackWith4BytesHeader(retPacket.MustMarshalToBytes())
			ConnMap[connID].Conn.Send(retPacketBytesWithHeader)
		}
	case *packets.PacketClientMove:
		// 协议判断
		if ConnMap[connID].ConnState != ConnStateGaming {
			ConnMap[connID].Conn.Close()
			return nil
		}

		// 拿到一些信息
		var gameContext = ConnMap[connID].Gcontext
		var selfContext *ConnContext = ConnMap[connID]
		var selfSide chess.Side
		var remoteContext *ConnContext
		var remoteSide chess.Side
		othertool.Ignore(remoteContext)
		othertool.Ignore(remoteSide)
		if gameContext.BlackConnContext == selfContext {
			remoteContext = gameContext.WhiteConnContext
			selfSide = chess.SideBlack
			remoteSide = chess.SideWhite
		} else {
			remoteContext = gameContext.BlackConnContext
			selfSide = chess.SideWhite
			remoteSide = chess.SideBlack
		}

		// 协议判断, 要求发送方确实是下棋的一方
		if (selfSide == chess.SideBlack && gameContext.Gstate != GameStateWaitingBlackPut) ||
			(selfSide == chess.SideWhite && gameContext.Gstate != GameStateWaitingWhitePut) {
			ConnMap[connID].Conn.Close()
			return nil
		}

		// 协议判断, 输入格式判断, 要求输入格式确实正确
		// 注意x,y两两相等的情况也是不合法的
		if !chesstool.CheckChessPostsionVaild(packet.FromX, packet.FromY) ||
			!chesstool.CheckChessPostsionVaild(packet.ToX, packet.ToY) ||
			(packet.FromX == packet.ToX && packet.FromY == packet.ToY) {
			ConnMap[connID].Conn.Close()
			return nil
		}

		result := chesstool.DoMove(gameContext.Table, selfSide, packet.FromX, packet.FromY, packet.ToX, packet.ToY)
		if !result.OK {
			moveFailedPacket := packets.PacketServerMoveResp{
				MoveRespType: packets.PacketTypeServerMoveRespTypeFailed,
				TableOnOK:    nil,
			}
			moveFailedPacketBytesWithHeader := packtool.DoPackWith4BytesHeader(moveFailedPacket.MustMarshalToBytes())
			selfContext.Conn.Send(moveFailedPacketBytesWithHeader)
			return nil
		}

		if !result.GameOver {
			// 处理兵的升变问题
			moveRespType := packets.PacketTypeServerMoveRespTypeOK
			if result.PawnUpgrade {
				moveRespType = packets.PacketTypeServerMoveRespTypePawnUpgrade
			}

			// move ok, 做通知
			moveOKPacket := packets.PacketServerMoveResp{
				MoveRespType: moveRespType,
				TableOnOK:    gameContext.Table,
			}
			moveOKPacketBytesWithHeader := packtool.DoPackWith4BytesHeader(moveOKPacket.MustMarshalToBytes())
			selfContext.Conn.Send(moveOKPacketBytesWithHeader)

			remoteMovePacket := packets.PacketServerNotifyRemoteMove{
				Table:             gameContext.Table,
				RemotePawnUpgrade: moveRespType == packets.PacketTypeServerMoveRespTypePawnUpgrade,
			}
			remoteContext.Conn.Send(remoteMovePacket)

			return nil
		}

		// game over, 发送消息, 清空资源
		gameOverPacket := packets.PacketServerGameOver{
			Table:      gameContext.Table,
			WinnerSide: selfSide,
		}
		gameOverPacketBytesWithHeader := packtool.DoPackWith4BytesHeader(gameOverPacket.MustMarshalToBytes())
		selfContext.Conn.Send(gameOverPacketBytesWithHeader)
		remoteContext.Conn.Send(gameOverPacketBytesWithHeader)

		selfContext.Gcontext = nil
		selfContext.ConnState = ConnStateNone
		remoteContext.Gcontext = nil
		remoteContext.ConnState = ConnStateNone

		return nil
	case *packets.PacketClientSendPawnUpgrade:
		othertool.Ignore(nil)
	case nil:
		// 协议错误, 直接关闭
		c.Close()
	}
	return nil
}

func OnTimeout() {
	var packet = packets.PacketHeartbeat{}
	heartPacketBytesWithHeader := packtool.DoPackWith4BytesHeader(packet.MustMarshalToBytes())

	ConnMapLock.Lock()
	for k := range ConnMap {
		ConnMap[k].Conn.Send(heartPacketBytesWithHeader)
		ConnMap[k].LoseHertbeatCount++
		if ConnMap[k].LoseHertbeatCount >= settings.MaxLoseHeartbeat {
			ConnMap[k].Conn.Close()
		}
	}
	ConnMapLock.Unlock()
}
