package internal

import (
	"fmt"
	"github.com/name5566/leaf/gate"
	GameMsg "leafServer/msg/game"
	"leafServer/services"
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
}

func rpcNewAgent(args []interface{}) {
	fmt.Println("-------rpcNewAgent--------" + fmt.Sprint(len(services.UserList)) + "---------------")
	a := args[0].(gate.Agent)
	_ = a
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	user := a.UserData()
	fmt.Println("------rpcCloseAgent---------" + fmt.Sprint(user.(*GameMsg.User)) + "---------------")
	services.OutUser(user.(*GameMsg.User).No)
	fmt.Println("-------rpcCloseAgent--------" + fmt.Sprint(len(services.UserList)) + "---------------")

}
