package opt

import (
	"fmt"
	"github.com/csby/wsf/opt/controller"
	"github.com/csby/wsf/server/configure"
	"github.com/csby/wsf/types"
	"net/http"
)

type Handler interface {
	Init(router types.Router, api func(path types.Path, router types.Router, tokenChecker types.RouterPreHandle) error) error
	NotFound(w http.ResponseWriter, r *http.Request, a types.Assistant)
}

func NewHandler(log types.Log, cfg *configure.Configure, db types.TokenDatabase, chs types.SocketChannelCollection) Handler {
	instance := &handler{}
	instance.SetLog(log)
	instance.cfg = cfg
	instance.dbToken = db
	instance.wsChannels = chs
	instance.svcMgr = &SvcUpdMgr{}

	return instance
}

type handler struct {
	types.Base

	cfg        *configure.Configure
	dbToken    types.TokenDatabase
	wsChannels types.SocketChannelCollection
	svcMgr     types.SvcUpdMgr

	auth      *controller.Auth
	monitor   *controller.Monitor
	service   *controller.Service
	update    *controller.Update
	site      *controller.Site
	websocket *controller.Websocket
}

func (s *handler) Init(router types.Router, api func(path types.Path, router types.Router, tokenChecker types.RouterPreHandle) error) error {
	if s.cfg == nil {
		return fmt.Errorf("opt: invalid config (cfg = nil)")
	}
	if s.dbToken == nil {
		return fmt.Errorf("opt: invalid token database (dbToken = nil)")
	}
	if s.wsChannels == nil {
		return fmt.Errorf("opt: invalid websocket channals (chs = nil)")
	}

	optApiPath.DefaultTokenCreate = s.createTokenForAccountPassword

	s.mapOptApi(optApiPath, router)
	s.mapOptSite(optWebPath, router, s.cfg.Operation.Root)
	s.mapWebappSite(webappWebPath, router, s.cfg.Webapp.Root)

	if api != nil {
		return api(optApiPath, router, s.auth.CheckToken)
	}

	return nil
}

func (s *handler) NotFound(w http.ResponseWriter, r *http.Request, a types.Assistant) {
	if r.Method == "GET" {
		http.FileServer(http.Dir(s.cfg.Root)).ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (s *handler) mapOptApi(path types.Path, router types.Router) {
	s.auth = controller.NewAuth(s.GetLog(), s.cfg, s.dbToken, s.wsChannels)
	s.monitor = controller.NewMonitor(s.GetLog(), s.cfg)
	s.service = controller.NewService(s.GetLog(), s.cfg, s.svcMgr)
	s.update = controller.NewUpdate(s.GetLog(), s.cfg, s.svcMgr)
	s.site = controller.NewSite(s.GetLog(), s.cfg, s.dbToken, s.wsChannels, optWebPath.Prefix, webappWebPath.Prefix)
	s.websocket = controller.NewWebsocket(s.GetLog(), s.cfg, s.dbToken, s.wsChannels)

	tokenChecker := s.auth.CheckToken
	// 获取验证码
	router.POST(path.New("/captcha").SetTokenType(types.TokenTypeNone),
		nil, s.auth.GetCaptcha, s.auth.GetCaptchaDoc)
	// 用户登陆
	router.POST(path.New("/login").SetTokenType(types.TokenTypeNone),
		nil, s.auth.Login, s.auth.LoginDoc)
	// 注销登陆
	router.POST(path.New("/logout"),
		tokenChecker, s.auth.Logout, s.auth.LogoutDoc)
	// 获取登录账号
	router.POST(path.New("/login/account"),
		tokenChecker, s.auth.GetLoginAccount, s.auth.GetLoginAccountDoc)
	// 获取在线用户
	router.POST(path.New("/online/users"),
		tokenChecker, s.auth.GetOnlineUsers, s.auth.GetOnlineUsersDoc)

	// 系统信息
	router.POST(path.New("/monitor/host"),
		tokenChecker, s.monitor.GetHost, s.monitor.GetHostDoc)
	router.POST(path.New("/monitor/network/interfaces"),
		tokenChecker, s.monitor.GetNetworkInterfaces, s.monitor.GetNetworkInterfacesDoc)
	router.POST(path.New("/monitor/network/listen/ports"),
		tokenChecker, s.monitor.GetNetworkListenPorts, s.monitor.GetNetworkListenPortsDoc)

	// 后台服务
	router.POST(path.New("/service/info"),
		tokenChecker, s.service.Info, s.service.InfoDoc)
	router.POST(path.New("/service/restart/enable"),
		tokenChecker, s.service.CanRestart, s.service.CanRestartDoc)
	router.POST(path.New("/service/restart"),
		tokenChecker, s.service.Restart, s.service.RestartDoc)
	router.POST(path.New("/service/update/enable"),
		tokenChecker, s.service.CanUpdate, s.service.CanUpdateDoc)
	router.POST(path.New("/service/update"),
		tokenChecker, s.service.Update, s.service.UpdateDoc)

	// 更新管理
	router.POST(path.New("/update/enable"),
		tokenChecker, s.update.Enable, s.update.EnableDoc)
	router.POST(path.New("/update/info"),
		tokenChecker, s.update.Info, s.update.InfoDoc)
	router.POST(path.New("/update/restart/enable"),
		tokenChecker, s.update.CanRestart, s.update.CanRestartDoc)
	router.POST(path.New("/update/restart"),
		tokenChecker, s.update.Restart, s.update.RestartDoc)
	router.POST(path.New("/update/upload/enable"),
		tokenChecker, s.update.CanUpdate, s.update.CanUpdateDoc)
	router.POST(path.New("/update/upload"),
		tokenChecker, s.update.Update, s.update.UpdateDoc)

	// 网站管理
	router.POST(path.New("/site/root/info"),
		tokenChecker, s.site.RootInfo, s.site.RootInfoDoc)
	router.POST(path.New("/site/root/file/upload"),
		tokenChecker, s.site.RootUploadFile, s.site.RootUploadFileDoc)
	router.POST(path.New("/site/root/file/delete"),
		tokenChecker, s.site.RootDeleteFile, s.site.RootDeleteFileDoc)
	router.POST(path.New("/site/opt/info"),
		tokenChecker, s.site.OptInfo, s.site.OptInfoDoc)
	router.POST(path.New("/site/opt/upload"),
		tokenChecker, s.site.OptUpload, s.site.OptUploadDoc)
	router.POST(path.New("/site/doc/info"),
		tokenChecker, s.site.DocInfo, s.site.DocInfoDoc)
	router.POST(path.New("/site/doc/upload"),
		tokenChecker, s.site.DocUpload, s.site.DocUploadDoc)
	router.POST(path.New("/site/webapp/info"),
		tokenChecker, s.site.WebappInfo, s.site.WebappInfoDoc)
	router.POST(path.New("/site/webapp/upload"),
		tokenChecker, s.site.WebappUpload, s.site.WebappUploadDoc)
	router.POST(path.New("/site/webapp/delete"),
		tokenChecker, s.site.WebappDelete, s.site.WebappDeleteDoc)

	// Websocket
	// 通知推送
	router.GET(path.New("/websocket/notify").SetTokenPlace(types.TokenPlaceQuery).SetWebSocket(true),
		tokenChecker, s.websocket.Notify, s.websocket.NotifyDoc)
}

func (s *handler) mapOptSite(path types.Path, router types.Router, root string) {
	router.ServeFiles(path.New("/*filepath"), nil, http.Dir(root), nil)
}

func (s *handler) mapWebappSite(path types.Path, router types.Router, root string) {
	router.ServeFiles(path.New("/*filepath"), nil, http.Dir(root), nil)
}

func (s *handler) createTokenForAccountPassword(items []types.TokenAuth, a types.Assistant) (string, types.ErrorCode, error) {
	if s.auth == nil {
		return "", types.ErrInternal, fmt.Errorf("not implement")
	}

	account := ""
	password := ""
	count := len(items)
	for i := 0; i < count; i++ {
		item := items[i]
		if item.Name == "account" {
			account = item.Value
		} else if item.Name == "password" {
			password = item.Value
		}
	}

	model, code, err := s.auth.Authenticate(a, account, password)
	if code != nil {
		return "", code, err
	}

	return model.Token, nil, nil
}
