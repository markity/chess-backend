package main

import (
	"chess-backend/comm/settings"
	"chess-backend/game"
	"fmt"
	"runtime"
	"time"

	"github.com/Allenxuxu/gev"
)

func main() {
	server, err := gev.NewServer(&game.ConnHandler{},
		gev.Address(fmt.Sprintf("%s:%d", settings.ServerListenIP, settings.ServerListenPort)),
		gev.Network("tcp"),
		gev.LoadBalance(gev.RoundRobin()),
		gev.NumLoops(runtime.NumCPU()),
	)
	if err != nil {
		panic(err)
	}

	// 这个是全局的心跳检测器
	server.RunEvery(time.Millisecond*settings.HeartbeatInterval, game.OnTimeout)

	server.Start()
}
