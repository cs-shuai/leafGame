package gate

import (
	"leafServer/game"
	"leafServer/msg"
	"leafServer/msg/game"
)

func init() {
	// 这里指定消息 Hello 路由到 game 模块
	// 模块间使用 ChanRPC 通讯，消息路由也不例外
	msg.Processor.SetRouter(&GameMsg.GameMessage{}, game.ChanRPC)
	msg.Processor.SetRouter(&GameMsg.User{}, game.ChanRPC)
	msg.Processor.SetRouter(&GameMsg.Login{}, game.ChanRPC)
	msg.Processor.SetRouter(&GameMsg.Room{}, game.ChanRPC)
	msg.Processor.SetRouter(&GameMsg.RoomOut{}, game.ChanRPC)
	msg.Processor.SetRouter(&GameMsg.Game{}, game.ChanRPC)
	msg.Processor.SetRouter(&GameMsg.Vote{}, game.ChanRPC)
	msg.Processor.SetRouter(&GameMsg.Logout{}, game.ChanRPC)
}
