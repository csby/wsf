package model

type Model struct {
	Name     string  `json:"name"`
	Children []*Item `json:"children"`
}

type Item struct {
	Name     string `json:"name"`     // 名称
	Type     string `json:"type"`     // 类型
	Note     string `json:"note"`     // 说明
	Required bool   `json:"required"` // 必填
}
