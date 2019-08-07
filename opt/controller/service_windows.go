package controller

import (
	"github.com/csby/wsf/types"
	"github.com/kardianos/service"
	"net/http"
	"os"
)

func (s *Service) CanRestart(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success(s.canRestart())
}

func (s *Service) Restart(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	if !s.canRestart() {
		a.Error(types.ErrNotSupport)
		return
	}

	err := s.svcMgr.RemoteRestart(s.cfg.Service.Name)
	if err != nil {
		a.Error(types.ErrInternal, err)
		return
	}

	a.Success(true)
}

func (s *Service) CanUpdate(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success(s.canRestart())
}

func (s *Service) Update(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	newBinFilePath, folder, ok := s.extractUploadFile(w, r, a)
	if !ok {
		if len(folder) > 0 {
			os.RemoveAll(folder)
		}
		return
	}

	if !s.canRestart() {
		os.RemoveAll(folder)
		a.Error(types.ErrNotSupport, "服务不支持在线更新")
		return
	}

	svcName := s.cfg.Service.Name
	svcPath := s.cfg.Module.Path
	err := s.svcMgr.RemoteUpdate(svcName, svcPath, newBinFilePath, folder)
	if err != nil {
		os.RemoveAll(folder)
		a.Error(types.ErrInternal, err)
		return
	}

	a.Success(true)
}

func (s *Service) canRestart() bool {
	if service.Interactive() {
		return false
	}

	_, err := s.svcMgr.RemoteInfo()
	if err != nil {
		return false
	}

	return true
}
