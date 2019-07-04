package model

type Form struct {
	Key       string      `json:"key" note:"标志"`
	Value     interface{} `json:"value" note:"值"`
	ValueKind int         `json:"valueKind" note:"类型: 0-text(文本); 1-file(文件)"`
	Required  bool        `json:"required" note:"是否必填"`
	Note      string      `json:"note" note:"说明"`
}
