package handler

import (
	"encoding/json"
	"fmt"
	"github.com/csby/wsf/router"
	"github.com/csby/wsf/types"
	"net"
	"net/http"
	"time"
)

type httpHandler struct {
	types.Base

	router          *router.Router
	handler         types.HttpHandler
	redirectToHttps bool

	rid types.RandNumber
}

func (s *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	a := s.newAssistant(w, r)
	s.LogDebug("new request: rid=", a.rid,
		", rip=", a.rip,
		", host=", r.Host,
		", schema=", a.schema,
		", method=", r.Method,
		", path=", a.path)

	if s.redirectToHttps {
		if a.schema == "http" {
			if r.Method == "GET" {
				redirectUrl := fmt.Sprintf("https://%s%s", r.Host, a.path)
				http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
				return
			}
		}
	}

	defer func(w http.ResponseWriter, r *http.Request, a *httpAssistant) {
		a.leaveTime = time.Now()
		go s.postRouting(w, r, a)
	}(w, r, a)

	defer func(a *httpAssistant) {
		if err := recover(); err != nil {
			s.LogError(a.schema,
				" request error(rid=", a.rid,
				", schema=", a.schema,
				", path=", a.path,
				"rip=", a.rip,
				"): ", err)

			a.Error(types.ErrException, err)
		}
	}(a)

	if s.preRouting(w, r, a) {
		return
	}

	a.path = r.URL.Path
	s.router.Serve(w, r, a)
}

func (s *httpHandler) preRouting(w http.ResponseWriter, r *http.Request, a *httpAssistant) bool {
	if s.handler == nil {
		return false
	}

	return s.handler.PreRouting(w, r, a)
}

func (s *httpHandler) postRouting(w http.ResponseWriter, r *http.Request, a *httpAssistant) {
	defer func() {
		if err := recover(); err != nil {
			s.LogError("postRouting", err)
		}
	}()

	if s.handler != nil {
		s.handler.PostRouting(w, r, a)
	}
}

func (s *httpHandler) newAssistant(w http.ResponseWriter, r *http.Request) *httpAssistant {
	instance := &httpAssistant{response: w, request: r, schema: "http"}
	instance.method = r.Method
	if r.TLS != nil {
		instance.schema = "https"
	}
	instance.keys = make(map[string]interface{})
	instance.log = false
	instance.enterTime = time.Now()
	instance.path = r.URL.Path
	instance.rid = s.rid.New()
	instance.rip, _, _ = net.SplitHostPort(r.RemoteAddr)
	instance.restart = s.restart()
	instance.token = r.Header.Get("token")
	if instance.token == "" {
		if r.Method == "GET" {
			instance.token = r.FormValue("token")
		}
	}
	if len(r.URL.Query()) > 0 {
		params := make([]*types.Query, 0)
		for k, v := range r.URL.Query() {
			param := &types.Query{Key: k}
			if len(v) > 0 {
				param.Value = v[0]
			}
			params = append(params, param)
		}
		instance.param, _ = json.Marshal(params)
	}

	return instance
}

func (s *httpHandler) restart() func() error {
	if s.handler == nil {
		return nil
	}

	extend := s.handler.Extend()
	if extend == nil {
		return nil
	}

	return extend.Restart()
}
