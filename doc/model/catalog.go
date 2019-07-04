package model

import (
	"fmt"
	"github.com/csby/wsf/types"
	"sort"
)

const (
	typeCatalog  = 0
	typeFunction = 1
)

type Catalog struct {
	ID       string       `json:"id"`       // 标识
	Name     string       `json:"name"`     // 名称
	Note     string       `json:"note"`     // 说明
	Type     int          `json:"type"`     // 类别: 0-catalog; 1-function
	Keywords string       `json:"keywords"` // 关键字, 用于过滤
	Children CatalogSlice `json:"children"`

	index         int
	onAddFunction func(fun *Function)
}

func (s *Catalog) AddChild(name string) types.Catalog {
	c := len(s.Children)
	for i := 0; i < c; i++ {
		item := s.Children[i]
		if item.Name == name {
			return item
		}
	}

	item := &Catalog{Name: name}
	item.Children = make(CatalogSlice, 0)
	item.Type = typeCatalog
	item.Keywords = name
	item.index = len(s.Children)
	item.onAddFunction = s.onAddFunction

	s.Children = append(s.Children, item)
	sort.Sort(s.Children)

	return item
}

func (s *Catalog) AddFunction(method string, httpPath types.HttpPath, name string) types.Function {
	path := httpPath.Path()
	item := &Catalog{Name: name}
	item.Children = make(CatalogSlice, 0)
	item.Type = typeFunction
	item.Keywords = fmt.Sprintf("%s%s", name, path)
	item.index = len(s.Children)

	s.Children = append(s.Children, item)
	sort.Sort(s.Children)

	fuc := &Function{
		Method:      method,
		Path:        path,
		Name:        name,
		TokenPlace:  httpPath.TokenPlace(),
		WebSocket:   httpPath.IsWebSocket(),
		TokenUI:     httpPath.TokenUI(),
		TokenCreate: httpPath.TokenCreate(),
	}
	if fuc.WebSocket {
		fuc.Method = "WEBSOCKET"
	}
	fuc.InputHeaders = make([]*Header, 0)
	fuc.InputQueries = make([]*Query, 0)
	fuc.InputForms = make([]*Form, 0)
	fuc.OutputHeaders = make([]*Header, 0)
	fuc.SetTokenType(httpPath.TokenType())
	if method == "POST" {
		fuc.SetInputContentType(types.ContentTypeJson)
		fuc.AddOutputHeader("access-control-allow-origin", "*")
		fuc.AddOutputHeader(headContentType, "application/json;charset=utf-8")
	}
	if s.onAddFunction != nil {
		s.onAddFunction(fuc)
		item.ID = fuc.ID
	}

	return fuc
}

func (s *Catalog) OnAddFunction(f func(fun *Function)) {
	s.onAddFunction = f
}

type CatalogSlice []*Catalog

func (s CatalogSlice) Len() int {
	return len(s)
}

func (s CatalogSlice) Less(i, j int) bool {
	if s[i].Type < s[j].Type {
		return true
	} else if s[i].Type == s[j].Type {
		return s[i].index < s[j].index
	}

	return false
}

func (s CatalogSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
