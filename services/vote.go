package services

import (
	"errors"
	"fmt"
	"github.com/name5566/leaf/log"
	GameMsg "leafServer/msg/game"
	"time"
)

func HandleVote(args []interface{}) {
	fmt.Println("----------请注意我!!!投票-----" + fmt.Sprint() + "---------------")
	// 接收参数处理
	i, a, user := getArgs(args)
	m := i.(*GameMsg.Vote)

	roomInfo := getRoom(user.RoomId)
	if roomInfo.GameInfo.Stage != GameMsg.GameStage_Vote {
		Error(a, errors.New("投票错误"))
		return
	}

	roomInfo.GameInfo.VoteChan <- m

	Success(a, "投票成功", "VoteSuccess", nil)
}

/**
 * 投票结束
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func VoteOver(roomInfo *GameMsg.Room, sleepTime int) {
	fmt.Println("--------VoteOver-------" + fmt.Sprint() + "---------------")
	// 等待投票结果时间设置
	time.Sleep(time.Duration(sleepTime) * time.Second)
	// 阶段不正确跳出
	if roomInfo.GameInfo.Stage != GameMsg.GameStage_Vote {
		return
	}

	// 投票积分
	var voteResult = make(map[string]int)
	for _, vote := range roomInfo.GameInfo.VoteList {
		if _, ok := voteResult[vote.VotePlayerNumber]; !ok {
			voteResult[vote.VotePlayerNumber] = 1
		} else {
			voteResult[vote.VotePlayerNumber]++
		}
	}

	log.Debug("投票结果:", voteResult)
	var highestScore = 0
	for userNo, score := range voteResult {
		// 平分
		if score == highestScore {
			roomInfo.GameInfo.OutUser = append(roomInfo.GameInfo.OutUser, roomInfo.UserList[userNo])
		}

		// 高于最高分
		if score > highestScore {
			highestScore = score
			roomInfo.GameInfo.OutUser = []*GameMsg.User{roomInfo.UserList[userNo]}
		}
	}

	// 最高分人数 不唯一 或 没人投票
	if len(roomInfo.GameInfo.OutUser) > 1 || highestScore == 0 {
		log.Debug("投票最高分:", highestScore)
		log.Debug("投票最高用户:", roomInfo.GameInfo.OutUser)

		// 平票次数过多
		if roomInfo.GameInfo.VoteNum > 3 {
			RoomSendMessage(roomInfo, "多次平票,无人淘汰")
			return
		}

		// 重新投票
		Vote(roomInfo)
		return
	}

	if roomInfo.GameInfo.OutUser[0] != nil {
		// 淘汰最高分
		delete(roomInfo.GameInfo.SurvivalUserList, roomInfo.GameInfo.OutUser[0].No)
	}

	// 计算游戏结果
	accounting(roomInfo)

	// 发送结果
	sendVoteOver(roomInfo)

}

func sendVoteOver(roomInfo *GameMsg.Room) {
	var msg = "游戏继续"
	if roomInfo.GameInfo.Stage == GameMsg.GameStage_Over {
		if roomInfo.GameInfo.WinRole == GameMsg.Role_Normal {
			msg = "好人胜利"
		} else {
			msg = "卧底胜利"
		}
	}

	// 发送投票结果
	sendVoteResult(roomInfo)

	var message = fmt.Sprintf("%s 淘汰, %s", roomInfo.GameInfo.OutUser[0].Name, msg)
	fmt.Println("---------roomInfo.UserList------" + fmt.Sprint(roomInfo.UserList) + "---------------")
	for _, userInfo := range roomInfo.UserList {
		Success(*userInfo.Agent, message, roomInfo.GameInfo.Stage, roomInfo.GameInfo)
	}
}

/**
 * 发送返回结果
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func sendVoteResult(roomInfo *GameMsg.Room) {
	result := fmt.Sprintf("第%d回合投票结果: <br>", roomInfo.GameInfo.Round)
	for _, vote := range roomInfo.GameInfo.VoteList {
		result += fmt.Sprintf("%s : 投票-> %s <br>", getUser(vote.UserId).Name, getUser(vote.VotePlayerNumber).Name)
	}

	RoomSendMessage(roomInfo, result)
}

func accounting(roomInfo *GameMsg.Room) {
	gameInfo := roomInfo.GameInfo
	// 游戏阶段为继续游戏
	gameInfo.Stage = GameMsg.GameStage_Game
	gameInfo.Round++

	// 计算人数
	if len(gameInfo.SurvivalUserList) <= gameInfo.UndercoverNum+1 {
		// 卧底胜利
		gameInfo.Stage = GameMsg.GameStage_Over
		gameInfo.WinRole = GameMsg.Role_Undercover
	}

	var UndercouverNum int
	for _, userInfo := range gameInfo.SurvivalUserList {
		if userInfo.Role == GameMsg.Role_Undercover {
			UndercouverNum++
		}
	}

	// 不存在卧底
	if UndercouverNum == 0 {
		// 好人胜利
		gameInfo.Stage = GameMsg.GameStage_Over
		gameInfo.WinRole = GameMsg.Role_Normal
	}
}
