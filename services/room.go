package services

import (
	"errors"
	"fmt"
	"github.com/name5566/leaf/log"
	GameMsg "leafServer/msg/game"
	"math/rand"
	"strconv"
	"time"
)

func HandleRoom(args []interface{}) {
	fmt.Println("----------请注意我!!!房间处理中-----" + fmt.Sprint() + "---------------")
	// 接收参数处理
	i, a, user := getArgs(args)
	m := i.(*GameMsg.Room)

	if user.No == "" {
		Error(a, errors.New("账号未登录"))
		return
	}

	var roomId string
	var room *GameMsg.Room
	if m.RoomId == "" {
		if m.TotalNumber == "" || m.UndercoverNumber == "" {
			Error(a, errors.New("操作异常"))
			return
		}
		room = CreateRoom(user, m)
		roomId = room.RoomId
	} else {
		roomId = m.RoomId
		room = getRoom(roomId)
	}

	// 是否准备
	if m.IsPrepare {
		b, err := Prepare(user, m)
		if err != nil {
			Error(a, err)
		}

		var str = "取消准备"
		if b {
			str = "已准备"
		}
		RoomSendMessage(room, user.Name+": "+str)
		Success(a, "", "Prepare", nil)
		return
	}

	room, err := JoinRoom(user, roomId)
	if err != nil {
		Error(a, err)
		return
	}

	var result = make(map[string]interface{})
	result["RoomInfo"] = room
	Success(a, "进入房间: "+roomId, "Room", result)
}

/**
 * 准备
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
func Prepare(user *GameMsg.User, room *GameMsg.Room) (b bool, err error) {
	var roomInfo = getRoom(room.RoomId)
	if roomInfo == nil {
		return b, errors.New("玩家为进入房间")
	}
	if _, ok := roomInfo.UserList[user.No]; !ok {
		return b, errors.New("未知错误")
	}

	if _, ok := roomInfo.PrepareList[user.No]; ok {
		delete(roomInfo.PrepareList, user.No)
		roomInfo.PrepareNum--
	} else {
		roomInfo.PrepareList[user.No] = user.No
		roomInfo.PrepareNum++
		b = true
	}

	matchPrepare(roomInfo, user)

	return b, nil
}
func CreateRoom(user *GameMsg.User, room *GameMsg.Room) *GameMsg.Room {
	var roomId string
	if room.RoomId == "" {
		for {
			roomId = MakeRoomCode()
			if getRoom(roomId) == nil {
				break
			}
		}
		room.RoomId = roomId
	} else {
		roomId = room.RoomId
	}

	if room.CreateUserId == "" {
		room.CreateUserId = user.No
		// common.Db.Where("room_id = ?", roomId).First(room)
		// if room.Id == 0 {
		// 	common.Db.Create(room)
		// }
		room.CreateUser = user
	}
	room.MsgChan = make(chan string, 10)
	room.PrepareList = make(map[string]string)
	go func(room *GameMsg.Room) {
		for {
			select {
			case msg := <-room.MsgChan:
				RoomSendMessage(room, msg)
			}
		}
	}(room)

	room.UserList = make(map[string]*GameMsg.User)

	// 添加房间到房间列表
	addRoom(room)

	return room
}

func JoinRoom(user *GameMsg.User, roomId string) (room *GameMsg.Room, err error) {
	if getRoom(roomId) == nil {
		return room, errors.New("房间不存在")
	}
	room = getRoom(roomId)
	totalNumber, err := strconv.Atoi(room.TotalNumber)
	if err != nil || room == nil {
		return room, err
	}

	// 验证房间人数
	if _, ok := room.UserList[user.No]; !ok && len(room.UserList) >= totalNumber {
		return room, errors.New("房间已满")
	}

	// 重置房主
	if user.No == room.CreateUserId {
		room.CreateUser = user
		matchPrepare(room, user)
	}

	// 取消准备
	if _, ok := room.PrepareList[user.No]; ok {
		delete(room.PrepareList, user.No)
		room.PrepareNum--
	}

	user.RoomId = roomId
	// common.Db.Model(user).Where("no = ?", user.No).UpdateColumn("room_id", roomId)
	if _, ok := room.UserList[user.No]; !ok {
		room.Number++
	}
	room.UserList[user.No] = user
	room.MsgChan <- user.Name + " 进入房间"
	log.Debug("房间: %s 当前人数: %d", room.RoomId, room.Number)
	return room, err
}

/**
 * 用户进入房间
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func UserJoinRoom(user *GameMsg.User, roomId string) {
	// 存在房间ID 房间列表中不存在
	if user.RoomId != "" && getRoom(roomId) == nil {
		CreateRoomForDb(user)
	}

	// 进入房间
	room, err := JoinRoom(user, roomId)
	log.Error("进入房间错误 %v", err)
	if err != nil {
		Error(*user.Agent, err)
		return
	}

	var result = make(map[string]interface{})
	result["RoomInfo"] = room
	// 如果游戏开始
	if room.GameInfo != nil {
		keywordMap := getUserKeyword(user.No, room)
		Success(*user.Agent, "游戏开始", "StartGame", keywordMap)

		time.Sleep(1 * time.Second)
		// 投票中
		if _, ok := room.GameInfo.VoteList[user.No]; !ok && room.GameInfo.Stage == GameMsg.GameStage_Vote {
			Success(*user.Agent, "开始投票", "Vote", room.GameInfo)
		}

		return
	}

	Success(*user.Agent, "进入房间: "+user.RoomId, "Room", result)
}

/**
 * 创建房间从数据库
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func CreateRoomForDb(user *GameMsg.User) {
	// 房间是否存在
	if getRoom(user.RoomId) == nil {
		NewRoom := new(GameMsg.Room)
		// common.Db.Where("room_id = ?", user.RoomId).First(NewRoom)
		fmt.Println("-----NewRoom----------" + fmt.Sprint(NewRoom) + "---------------")
		CreateRoom(user, NewRoom)
	}
}

/**
 * 生成房间号
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func MakeRoomCode() string {
	rand.Seed(time.Now().Unix())
	r := rand.Intn(999999)

	s := fmt.Sprintf("%0*d", 6, r)

	return s
}

/**
 * 发布房间消息
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func RoomSendMessage(room *GameMsg.Room, msg string) {
	fmt.Println("---------------" + fmt.Sprint(room.UserList) + "---------------")
	for _, userInfo := range room.UserList {
		Success(*userInfo.Agent, msg, "Messagee", nil)
	}
}

/**
 * 用户发送消息
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func UserSendMessage(room *GameMsg.Room, msg string, user *GameMsg.User) {
	msg = fmt.Sprintf("%s: %s", user.Name, msg)
	RoomSendMessage(room, msg)
}

/**
 * 退出房间
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func RoomOut(args []interface{}) {
	// 接收参数处理
	_, a, user := getArgs(args)
	room := getRoom(user.RoomId)
	fmt.Println("---------user------" + fmt.Sprint(user) + "---------------")
	delete(room.UserList, user.No)
	user.RoomId = ""
	room.Number--
	// common.Db.Model(room).Where("room_id = ?", room.RoomId).UpdateColumn("number", room.Number)
	// common.Db.Model(user).Where("no = ?", user.No).UpdateColumn("room_id", "")
	RoomSendMessage(room, user.Name+"退出房间")
	Success(a, "退出房间", "RoomOut", nil)
}

/**
 * 重置准备
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func resetPrepare(room *GameMsg.Room) {
	room.PrepareNum = 0
	room.PrepareList = map[string]string{}
}

/**
 * 计算投票
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
func matchPrepare(room *GameMsg.Room, user *GameMsg.User) {
	totalNumber, _ := strconv.Atoi(room.TotalNumber)
	// fmt.Println("--------PrepareNum-------" + fmt.Sprint(room.PrepareNum) + "---------------")
	// fmt.Println("--------totalNumber-------" + fmt.Sprint(totalNumber) + "---------------")
	// fmt.Println("--------room.Number-------" + fmt.Sprint(room.Number) + "---------------")
	if room.PrepareNum == totalNumber-1 && totalNumber == room.Number {
		Success(*room.CreateUser.Agent, "开始", "Start", nil)
	} else {
		if room.CreateUser != nil {
			Success(*room.CreateUser.Agent, "准备", "Start", nil)
		}
	}
}

/**
 * 通过房间号获取房间信息
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func getRoom(roomId string) *GameMsg.Room {
	if _, ok := RoomList[roomId]; ok {
		return RoomList[roomId]
	}

	return nil
}

/**
 * 添加房间到房间列表
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func addRoom(room *GameMsg.Room) {
	RoomList[room.RoomId] = room
}
