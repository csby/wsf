package types

import (
	"time"
)

const (
	TokenName = "token"
)

const (
	TokenTypeNone            = 0 // 不需要凭证
	TokenTypeAccountPassword = 1 // 账号及密码
)

const (
	TokenPlaceHeader = 0 // 凭证在头部
	TokenPlaceQuery  = 1 // 凭证在参数
)

const (
	TokenValueKindEdit      = 0 // 编辑筐
	TokenValueKindPassword  = 1 // 密码筐
	TokenValueKindSelection = 2 // 选择筐
)

type TokenDatabase interface {
	Name() string
	Set(key string, data interface{})
	Get(key string, delay bool) (interface{}, bool)
	Del(key string) bool
	Lst(key string) []interface{}
	Permanent(key string, val bool) bool
}

type TokenAuth struct {
	Name  string `json:"name" note:"名称"`
	Value string `json:"value" note:"值"`
}

type TokenUI struct {
	TokenAuth

	Required   bool                `json:"required" note:"是否必填"`
	Label      string              `json:"label" note:"标签"`
	ValueKind  int                 `json:"valueKind" note:"0-编辑筐(edit); 1-密码筐(password); 2-选择筐(selection)"`
	Selections []TokenUISelectItem `json:"selections" note:"ValueType == 2时的可选项"`
}

type TokenUISelectItem struct {
	Name  string `json:"name" note:"显示名称"`
	Value string `json:"value" note:"值"`
}

var (
	tokenUIForAccountPassword = []TokenUI{
		{
			TokenAuth: TokenAuth{
				Name:  "account",
				Value: "",
			},
			Required:  true,
			Label:     "账号",
			ValueKind: TokenValueKindEdit,
		},
		{
			TokenAuth: TokenAuth{
				Name:  "password",
				Value: "",
			},
			Label:     "密码",
			ValueKind: TokenValueKindPassword,
		},
	}
)

func TokenUIForAccountPassword() []TokenUI {
	return tokenUIForAccountPassword
}

type Token struct {
	ID          string    `json:"id" note:"标识ID"`
	UserAccount string    `json:"userAccount" note:"用户账号"`
	UserName    string    `json:"userName" note:"用户姓名"`
	LoginIP     string    `json:"loginIp" note:"用户登陆IP"`
	LoginTime   time.Time `json:"loginTime" note:"登陆时间"`
	ActiveTime  time.Time `json:"activeTime" note:"最近激活时间"`
	Usage       int       `json:"usage" "使用次数"`

	Ext interface{} `json:"ext" note:"扩展信息"`
}

type TokenFilter struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	FunId    string `json:"funId"`
}

type OnlineUser struct {
	Token         string   `json:"token" note:"凭证"`
	UserAccount   string   `json:"userAccount" note:"用户账号"`
	UserName      string   `json:"userName" note:"用户姓名"`
	LoginIP       string   `json:"loginIp" note:"用户登陆IP"`
	LoginTime     DateTime `json:"loginTime" note:"登陆时间"`
	LoginDuration string   `json:"loginDuration" note:"登陆时时长"`
}

func (s *OnlineUser) CopyFrom(token *Token) {
	if token == nil {
		return
	}

	s.Token = token.ID
	s.UserAccount = token.UserAccount
	s.UserName = token.UserName
	s.LoginIP = token.LoginIP
	s.LoginTime = DateTime(token.LoginTime)
}
