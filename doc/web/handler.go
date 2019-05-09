package web

import (
	"github.com/csby/wsf/types"
	"net/http"
)

const (
	SitePath = "/doc"
	ApiPath  = "/doc.api"
)

const (
	ApiPathCatalogTree    = "/catalog/tree"
	ApiPathFunctionDetail = "/function/:id"
)

// rootPath: site path in location
// sitePath: document site prefix path, http(s)://ip/[SitePath]/*
// apiPath: document api prefix path, http(s)://ip/[ApiPath]/*
func NewHandler(rootPath, sitePrefix, apiPrefix string) Handler {
	return &handler{
		rootPath:   rootPath,
		sitePrefix: sitePrefix,
		apiPrefix:  apiPrefix,
	}
}

type Handler interface {
	Init(router types.Router)
}

type handler struct {
	rootPath   string
	sitePrefix string
	apiPrefix  string
}

func (s *handler) Init(router types.Router) {
	// site
	sitePath := types.Path{Prefix: s.sitePrefix}
	router.ServeFiles(sitePath.New("/*filepath"), nil, http.Dir(s.rootPath), nil)

	// api
	apiPath := types.Path{Prefix: s.apiPrefix}
	ctrl := &controller{doc: router.Document()}

	// 获取接口目录信息
	router.POST(apiPath.New(ApiPathCatalogTree), nil, ctrl.GetCatalogTree, nil)

	// 获取接口定义信息
	router.POST(apiPath.New(ApiPathFunctionDetail), nil, ctrl.GetFunctionDetail, nil)
}
