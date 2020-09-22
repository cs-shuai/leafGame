package msg

import (
	"github.com/name5566/leaf/network/json"
	"leafServer/msg/game"
)

// var Processorbak network.Processor

// 使用默认的 JSON 消息处理器（默认还提供了 protobuf 消息处理器）
var Processor = json.NewProcessor()

func init() {
	// 这里我们注册了一个 JSON 消息 Hello
	Processor.Register(&GameMsg.User{})
	Processor.Register(&GameMsg.GameMessage{})
	Processor.Register(&GameMsg.Login{})
	Processor.Register(&GameMsg.Room{})
	Processor.Register(&GameMsg.Game{})
	Processor.Register(&GameMsg.Vote{})
	Processor.Register(&GameMsg.RoomOut{})
	Processor.Register(&GameMsg.Logout{})

}
