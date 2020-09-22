package internal

import (
	GameMsg "leafServer/msg/game"
	"leafServer/services"
	"reflect"
)

func init() {
	// 向当前模块（game 模块）注册 Hello 消息的消息处理函数 handleHello
	handler(&GameMsg.Game{}, services.HandleGame)
	handler(&GameMsg.Login{}, services.HandleLogin)
	handler(&GameMsg.Room{}, services.HandleRoom)
	handler(&GameMsg.GameMessage{}, services.HandleMessage)
	handler(&GameMsg.Vote{}, services.HandleVote)
	handler(&GameMsg.RoomOut{}, services.RoomOut)
	handler(&GameMsg.Logout{}, services.HanldeLogout)
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}
