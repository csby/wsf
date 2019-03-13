package handler

import (
	"github.com/csby/wsf/doc"
	"github.com/csby/wsf/router"
	"github.com/csby/wsf/types"
	"net/http"
)

func NewHttpHandler(log types.Log, handler types.HttpHandler) (http.Handler, error) {
	instance := &httpHandler{handler: handler, router: router.New()}
	instance.SetLog(log)
	instance.rid = &randNumber{id: 0, max: 0}

	if handler != nil {
		instance.redirectToHttps = handler.RedirectToHttps()
		instance.router.NotFound = handler.NotFound()
		instance.router.Doc = doc.NewDoc(handler.EnableDocument())

		handler.Map(instance.router)
	}

	return instance, nil
}
