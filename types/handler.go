package types

import "net/http"

type HttpHandler interface {
	Map(router Router)
	PreRouting(w http.ResponseWriter, r *http.Request, a Assistant) bool
	PostRouting(w http.ResponseWriter, r *http.Request, a Assistant)
	NotFound() func(http.ResponseWriter, *http.Request, Assistant)

	Extend() HttpHandlerExtend
}

type HttpHandlerExtend interface {
	Restart() func() error
	RedirectToHttps() bool
	DocumentEnabled() bool
	DocumentRoot() string
	ServerInfo() *ServerInformation
}

type TcpHandler interface {
}
