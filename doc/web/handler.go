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
	ApiPathInformation    = "/info"
	ApiPathCatalogTree    = "/catalog/tree"
	ApiPathFunctionDetail = "/function/:id"
	ApiPathTokenUI        = "/token/ui/:id"
	ApiPathTokenCreate    = "/token/create/:id"
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
	Init(router types.Router, info *types.ServerInformation)
}

type handler struct {
	rootPath   string
	sitePrefix string
	apiPrefix  string
}

func (s *handler) Init(router types.Router, info *types.ServerInformation) {
	// site
	sitePath := types.Path{Prefix: s.sitePrefix}
	router.ServeFiles(sitePath.New("/*filepath"), nil, http.Dir(s.rootPath), nil)

	// api
	apiPath := types.Path{Prefix: s.apiPrefix}
	ctrl := &controller{doc: router.Document()}
	if info != nil {
		ctrl.info.Name = info.Name
		ctrl.info.Version = info.Version
	}

	// 获取服务信息
	router.POST(apiPath.New(ApiPathInformation), nil, ctrl.GetInformation, nil)

	// 获取接口目录信息
	router.POST(apiPath.New(ApiPathCatalogTree), nil, ctrl.GetCatalogTree, nil)

	// 获取接口定义信息
	router.POST(apiPath.New(ApiPathFunctionDetail), nil, ctrl.GetFunctionDetail, nil)

	// 获取创建凭证的输入项目
	router.POST(apiPath.New(ApiPathTokenUI), nil, ctrl.GetTokenUI, nil)

	// 创建凭证
	router.POST(apiPath.New(ApiPathTokenCreate), nil, ctrl.CreateToken, nil)
}
