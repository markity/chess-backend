# 象棋游戏-后端

### 1.接口协议:

tcp分包, 前4字节指示分包的长度, 紧接着的是分包的字节。

### 2. 分包类型:

- PacketTypeHeartbeat: 心跳包, 客户端服务端维持200ms的心跳, 5次丢失算做断线, 此时客户端/服务端自动断开连接
- PacketTypeClientStartMatch: 客户端要求开始匹配
- PacketTypeServerMatchedOK: 服务端告知用户匹配完毕
- PacketTypeClientMove: 客户端告知自己的下棋动作, 包括两个坐标和是否仪和
- PacketTypeServerMoveResp: 服务端告知客户端上个动作的结果, 比如不合法的移动, 或者现在有兵的升变
- PacketTypeClientSendPawnUpgrade: 客户端告知服务端自己的兵想要升变成什么
- PacketTypeServerGameOver: 服务端告知客户端游戏结束, 有四种可能, 投降, 平局, 对方认输, 正常分出胜负
- PacketTypeServerRemoteLoseConnection: 如果对方断线, 这个用来告知游戏者对端连接断开
- PacketTypeServerNotifyRemoteMove: 告知游戏者对方的动作, 包括两个坐标和对方是否仪和, 或者对方是否正在进行兵的升变
- PacketTypeClientWheatherAcceptDraw: 如果对方要求和棋, 客户端发送这个包来确认是否同意和棋
- PacketTypeClientDoSurrender: 主动认输
- PacketTypeServerRemoteUpgradeOK: 告知对方的兵的升变已经完成
- PacketTypeServerUpgradeOK: 告知服务端自己兵应该升变成什么

### 3. 游戏玩法

```plaintext
mov a2 a3                                 移动
dmov a2 a3                               移动并提出议和
accept                                       接受对方的议和
refuse                                        拒绝对方的议和
swi bishop/knight/rook/queen  进行一个兵的升变
sur                                             直接投降
```

### 4. 说明

我已经开启服务器, exe文件在邮件的压缩包里面, 打开就直接能玩。