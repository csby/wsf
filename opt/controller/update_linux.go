package controller

import (
	"github.com/csby/wsf/types"
	"net/http"
)

func (s *Update) Enable(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success(false)
}

func (s *Update) Info(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Error(types.ErrNotSupport)
}

func (s *Update) CanRestart(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success(false)
}

func (s *Update) Restart(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Error(types.ErrNotSupport)
}

func (s *Update) CanUpdate(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success(false)
}

func (s *Update) Update(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Error(types.ErrNotSupport)
}

func (s *Update) executeFileName() string {
	return "wsfupd"
}
