package controller

import (
	"github.com/csby/wsf/types"
	"net/http"
	"os"
	"time"
)

func (s *Service) CanRestart(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success(a.CanRestart())
}

func (s *Service) Restart(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	if !a.CanRestart() {
		a.Error(types.ErrNotSupport, "当前不在服务模式下运行")
		return
	}

	go func(a types.Assistant) {
		time.Sleep(2 * time.Second)
		err := a.Restart()
		if err != nil {
			s.LogError("重启服务失败:", err)
		}
		os.Exit(1)
	}(a)

	a.Success(true)
}

func (s *Service) CanUpdate(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success(a.CanUpdate())
}

func (s *Service) Update(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	if !s.update(w, r, a) {
		return
	}

	go func(a types.Assistant) {
		time.Sleep(2 * time.Second)
		err := a.Restart()
		if err != nil {
			s.LogError("更新服务后重启失败:", err)
		}
		os.Exit(0)
	}(a)

	a.Success(nil)
}

func (s *Service) update(w http.ResponseWriter, r *http.Request, a types.Assistant) bool {
	newBinFilePath, folder, ok := s.extractUploadFile(w, r, a)
	if !ok {
		if len(folder) > 0 {
			os.RemoveAll(folder)
		}
		return false
	}
	defer os.RemoveAll(folder)

	if !a.CanUpdate() {
		a.Error(types.ErrNotSupport, "服务不支持在线更新")
		return false
	}

	oldBinFilePath := s.cfg.Module.Path
	err := os.Remove(oldBinFilePath)
	if err != nil {
		a.Error(types.ErrInternal, err)
		return false
	}
	_, err = s.copyFile(newBinFilePath, oldBinFilePath)
	if err != nil {
		a.Error(types.ErrInternal, err)
		return false
	}

	return true
}
