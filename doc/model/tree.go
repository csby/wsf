package model

type Tree struct {
	ID       string `json:"id"`       // 标识
	Name     string `json:"name"`     // 名称
	Note     string `json:"note"`     // 说明
	Type     int    `json:"type"`     // 类别: 0-catalog; 1-function
	Keywords string `json:"keywords"` // 关键字, 用于过滤

	Children []*Tree `json:"children"`
}
