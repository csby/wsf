package controller

import (
	"fmt"
	"github.com/csby/wsf/types"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (s *Update) Enable(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success(true)
}

func (s *Update) Info(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	data := &types.SvcUpdInfo{}
	data.Name = s.serviceName()
	data.Status = 0

	info, err := s.svcMgr.RemoteInfo()
	if err == nil {
		data.Version = info.Version
		data.Remark = info.Remark
		data.BootTime = info.BootTime
		data.Status = 2
	} else {
		status, err := s.svcMgr.Status(data.Name)
		if err == nil {
			if status == types.ServerStatusStopped {
				data.Status = 1
			} else if status == types.ServerStatusRunning {
				data.Status = 2
			}
		}
	}
	a.Success(data)
}

func (s *Update) CanRestart(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	_, err := s.svcMgr.Status(s.serviceName())
	if err != nil {
		a.Success(false)
	} else {
		a.Success(true)
	}
}

func (s *Update) Restart(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	status, err := s.svcMgr.Status(s.serviceName())
	if err != nil {
		a.Error(types.ErrInternal, err)
		return
	}

	if status == types.ServerStatusRunning {
		err = s.svcMgr.Restart(s.serviceName())
	} else {
		err = s.svcMgr.Start(s.serviceName())
	}

	if err != nil {
		a.Error(types.ErrInternal, err)
		return
	}

	a.Success(true)
}

func (s *Update) CanUpdate(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success(s.canUpdate())
}

func (s *Update) Update(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	newBinFilePath, folder, ok := s.extractUploadFile(w, r, a)
	if !ok {
		if len(folder) > 0 {
			os.RemoveAll(folder)
		}
		return
	}
	defer os.RemoveAll(folder)

	if !s.canUpdate() {
		a.Error(types.ErrNotSupport, "服务不支持在线更新")
		return
	}

	svcName := s.serviceName()
	status, err := s.svcMgr.Status(svcName)
	if err != nil {
		if strings.Contains(err.Error(), "Access is denied") {
			a.Error(types.ErrInternal, err)
			return
		}

		binFileFolder, _ := filepath.Split(s.cfg.Module.Path)
		binFilePath := filepath.Join(binFileFolder, s.executeFileName())
		_, err := s.copyFile(newBinFilePath, binFilePath)
		if err != nil {
			a.Error(types.ErrInternal, fmt.Errorf("copy file error: %v", err))
			return
		}

		err = s.svcMgr.Install(svcName, binFilePath)
		if err != nil {
			a.Error(types.ErrInternal, fmt.Errorf("install service '%s' error: %v", svcName, err))
			return
		}

		err = s.svcMgr.Start(svcName)
		if err != nil {
			a.Error(types.ErrInternal, fmt.Errorf("start service '%s' error: %v", svcName, err))
			return
		}
	} else {
		if status != types.ServerStatusRunning {
			err = s.svcMgr.Start(svcName)
			if err != nil {
				a.Error(types.ErrInternal, fmt.Errorf("stop service '%s' error: %v", svcName, err))
				return
			}
		}

		info, err := s.svcMgr.RemoteInfo()
		if err != nil {
			a.Error(types.ErrInternal, fmt.Errorf("get service '%s' info error: %v", svcName, err))
			return
		}

		err = s.svcMgr.Stop(svcName)
		if err != nil {
			a.Error(types.ErrInternal, fmt.Errorf("stop service '%s' error: %v", svcName, err))
			return
		}

		err = os.Remove(info.Path)
		if err != nil {
			a.Error(types.ErrInternal, fmt.Errorf("remove service '%s' execute file error: %v", svcName, err))
			return
		}

		_, err = s.copyFile(newBinFilePath, info.Path)
		if err != nil {
			a.Error(types.ErrInternal, fmt.Errorf("copy file error: %v", err))
			return
		}

		err = s.svcMgr.Start(svcName)
		if err != nil {
			a.Error(types.ErrInternal, fmt.Errorf("start service '%s' error: %v", svcName, err))
			return
		}
	}

	a.Success(true)
}

func (s *Update) executeFileName() string {
	return "wsfupd.exe"
}

func (s *Update) canUpdate() bool {
	info, err := s.svcMgr.RemoteInfo()
	if err == nil {
		if info.Interactive {
			return false
		}
	}

	return true
}
