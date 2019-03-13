package types

import "net/http"

type HttpHandler interface {
	Map(router Router)
	PreRouting(w http.ResponseWriter, r *http.Request, a Assistant) bool
	PostRouting(w http.ResponseWriter, r *http.Request, a Assistant)

	NotFound() func(http.ResponseWriter, *http.Request, Assistant)
	Restart() func() error
	RedirectToHttps() bool
	EnableDocument() bool
}

type TcpHandler interface {
}
