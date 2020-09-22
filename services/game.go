package services

import (
	"fmt"
	"github.com/name5566/leaf/log"
	GameMsg "leafServer/msg/game"
	"strconv"
)

func HandleGame(args []interface{}) {
	fmt.Println("----------请注意我!!!游戏开始-----" + fmt.Sprint() + "---------------")
	// 接收参数处理
	m, _, user := getArgs(args)
	// 获取房间信息
	roomInfo := getRoom(user.RoomId)

	switch m.(*GameMsg.Game).Stage {
	case GameMsg.GameStage_Start:
		Start(roomInfo)
		break
	case GameMsg.GameStage_Vote:
		Vote(roomInfo)
		break
	}

}

func Vote(roomInfo *GameMsg.Room) {
	// 创建投票通道 长度为剩余人数
	roomInfo.GameInfo.VoteChan = make(chan *GameMsg.Vote, len(roomInfo.GameInfo.SurvivalUserList))
	roomInfo.GameInfo.VoteList = make(map[string]*GameMsg.Vote)
	roomInfo.GameInfo.VoteNum++

	// 等待到时统计投票结果
	go VoteOver(roomInfo, roomInfo.GameInfo.VoteTime)
	// 监听投票 全部投票结束 统计投票结果
	go voteListen(roomInfo)

	// 发布投票消息
	roomInfo.GameInfo.Stage = GameMsg.GameStage_Vote
	for userNo, _ := range roomInfo.GameInfo.SurvivalUserList {
		userInfo := roomInfo.UserList[userNo]
		Success(*userInfo.Agent, "开始投票", "Vote", roomInfo.GameInfo)
	}
}

func voteListen(roomInfo *GameMsg.Room) {
	for {
		select {
		case vote := <-roomInfo.GameInfo.VoteChan:
			log.Debug("有人投票 %v", vote)
			roomInfo.GameInfo.VoteList[vote.UserId] = vote
			if len(roomInfo.GameInfo.VoteList) >= len(roomInfo.GameInfo.SurvivalUserList) {
				log.Debug("完成")
				go VoteOver(roomInfo, 0)
			}
		}
	}
}

func Start(roomInfo *GameMsg.Room) {
	// 初始化游戏
	initGame(roomInfo)

	// 通知玩家
	for _, userInfo := range roomInfo.UserList {
		// 获取词组
		keywordMap := getUserKeyword(userInfo.No, roomInfo)
		Success(*userInfo.Agent, "游戏开始", "StartGame", keywordMap)
	}
}

/**
 * 获取用户词组
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
func getUserKeyword(userNo string, roomInfo *GameMsg.Room) (res keywordResult) {
	res.Stage = "start"
	res.CreateUserId = roomInfo.CreateUserId
	if roomInfo.GameInfo.SurvivalUserList[userNo].Role == GameMsg.Role_Normal {
		res.Keyword = roomInfo.GameInfo.Keyword.NormalWord
	} else {
		res.Keyword = roomInfo.GameInfo.Keyword.UndercoverWord
	}

	return res
}

/**
 * 创建新游戏
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
func initGame(roomInfo *GameMsg.Room) {
	game := new(GameMsg.Game)
	game.RoomId = roomInfo.RoomId
	game.UndercoverNum, _ = strconv.Atoi(roomInfo.UndercoverNumber)
	game.Round = 1
	game.ActionTime = 60
	game.VoteTime = 60
	game.VoteNum = 1
	game.Stage = "start"

	roomInfo.GameInfo = game

	// 生成词组
	GetKeyword(roomInfo)

	// 分配角色
	distributionRole(roomInfo)

	// 重置投票
	resetPrepare(roomInfo)
}

func GetKeyword(roomInfo *GameMsg.Room) {
	for {
		keywords := getRandomDictionary()
		key := fmt.Sprintf("%s-%s", keywords[0], keywords[1])
		if _, ok := roomInfo.UsedWord[key]; !ok {
			keyword := new(GameMsg.Keyword)
			keyword.NormalWord = keywords[0]
			keyword.UndercoverWord = keywords[1]
			roomInfo.GameInfo.Keyword = keyword
			break
		}
	}
}

/**
 *
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
func distributionRole(roomInfo *GameMsg.Room) {
	var UndercoverNum = 0
	for _, userInfo := range roomInfo.UserList {
		if UndercoverNum < roomInfo.GameInfo.UndercoverNum {
			userInfo.Role = GameMsg.Role_Undercover
			UndercoverNum++
		} else {
			userInfo.Role = GameMsg.Role_Normal
		}
	}

	roomInfo.GameInfo.SurvivalUserList = make(map[string]*GameMsg.User)
	// 克隆用户列表
	for userNo, userInfo := range roomInfo.UserList {
		roomInfo.GameInfo.SurvivalUserList[userNo] = userInfo
	}
}
