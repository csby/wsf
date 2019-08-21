package controller

import (
	"github.com/csby/wsf/server/configure"
	"github.com/csby/wsf/types"
	"github.com/mojocn/base64Captcha"
	"net/http"
	"strings"
	"time"
)

type Auth struct {
	controller

	errorCount map[string]int
	ldap       Ldap
}

func NewAuth(log types.Log, cfg *configure.Configure, db types.TokenDatabase, chs types.SocketChannelCollection) *Auth {
	instance := &Auth{}
	instance.SetLog(log)
	instance.cfg = cfg
	instance.dbToken = db
	instance.wsChannels = chs
	instance.errorCount = make(map[string]int)

	if cfg != nil {
		instance.ldap.Enable = cfg.Operation.Ldap.Enable
		instance.ldap.Host = cfg.Operation.Ldap.Host
		instance.ldap.Port = cfg.Operation.Ldap.Port
		instance.ldap.Base = cfg.Operation.Ldap.Base
	}

	if chs != nil {
		chs.AddFilter(instance.onWebsocketWriteFilter)
	}

	return instance
}

func (s *Auth) GetCaptcha(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	filter := &types.CaptchaFilter{
		Mode:   base64Captcha.CaptchaModeNumberAlphabet,
		Length: 4,
		Width:  100,
		Height: 30,
	}
	err := a.GetJson(filter)
	if err != nil {
		a.Error(types.ErrInput, err)
		return
	}

	captchaConfig := base64Captcha.ConfigCharacter{
		Mode:               filter.Mode,
		Height:             filter.Height,
		Width:              filter.Width,
		CaptchaLen:         filter.Length,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		IsUseSimpleFont:    true,
	}
	captchaId, captchaValue := base64Captcha.GenerateCaptcha("", captchaConfig)

	data := &types.Captcha{
		ID:           captchaId,
		Value:        base64Captcha.CaptchaWriteToBase64Encoding(captchaValue),
		Required:     s.captchaRequired(a.RIP()),
		RsaPublicKey: a.RSAPublicKey(),
	}

	a.Success(data)
}

func (s *Auth) GetCaptchaDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "权限管理")
	function := catalog.AddFunction(method, path, "获取验证码")
	function.SetNote("获取用户登陆需要的验证码信息")
	function.SetInputExample(&types.CaptchaFilter{
		Mode:   base64Captcha.CaptchaModeNumberAlphabet,
		Length: 4,
		Width:  100,
		Height: 30,
	})
	function.SetOutputDataExample(&types.Captcha{
		ID:           "GKSVhVMRAHsyVuXSrMYs",
		Value:        "data:image/png;base64,iVBOR...",
		RsaPublicKey: "-----BEGIN PUBLIC KEY-----...-----END PUBLIC KEY-----",
		Required:     false,
	})
	function.AddOutputError(types.ErrInput)
}

func (s *Auth) Login(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	filter := &types.LoginFilter{}
	err := a.GetJson(filter)
	if err != nil {
		a.Error(types.ErrInput, err)
		return
	}

	requireCaptcha := s.captchaRequired(a.RIP())
	err = filter.Check(requireCaptcha)
	if err != nil {
		a.Error(types.ErrInputInvalid, err)
		return
	}

	if requireCaptcha {
		if !base64Captcha.VerifyCaptcha(filter.CaptchaId, filter.CaptchaValue) {
			a.Error(types.ErrLoginCaptchaInvalid)
			return
		}
	}

	pwd := filter.Password
	if strings.ToLower(filter.Encryption) == "rsa" {
		decryptedPwd, err := a.RSADecrypt(pwd)
		if err != nil {
			a.Error(types.ErrLoginPasswordInvalid, err)
			s.increaseErrorCount(a.RIP())
			return
		}
		pwd = string(decryptedPwd)
	}

	login, be, err := s.Authenticate(a, filter.Account, pwd)
	if be != nil {
		a.Error(be, err)
		s.increaseErrorCount(a.RIP())
		return
	}

	a.Success(login)
	s.clearErrorCount(a.RIP())
}

func (s *Auth) LoginDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "权限管理")
	function := catalog.AddFunction(method, path, "用户登录")
	function.SetNote("通过用户账号及密码进行登录获取凭证")
	function.SetInputExample(&types.LoginFilter{
		Account:      "admin",
		Password:     "1",
		CaptchaId:    "r4kcmz2E12e0qJQOvqRB",
		CaptchaValue: "1e35",
		Encryption:   "",
	})
	function.SetOutputDataExample(&types.Login{
		Token: "71b9b7e2ac6d4166b18f414942ff3481",
	})
	function.AddOutputError(types.ErrInput)
	function.AddOutputError(types.ErrInputInvalid)
	function.AddOutputError(types.ErrLoginCaptchaInvalid)
	function.AddOutputError(types.ErrLoginPasswordInvalid)
}

func (s *Auth) Logout(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	tv := a.Token()
	if len(tv) < 1 {
		a.Error(types.ErrTokenEmpty)
		return
	}
	_, ok := s.dbToken.Get(tv, false)
	if !ok {
		a.Error(types.ErrTokenInvalid)
		return
	}

	s.writeWebSocketMessage(a.Token(), types.WSOptUserLogout, nil)
	ok = s.dbToken.Del(tv)
	if ok {
	}

	a.Success(nil)
}

func (s *Auth) LogoutDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "权限管理")
	function := catalog.AddFunction(method, path, "退出登录")
	function.SetNote("退出登录, 使当前凭证失效")
	function.SetOutputDataExample(nil)
	function.SetInputContentType("")
	function.AddOutputError(types.ErrTokenInvalid)
}

func (s *Auth) GetLoginAccount(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	token := s.getToken(a.Token())
	if token == nil {
		a.Error(types.ErrInternal, "凭证无效")
		return
	}
	a.Success(&types.LoginAccount{
		Account:   token.UserAccount,
		Name:      token.UserName,
		LoginTime: types.DateTime(token.LoginTime),
	})
}

func (s *Auth) GetLoginAccountDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "权限管理")
	function := catalog.AddFunction(method, path, "获取登录账号")
	function.SetNote("获取当前登录账号基本信息")
	function.SetOutputDataExample(&types.LoginAccount{
		Account:   "admin",
		Name:      "管理员",
		LoginTime: types.DateTime(time.Now()),
	})
	function.SetInputContentType("")
	function.AddOutputError(types.ErrInternal)
	function.AddOutputError(types.ErrTokenInvalid)
}

func (s *Auth) GetOnlineUsers(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	if s.wsChannels == nil {
		a.Error(types.ErrInternal, "websocket channels is nil")
		return
	}
	a.Success(s.wsChannels.OnlineUsers())
}

func (s *Auth) GetOnlineUsersDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "权限管理")
	function := catalog.AddFunction(method, path, "获取在线用户")
	function.SetNote("获取当前所有在线用户")
	function.SetOutputDataExample([]types.OnlineUser{
		{
			UserAccount:   "admin",
			UserName:      "管理员",
			LoginIP:       "192.168.1.8",
			LoginTime:     types.DateTime(time.Now()),
			LoginDuration: "7秒",
		},
	})
	function.SetInputContentType("")
	function.AddOutputError(types.ErrInternal)
	function.AddOutputError(types.ErrTokenInvalid)
}

func (s *Auth) Authenticate(a types.Assistant, account, password string) (*types.Login, types.ErrorCode, error) {
	act := strings.ToLower(account)
	pwd := password

	var user *configure.User = nil
	userCount := len(s.cfg.Operation.Users)
	for index := 0; index < userCount; index++ {
		if act == strings.ToLower(s.cfg.Operation.Users[index].Account) {
			user = &s.cfg.Operation.Users[index]
			break
		}
	}

	var err error = nil
	userName := account
	if user != nil {
		if pwd != user.Password {
			return nil, types.ErrLoginPasswordInvalid, nil
		}
	} else {
		if s.ldap.Enable {
			err = s.ldap.Authenticate(account, password)
			if err != nil {
				return nil, types.ErrLoginAccountOrPasswordInvalid, err
			}
		} else {
			return nil, types.ErrLoginAccountNotExit, nil
		}
	}

	now := time.Now()
	token := &types.Token{
		ID:          a.NewGuid(),
		UserAccount: account,
		UserName:    userName,
		LoginIP:     a.RIP(),
		LoginTime:   now,
		ActiveTime:  now,
		Usage:       0,
	}
	s.dbToken.Set(token.ID, token)

	login := &types.Login{
		Token:   token.ID,
		Account: token.UserAccount,
		Name:    token.UserName,
	}

	return login, nil, err
}

func (s *Auth) CheckToken(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) bool {
	tokenValue := a.Token()
	if len(tokenValue) < 1 {
		a.Error(types.ErrTokenEmpty)
		return true
	}

	token, ok := s.dbToken.Get(tokenValue, true)
	if !ok {
		a.Error(types.ErrTokenInvalid)
		return true
	}

	tokenModel, ok := token.(*types.Token)
	if !ok {
		a.Error(types.ErrInternal, "类型转换错误(*types.Token)")
		return true
	}

	if tokenModel.LoginIP != a.RIP() {
		a.Error(types.ErrTokenIllegal, "IP不匹配")
		return true
	}

	return false
}

func (s *Auth) CreateTokenForAccountPassword(items []types.TokenAuth, a types.Assistant) (string, types.ErrorCode, error) {
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

	model, code, err := s.Authenticate(a, account, password)
	if code != nil {
		return "", code, err
	}

	return model.Token, nil, nil
}

func (s *Auth) onWebsocketWriteFilter(message *types.SocketMessage, channel types.SocketChannel, token *types.Token) bool {
	if message == nil {
		return false
	}

	if channel == nil {
		return false
	}

	if token == nil {
		return false
	}
	channelToken := channel.Token()
	if channelToken == nil {
		return false
	}

	if message.ID == types.WSOptUserLogin {
		if channelToken.ID == token.ID {
			return true
		}
	} else if message.ID == types.WSOptUserLogout {
		if channelToken.ID != token.ID {
			return true
		}
	}

	return false
}

func (s *Auth) captchaRequired(ip string) bool {
	if s.errorCount == nil {
		return false
	}

	count, ok := s.errorCount[ip]
	if ok {
		if count < 3 {
			return false
		} else {
			return true
		}
	}

	return false
}

func (s *Auth) increaseErrorCount(ip string) {
	if s.errorCount == nil {
		return
	}

	count := 1
	v, ok := s.errorCount[ip]
	if ok {
		count += v
	}

	s.errorCount[ip] = count
}

func (s *Auth) clearErrorCount(ip string) {
	if s.errorCount == nil {
		return
	}

	_, ok := s.errorCount[ip]
	if ok {
		delete(s.errorCount, ip)
	}
}
