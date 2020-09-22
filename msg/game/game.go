package GameMsg

import "github.com/name5566/leaf/gate"

// 用户
type User struct {
	Id     int         `gorm:"AUTO_INCREMENT"` // Id
	Openid string      `gorm:"size:255"`       // 标识
	No     string      `gorm:"size:255"`       // 序号
	Name   string      `gorm:"size:255"`       // 用户名
	Word   string      `gorm:"-"`              // 词语名称
	Role   string      `gorm:"-"`              // 角色
	Status int         `gorm:"size:1"`         // 状态
	RoomId string      `gorm:"size:255"`       // 房间号
	Agent  *gate.Agent `gorm:"-"`
}

// 词语
type Word struct {
	Id          int    // Id
	Word        string // 词
	AnotherWord string // 另一个词
}

// 消息
type GameMessage struct {
	UserId   string // 用户Id
	UserName string // 用户名称
	RoomId   string // 房间id
	Msg      string // 消息
	Status   int    // 状态
	Data     interface{}
	Type     string
}

type Login struct {
	UserName string // 用户名称
	UserId   string // 用户名称
}

type Logout struct {
	UserName string // 用户名称
	UserId   string // 用户名称
}

type Room struct {
	Id               int                 `gorm:"AUTO_INCREMENT"`
	CreateUser       *User               `gorm:"-"`
	CreateUserId     string              `gorm:"size:255"` // 用户ID
	Msg              string              `gorm:"-"`
	RoomId           string              `gorm:"size:255"` // 房间ID
	Password         string              `gorm:"size:255"` // 房间密码
	TotalNumber      string              `gorm:"size:255"` // 总人数
	Number           int                 `gorm:"size:11"`  // 当前人数
	UndercoverNumber string              `gorm:"size:11"`  // 卧底人数
	UserList         map[string]*User    `gorm:"-"`
	GameInfo         *Game               `gorm:"-"`
	MsgChan          chan string         `json:"-" gorm:"-"`
	PrepareList      map[string]string   `gorm:"-"`
	UsedWord         map[string]*Keyword `gorm:"-"` // 使用过词语
	PrepareNum       int                 `gorm:"-"`
	IsPrepare        bool                `gorm:"-"`
}

type RoomOut struct {
	RoomId string
	UserId string
}

// 游戏
type Game struct {
	Round            int              // 回合数
	SurvivalUserList map[string]*User `gorm:"-"` // 存活用户列表
	Keyword          *Keyword         // 词语
	UndercoverNum    int              // 卧底数量
	Stage            string           // 阶段
	ActionTime       int              // 操作时间
	VoteTime         int              // 投票等待时间 (秒)
	VoteList         map[string]*Vote // 投票列表
	RoomId           string           // 房间号
	VoteChan         chan *Vote       `json:"-" gorm:"-"` // 投票通道
	VoteNum          int              `gorm:"-"`          // 投票次数
	WinRole          string           // 胜利方
	OutUser          []*User          `gorm:"-"`
}

type Vote struct {
	Round            int
	UserId           string
	VotePlayerNumber string
	IsPrepare        bool
	RoomId           string
	GameId           string
}

/**
 * 词组
 * @Author: cs_shuai
 * @Date: 2020-09-11
 */
type Keyword struct {
	Code           string
	NormalWord     string
	UndercoverWord string
	Vension        int64
}

const (
	Role_Normal     = "Normal"     // 正常
	Role_Undercover = "Undercover" // 卧底

	GameStage_Start = "Start" // 准备阶段
	GameStage_Vote  = "Vote"  // 投票阶段
	GameStage_Game  = "Game"  // 游戏阶段
	GameStage_Over  = "Over"  // 完成阶段

)
