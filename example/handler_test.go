package example

import (
	"github.com/csby/wsf/types"
	"net/http"
)

const (
	uriHello = "/hello"
)

var (
	testPath = types.Path{Prefix: "/test", DefaultShortenUrl: true, DefaultTokenType: 1}
)

type HttpHandler struct {
	controller *Controller
}

func (s *HttpHandler) Map(router types.Router) {
	s.controller = &Controller{}

	router.POST(testPath.New(uriHello), nil, s.controller.Hello, s.controller.HelloDoc)

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
	log.Info("PostRouting: ", a.Path())
}

func (s *HttpHandler) NotFound() func(http.ResponseWriter, *http.Request, types.Assistant) {
	return nil
}

var httpHandlerExtend = &HttpHandlerExtend{}

func (s *HttpHandler) Extend() types.HttpHandlerExtend {
	return httpHandlerExtend
}

type HttpHandlerExtend struct {
	restart func() error
}

func (s *HttpHandlerExtend) Restart() func() error {
	return s.restart
}

func (s *HttpHandlerExtend) RedirectToHttps() bool {
	return false
}

func (s *HttpHandlerExtend) DocumentEnabled() bool {
	return cfg.Server.Document.Enabled
}

func (s *HttpHandlerExtend) DocumentRoot() string {
	return cfg.Server.Document.Root
}

func (s *HttpHandlerExtend) ServerInfo() *types.ServerInformation {
	return &types.ServerInformation{Name: "unit-test", Version: "1.0.1.0"}
}
