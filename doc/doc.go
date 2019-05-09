package doc

import (
	"encoding/hex"
	"fmt"
	"github.com/csby/wsf/doc/model"
	"github.com/csby/wsf/types"
	"hash/adler32"
)

func NewDoc(enable bool) types.Doc {
	return &doc{
		enable:    enable,
		catalogs:  make(model.CatalogSlice, 0),
		functions: make(map[string]*model.Function),
	}
}

type doc struct {
	enable    bool
	catalogs  model.CatalogSlice
	functions map[string]*model.Function

	onFunctionReady func(index int, method, path, name string)
}

func (s *doc) Enable() bool {
	return s.enable
}

func (s *doc) AddCatalog(name string) types.Catalog {
	c := len(s.catalogs)
	for i := 0; i < c; i++ {
		item := s.catalogs[i]
		if item.Name == name {
			return item
		}
	}

	item := &model.Catalog{Name: name}
	item.OnAddFunction(s.onNewFunction)
	item.Children = make(model.CatalogSlice, 0)

	s.catalogs = append(s.catalogs, item)

	return item
}

func (s *doc) OnFunctionReady(f func(index int, method, path, name string)) {
	s.onFunctionReady = f
}

func (s *doc) Catalogs() interface{} {
	return s.catalogs
}

func (s *doc) Function(id, addr string) (interface{}, error) {
	fun, ok := s.functions[id]
	if ok {
		fun.FullPath = fmt.Sprintf("%s%s", addr, fun.Path)
		return fun, nil
	} else {
		return nil, fmt.Errorf("id not '%s' exist", id)
	}
}

func (s *doc) onNewFunction(fun *model.Function) {
	id := s.generateFunctionId(fun.Method, fun.Path)
	_, ok := s.functions[id]
	if ok {
		panic(fmt.Sprintf("a document handle is already registered for path '%s: %s'", fun.Method, fun.Path))
	}

	fun.ID = id
	s.functions[id] = fun

	if s.onFunctionReady != nil {
		s.onFunctionReady(len(s.functions), fun.Method, fun.Path, fun.Name)
	}
}

func (s *doc) generateFunctionId(method, path string) string {
	h := adler32.New()
	_, err := h.Write([]byte(method + path))
	if err != nil {
		return path
	}

	return hex.EncodeToString(h.Sum(nil))
}
