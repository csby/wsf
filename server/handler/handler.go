package handler

import (
	"fmt"
	"github.com/csby/wsf/doc"
	"github.com/csby/wsf/doc/web"
	"github.com/csby/wsf/router"
	"github.com/csby/wsf/types"
	"net/http"
)

func NewHttpHandler(log types.Log, handler types.HttpHandler) (http.Handler, error) {
	instance := &httpHandler{handler: handler, router: router.New()}
	instance.SetLog(log)
	instance.rid = &randNumber{id: 0, max: 0}

	if handler != nil {
		var server1erInfo *types.ServerInformation = nil
		redirectToHttps := false
		documentEnabled := false
		documentRoot := ""
		extend := handler.Extend()
		if extend != nil {
			redirectToHttps = extend.RedirectToHttps()
			documentEnabled = extend.DocumentEnabled()
			documentRoot = extend.DocumentRoot()
			server1erInfo = extend.ServerInfo()
		}

		instance.router.NotFound = handler.NotFound()
		instance.redirectToHttps = redirectToHttps
		instance.router.Doc = doc.NewDoc(documentEnabled)

		instance.router.Doc.OnFunctionReady(func(index int, method, path, name string) {
			if log != nil {
				log.Debug(fmt.Sprintf("api-%03d", index), ": [", method, "] ", path, " (", name, ") has been ready")
			}
		})

		handler.Map(instance.router)

		if documentEnabled {
			docHandler := web.NewHandler(documentRoot, web.SitePath, web.ApiPath)
			docHandler.Init(instance.router, server1erInfo)

			if log != nil {
				log.Info("document for api is enabled")
				log.Info("document information api path: [POST] ", web.ApiPath, web.ApiPathInformation)
				log.Info("document catalog api path: [POST] ", web.ApiPath, web.ApiPathCatalogTree)
				log.Info("document function api path: [POST] ", web.ApiPath, web.ApiPathFunctionDetail)
				log.Info("document token ui api path: [POST] ", web.ApiPath, web.ApiPathTokenUI)
				log.Info("document token create api path: [POST] ", web.ApiPath, web.ApiPathTokenCreate)
			}
		}
	}

	return instance, nil
}
