package types

import "encoding/json"

type Result struct {
	Code   int         `json:"code" note:"结果状态码, 0-表示成功; 其它-表示失败(如: 101-凭证失效须重新登陆)"`
	Serial uint64      `json:"serial" note:"请求序号, 一般用于错误定位或异步调用结果查询标识"`
	Elapse string      `json:"elapse" note:"耗时, 接口在服务端执行时间"`
	Error  ResultError `json:"error" note:"失败时的错误信息"`
	Data   interface{} `json:"data" note:"成功时的结果数据"`
}

type ResultError struct {
	Summary string `json:"summary" note:"描述信息, 用于错误信息提示"`
	Detail  string `json:"detail" note:"详细信息, 用于错误排查"`
}

func (s *Result) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Result) Unmarshal(v []byte) error {
	return json.Unmarshal(v, s)
}

func (s *Result) GetData(v interface{}) error {
	data, err := json.Marshal(s.Data)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func (s *Result) FormatString() string {
	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return ""
	}

	return string(bytes[:])
}
