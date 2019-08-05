package controller

import (
	"bytes"
	"fmt"
	"github.com/csby/wsf/file"
	"github.com/csby/wsf/types"
	"net/http"
	"os"
	"path/filepath"
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
	if !a.CanUpdate() {
		a.Error(types.ErrNotSupport, "服务不支持在线更新")
		return false
	}

	uploadFile, _, err := r.FormFile("file")
	if err != nil {
		a.Error(types.ErrInputInvalid, "invalid file: ", err)
		return false
	}
	defer uploadFile.Close()

	buf := &bytes.Buffer{}
	fileSize, err := buf.ReadFrom(uploadFile)
	if err != nil {
		a.Error(types.ErrInputInvalid, "read file error: ", err)
		return false
	}
	if fileSize < 1 {
		a.Error(types.ErrInputInvalid, "invalid file: size is zero")
		return false
	}

	oldBinFilePath := s.cfg.Module.Path
	oldBinFileFolder, oldBinFileName := filepath.Split(oldBinFilePath)
	tempFolder := filepath.Join(oldBinFileFolder, a.NewGuid())
	err = os.MkdirAll(tempFolder, 0777)
	if err != nil {
		a.Error(types.ErrInternal, fmt.Sprintf("create temp folder '%s' error:", tempFolder), err)
		return false
	}
	defer os.RemoveAll(tempFolder)

	fileData := buf.Bytes()
	zipFile := &file.Zip{}
	err = zipFile.DecompressMemory(fileData, tempFolder)
	if err != nil {
		tarFile := &file.Tar{}
		err = tarFile.DecompressMemory(fileData, tempFolder)
		if err != nil {
			a.Error(types.ErrInternal, "decompress file error: ", err)
			return false
		}
	}

	newBinFilePath, err := s.getBinFilePath(tempFolder, oldBinFileName)
	if err != nil {
		a.Error(types.ErrInputInvalid, err)
		return false
	}
	module := &types.Module{Path: newBinFilePath}
	moduleName := module.Name()
	if moduleName != s.cfg.Module.Name {
		a.Error(types.ErrInputInvalid, fmt.Sprintf("模块名称(%s)无效", moduleName))
		return false
	}
	moduleType := module.Type()
	if moduleType != s.cfg.Module.Type {
		a.Error(types.ErrInputInvalid, fmt.Sprintf("模块名称(%s)无效", moduleType))
		return false
	}

	err = os.Remove(oldBinFilePath)
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
