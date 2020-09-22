package services

import (
	"crypto/sha1"
	"fmt"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	GameMsg "leafServer/msg/game"
	"time"
)

func HandleLogin(args []interface{}) {
	fmt.Println("----------请注意我!!!我要登录了-----" + fmt.Sprint() + "---------------")
	// 接收参数处理
	i, a, user := getArgs(args)
	m := i.(*GameMsg.Login)

	fmt.Println("---------------" + fmt.Sprint(i) + "---------------")
	fmt.Println("---------------" + fmt.Sprint(m) + "---------------")
	// 会员标识
	var userNo = m.UserId
	if userNo == "" {
		userNo = MakeUserId()
	}

	// 获取用户信息
	user = GetUser(userNo, m.UserName)

	// 用户登录
	loginUser(user, a)

	// 进入房间
	if user.RoomId != "" {
		UserJoinRoom(user, user.RoomId)
	}

	Success(a, "登录成功", "Login", GameMsg.Login{
		UserId:   user.No,
		UserName: user.Name,
	})
}

/**
 * 获取用户
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func GetUser(userNo, userName string) *GameMsg.User {
	user := new(GameMsg.User)
	// 存在返回
	if getUser(userNo) != nil {
		return getUser(userNo)
	}

	// 数据库查询
	// common.Db.Where("no = ?", userNo).First(user)
	// if  user.Id  != 0 {
	// 	return user
	// }

	// 创建用户
	user = createUser(userNo, userName)

	return user
}

/**
 * 创建会员
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func createUser(userNo, userName string) *GameMsg.User {
	user := new(GameMsg.User)
	fmt.Println("------userName---------" + fmt.Sprint(userName) + "---------------")
	if userName == "" {
		userName = MakeRoomCode() + "玩家"
	}
	user.Name = userName
	user.No = userNo
	// common.Db.Create(user)

	return user
}

/**
 * 用户登录
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func loginUser(user *GameMsg.User, agent gate.Agent) {
	user.Agent = &agent
	agent.SetUserData(user)

	addUser(user)
}

/**
 * 生成用户编码
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
func MakeUserId() string {
	n := time.Now().UnixNano()

	h := sha1.New()
	h.Write([]byte(fmt.Sprint(n)))
	l := fmt.Sprintf("%x", h.Sum(nil))
	return l
}

/**
 * 获取用户从用户列表
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func getUser(userNo string) *GameMsg.User {
	if _, ok := UserList[userNo]; ok {
		return UserList[userNo]
	}

	return nil
}

/**
 * 添加用户到用户列表
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func addUser(user *GameMsg.User) {
	UserList[user.No] = user
	log.Debug("当前在线人数: %d", len(UserList))
}

/**
 * 获取用户从用户列表
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
func OutUser(userNo string) {
	delete(UserList, userNo)
}

func HanldeLogout(args []interface{}) {
	// 接收参数处理
	_, _, user := getArgs(args)
	OutUser(user.No)
	// common.Db.Delete(user.No)
}
