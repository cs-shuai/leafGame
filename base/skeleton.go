package base

import (
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/module"
	"github.com/name5566/leaf/timer"
	"leafServer/conf"
	GameMsg "leafServer/msg/game"
	"leafServer/services"
)

func NewSkeleton() *module.Skeleton {
	skeleton := &module.Skeleton{
		GoLen:              conf.GoLen,
		TimerDispatcherLen: conf.TimerDispatcherLen,
		AsynCallLen:        conf.AsynCallLen,
		ChanRPCServer:      chanrpc.NewServer(conf.ChanRPCLen),
	}
	skeleton.Init()

	// 清除房间定时
	cronExpr, _ := timer.NewCronExpr("*/30 * * * *")
	skeleton.CronFunc(cronExpr, func() {
		if len(services.UserList) == 0 && len(services.RoomList) != 0 {
			services.RoomList = make(map[string]*GameMsg.Room)
		}
	})
	return skeleton
}
