package example

import (
	"github.com/csby/wsf/types"
	"net/http"
)

type Controller struct {
}

func (s *Controller) Hello(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success("Hello")
}

func (s *Controller) HelloDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := doc.AddCatalog("UnitTest").AddChild("Restful")
	function := catalog.AddFunction(method, path, "Hello World")
	function.SetNote("restful api, return data with 'Hello'")
	function.SetOutputDataExample("Hello")
}
