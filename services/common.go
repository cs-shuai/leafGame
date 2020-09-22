package services

import (
	"errors"
	"fmt"
	"github.com/name5566/leaf/gate"
	"io/ioutil"
	GameMsg "leafServer/msg/game"
	"math/rand"
	"strings"
)

var Dictionary [][]string

func init() {
	// 初始化字典
	Dictionary = GetDictionary()
}

/**
 * 词组返回
 * @Author: cs_shuai
 * @Date: 2020-09-21
 */
type keywordResult struct {
	Keyword      string
	Stage        string
	CreateUserId string
}

/**
 * 获取字典
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
func getRandomDictionary() []string {
	index := rand.Intn(len(Dictionary) - 1)
	return Dictionary[index]
}

/**
 * 用户列表
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
var UserList = make(map[string]*GameMsg.User)

/**
 * 房间列表
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
var RoomList = make(map[string]*GameMsg.Room)

/**
 * 错误提示
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
func Error(a gate.Agent, err error) {
	a.WriteMsg(&GameMsg.GameMessage{
		Msg:    fmt.Sprint(err),
		Status: 0,
	})
}

/**
 * 成功提示
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
func Success(a gate.Agent, msg string, t string, data interface{}) {
	a.WriteMsg(&GameMsg.GameMessage{
		Msg:    msg,
		Data:   data,
		Type:   t,
		Status: 1,
	})
}

/**
 * 发送消息
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
func HandleMessage(args []interface{}) {
	// 接收参数处理
	i, a, user := getArgs(args)
	m := i.(*GameMsg.GameMessage)

	// 验证房间信息
	if m.RoomId == "" || getRoom(m.RoomId) == nil {
		Error(a, errors.New("消息发送失败"))
		return
	}

	// 发送消息
	UserSendMessage(getRoom(m.RoomId), m.Msg, user)
}

/**
 * 获取字典
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
func GetDictionary() (dictionary [][]string) {
	// TODO 配置文件
	data, err := ioutil.ReadFile("./keyword/keyword.text")
	if err != nil {
		panic(err)
	}

	// 读取
	for _, keywords := range strings.Split(string(data), "\n") {
		keywordArr := strings.Split(keywords, "——")
		dictionary = append(dictionary, keywordArr)
	}

	return dictionary
}

/**
 * 获取参数
 * @Author: cs_shuai
 * @Date: 2020-09-19
 */
func getArgs(args []interface{}) (i interface{}, a gate.Agent, user *GameMsg.User) {
	// 收到的 Hello 消息
	i = args[0]
	fmt.Println("---------------" + fmt.Sprint(args[0]) + "---------------")
	// 消息的发送者
	a = args[1].(gate.Agent)
	if a.UserData() != nil {
		user = a.UserData().(*GameMsg.User)
	} else {
		user = new(GameMsg.User)
	}

	return
}
