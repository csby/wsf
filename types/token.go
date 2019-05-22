package types

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
