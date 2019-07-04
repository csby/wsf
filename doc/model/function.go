package model

import (
	"github.com/csby/wsf/types"
	"strings"
)

const (
	headContentType = "content-type"
)

var (
	modelArgument = &Argument{}
)

type Function struct {
	ID            string      `json:"id"`            // 接口标识
	Name          string      `json:"name"`          // 接口名称
	Note          string      `json:"note"`          // 接口说明
	Method        string      `json:"method"`        // 接口方法
	Path          string      `json:"path"`          // 接口地址
	FullPath      string      `json:"fullPath"`      // 接口地址
	TokenType     int         `json:"tokenType"`     // 凭证类型
	TokenPlace    int         `json:"-"`             // 凭证位置
	WebSocket     bool        `json:"webSocket"`     // 是否为websocket接口
	InputHeaders  []*Header   `json:"inputHeaders"`  // 输入头部
	InputQueries  []*Query    `json:"inputQueries"`  // 输入参数
	InputForms    []*Form     `json:"inputForms"`    // 输入表单
	InputModel    *Argument   `json:"inputModel"`    // 输入数据
	InputSample   interface{} `json:"inputSample"`   // 输入数据示例
	OutputHeaders []*Header   `json:"outputHeaders"` // 输出头部
	OutputModel   *Argument   `json:"outputModel"`   // 输出数据
	OutputSample  interface{} `json:"outputSample"`  // 输出数据示例

	TokenUI     func() []types.TokenUI                                                            `json:"-"`
	TokenCreate func(items []types.TokenAuth, a types.Assistant) (string, types.ErrorCode, error) `json:"-"`
}

func (s *Function) SetNote(v string) {
	s.Note = v
}

func (s *Function) SetTokenType(v int) {
	if s.TokenType == v {
		return
	}
	s.TokenType = v
	s.tokenTypeChanged()
}

func (s *Function) SetInputContentType(v string) {
	items := make([]*Header, 0)
	c := len(s.InputHeaders)
	for i := 0; i < c; i++ {
		item := s.InputHeaders[i]
		if item == nil {
			continue
		}
		if item.Name == headContentType {
			continue
		}

		items = append(items, item)
	}

	if len(v) > 0 {
		h := &Header{
			Name:         headContentType,
			Note:         "内容类型",
			Required:     true,
			Values:       []string{v},
			DefaultValue: v,
		}

		items = append(items, h)
	}

	s.InputHeaders = items
}

func (s *Function) AddInputHeader(required bool, name, note, defaultValue string, optionValues ...string) {
	header := s.GetInputHeader(name)
	if header != nil {
		header.Required = required
		header.Note = note
		header.DefaultValue = defaultValue
		header.Values = optionValues
	} else {
		s.InputHeaders = append(s.InputHeaders, &Header{
			Name:         name,
			Note:         note,
			Required:     required,
			Values:       optionValues,
			DefaultValue: defaultValue,
		})
	}
}

func (s *Function) ClearInputHeader() {
	s.InputHeaders = make([]*Header, 0)
}

func (s *Function) RemoveInputHeader(name string) {
	items := make([]*Header, 0)
	c := len(s.InputHeaders)
	for i := 0; i < c; i++ {
		item := s.InputHeaders[i]
		if item == nil {
			continue
		}
		if item.Name == types.TokenName {
			continue
		}
		if item.Name == name {
			continue
		}

		items = append(items, item)
	}

	s.InputHeaders = items
}

func (s *Function) AddInputQuery(required bool, name, note, defaultValue string, optionValues ...string) {
	query := s.GetInputQuery(name)
	if query != nil {
		query.Required = required
		query.Note = note
		query.DefaultValue = defaultValue
		query.Values = optionValues
	} else {
		s.InputQueries = append(s.InputQueries, &Query{
			Name:         name,
			Note:         note,
			Required:     required,
			Values:       optionValues,
			DefaultValue: defaultValue,
		})
	}
}

func (s *Function) RemoveInputQuery(name string) {
	items := make([]*Query, 0)
	c := len(s.InputQueries)
	for i := 0; i < c; i++ {
		item := s.InputQueries[i]
		if item == nil {
			continue
		}
		if item.Name == types.TokenName {
			continue
		}
		if item.Name == name {
			continue
		}

		items = append(items, item)
	}

	s.InputQueries = items
}

func (s *Function) AddInputForm(required bool, key, note string, valueKind int, defaultValue interface{}) {
	form := s.GetInputForm(key)
	if form != nil {
		form.Required = required
		form.Note = note
		form.Value = defaultValue
		form.ValueKind = valueKind
	} else {
		s.InputForms = append(s.InputForms, &Form{
			Key:       key,
			Note:      note,
			Required:  required,
			Value:     defaultValue,
			ValueKind: valueKind,
		})
	}
}

func (s *Function) RemoveInputForm(key string) {
	items := make([]*Query, 0)
	c := len(s.InputQueries)
	for i := 0; i < c; i++ {
		item := s.InputQueries[i]
		if item == nil {
			continue
		}
		if item.Name == key {
			continue
		}

		items = append(items, item)
	}

	s.InputQueries = items
}

func (s *Function) SetInputExample(v interface{}) {
	s.InputSample = v
	s.InputModel = modelArgument.FromExample(v)
}

func (s *Function) AddOutputHeader(name, value string) {
	header := s.GetOutputHeader(name)
	if header != nil {
		header.DefaultValue = value
	} else {
		s.OutputHeaders = append(s.OutputHeaders, &Header{
			Name:         name,
			DefaultValue: value,
		})
	}
}

func (s *Function) ClearOutputHeader() {
	s.OutputHeaders = make([]*Header, 0)
}

func (s *Function) SetOutputExample(v interface{}) {
	s.OutputSample = v
	s.OutputModel = modelArgument.FromExample(v)
}

func (s *Function) SetOutputDataExample(v interface{}) {
	s.OutputSample = &types.Result{
		Code:   0,
		Serial: 201805161315480008,
		Error: types.ResultError{
			Summary: "",
			Detail:  "",
		},
		Data: v,
	}
	s.OutputModel = modelArgument.FromExample(s.OutputSample)
}

func (s *Function) GetInputHeader(name string) *Header {
	c := len(s.InputHeaders)
	for i := 0; i < c; i++ {
		item := s.InputHeaders[i]
		if strings.ToLower(item.Name) == strings.ToLower(name) {
			return item
		}
	}

	return nil
}

func (s *Function) GetInputQuery(name string) *Query {
	c := len(s.InputQueries)
	for i := 0; i < c; i++ {
		item := s.InputQueries[i]
		if strings.ToLower(item.Name) == strings.ToLower(name) {
			return item
		}
	}

	return nil
}

func (s *Function) GetInputForm(key string) *Form {
	c := len(s.InputForms)
	for i := 0; i < c; i++ {
		item := s.InputForms[i]
		if strings.ToLower(item.Key) == strings.ToLower(key) {
			return item
		}
	}

	return nil
}

func (s *Function) GetOutputHeader(name string) *Header {
	c := len(s.OutputHeaders)
	for i := 0; i < c; i++ {
		item := s.OutputHeaders[i]
		if strings.ToLower(item.Name) == strings.ToLower(name) {
			return item
		}
	}

	return nil
}

func (s *Function) tokenTypeChanged() {
	if s.TokenType != types.TokenTypeNone {
		if s.TokenPlace == types.TokenPlaceQuery {
			s.RemoveInputHeader(types.TokenName)
			inputQuery := s.GetInputQuery(types.TokenName)
			if inputQuery != nil {
				inputQuery.Token = true
			} else {
				s.InputQueries = append(s.InputQueries, &Query{
					Name:     types.TokenName,
					Note:     "凭证",
					Required: true,
					Token:    true,
				})
			}
		} else {
			s.RemoveInputQuery(types.TokenName)
			inputHeader := s.GetInputHeader(types.TokenName)
			if inputHeader != nil {
				inputHeader.Token = true
			} else {
				s.InputHeaders = append(s.InputHeaders, &Header{
					Name:     types.TokenName,
					Note:     "凭证",
					Required: true,
					Token:    true,
				})
			}
		}
	} else {
		s.RemoveInputHeader(types.TokenName)
		s.RemoveInputQuery(types.TokenName)
	}
}
