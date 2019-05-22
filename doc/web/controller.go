package web

import (
	"github.com/csby/wsf/types"
	"net/http"
)

type controller struct {
	doc  types.Doc
	info types.ServerInformation
}

func (s *controller) GetInformation(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success(&s.info)
}

func (s *controller) GetCatalogTree(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	if s.doc == nil {
		a.Error(types.ErrInternal, "doc is nil")
		return
	}

	a.Success(s.doc.Catalogs())
}

func (s *controller) GetFunctionDetail(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	if s.doc == nil {
		a.Error(types.ErrInternal, "doc is nil")
		return
	}

	id := p.ByName("id")
	fun, err := s.doc.Function(id, a.Schema(), r.Host)
	if err != nil {
		a.Error(types.ErrInput, err)
		return
	}

	a.Success(fun)
}

func (s *controller) GetTokenUI(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	if s.doc == nil {
		a.Error(types.ErrInternal, "doc is nil")
		return
	}

	id := p.ByName("id")
	ui, err := s.doc.TokenUI(id)
	if err != nil {
		a.Error(types.ErrInput, err)
		return
	}

	a.Success(ui)
}

func (s *controller) CreateToken(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	if s.doc == nil {
		a.Error(types.ErrInternal, "doc is nil")
		return
	}

	items := make([]types.TokenAuth, 0)
	err := a.GetJson(&items)
	if err != nil {
		a.Error(types.ErrInput, err)
		return
	}

	id := p.ByName("id")
	token, code, err := s.doc.TokenCreate(id, items, a)
	if code != nil {
		a.Error(code, err)
		return
	}

	a.Success(token)
}
