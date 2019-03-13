package web

import (
	"fmt"
	"github.com/csby/wsf/types"
	"net/http"
)

type controller struct {
	doc types.Doc
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
	addr := fmt.Sprintf("%s://%s", a.Schema(), r.Host)
	fun, err := s.doc.Function(id, addr)
	if err != nil {
		a.Error(types.ErrInput, err)
		return
	}

	a.Success(fun)
}
