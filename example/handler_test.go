package example

import (
	"github.com/csby/wsf/doc/web"
	"github.com/csby/wsf/types"
	"net/http"
)

const (
	uriHello = "/hello"
)

var (
	testPath = types.Path{Prefix: "/test", TokenKind: 3, DefaultShortenUrl: true, DefaultTokenType: 1}
)

type HttpHandler struct {
	restart func() error

	controller *Controller
}

func (s *HttpHandler) Map(router types.Router) {
	doc := router.Document()
	if doc != nil {
		doc.OnFunctionReady(func(index int, method, path, name string) {
			log.Debug("api-", index, ": [", method, "] ", path, " (", name, ") has been ready")
		})
	}

	s.controller = &Controller{}

	router.POST(testPath.New(uriHello), nil, s.controller.Hello, s.controller.HelloDoc)

	if cfg.Server.Document.Enabled {
		docHandler := web.NewHandler(cfg.Server.Document.Root, web.SitePath, web.ApiPath)
		docHandler.Init(router)
	}
}

func (s *HttpHandler) PreRouting(w http.ResponseWriter, r *http.Request, a types.Assistant) bool {
	// enable across access
	if r.Method == "OPTIONS" {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type,token")
		return true
	}

	return false
}

func (s *HttpHandler) PostRouting(w http.ResponseWriter, r *http.Request, a types.Assistant) {

}

func (s *HttpHandler) NotFound() func(http.ResponseWriter, *http.Request, types.Assistant) {
	return nil
}

func (s *HttpHandler) Restart() func() error {
	return s.restart
}

func (s *HttpHandler) RedirectToHttps() bool {
	return false
}

func (s *HttpHandler) EnableDocument() bool {
	return cfg.Server.Document.Enabled
}
