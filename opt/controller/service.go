package controller

import (
	"bytes"
	"fmt"
	"github.com/csby/wsf/file"
	"github.com/csby/wsf/server/configure"
	"github.com/csby/wsf/types"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Service struct {
	controller

	bootTime time.Time
}

func NewService(log types.Log, cfg *configure.Configure) *Service {
	instance := &Service{}
	instance.SetLog(log)
	instance.cfg = cfg

	if cfg != nil {
		instance.bootTime = cfg.Service.BootTime
	} else {
		instance.bootTime = time.Now()
	}

	return instance
}

func (s *Service) Info(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	data := &types.SvcInfo{BootTime: types.DateTime(s.bootTime)}
	cfg := s.cfg
	if cfg != nil {
		data.Name = cfg.Module.Name
		data.Version = cfg.Module.Version
		data.Remark = cfg.Module.Remark
	}

	a.Success(data)
}

func (s *Service) InfoDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "后台服务")
	function := catalog.AddFunction(method, path, "获取服务信息")
	function.SetNote("获取当前服务信息")
	function.SetOutputDataExample(&types.SvcInfo{
		Name:     "server",
		BootTime: types.DateTime(time.Now()),
		Version:  "1.0.1.0",
		Remark:   "XXX服务",
	})
	function.SetInputContentType("")
}

func (s *Service) CanRestart(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success(a.CanRestart())
}

func (s *Service) CanRestartDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "后台服务")
	function := catalog.AddFunction(method, path, "是否可在线重启")
	function.SetNote("判断当前服务是否可以在线重启")
	function.SetOutputDataExample(true)
	function.SetInputContentType("")
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

func (s *Service) RestartDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "后台服务")
	function := catalog.AddFunction(method, path, "重启服务")
	function.SetNote("重新启动当前服务")
	function.SetOutputDataExample(true)
	function.SetInputContentType("")
}

func (s *Service) CanUpdate(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	a.Success(a.CanUpdate())
}

func (s *Service) CanUpdateDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "后台服务")
	function := catalog.AddFunction(method, path, "是否可在线更新")
	function.SetNote("判断当前服务是否可以在线更新")
	function.SetOutputDataExample(true)
	function.SetInputContentType("")
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

func (s *Service) UpdateDoc(doc types.Doc, method string, path types.HttpPath) {
	_, fileName := filepath.Split(s.cfg.Module.Path)
	note := fmt.Sprintf("安装包(必须包含文件'%s')", fileName)

	catalog := s.createCatalog(doc, "后台服务")
	function := catalog.AddFunction(method, path, "更新服务")
	function.SetNote("上传并更新当前服务")
	function.SetOutputDataExample(nil)
	function.SetInputContentType("multipart/form-data")
	function.AddInputForm(true, "file", note, 1, nil)
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

func (s *Service) copyFile(source, dest string) (int64, error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return 0, err
	}
	defer sourceFile.Close()

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return 0, err
	}

	destFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sourceFileInfo.Mode())
	if err != nil {
		return 0, err
	}
	defer destFile.Close()

	return io.Copy(destFile, sourceFile)
}

func (s *Service) getBinFilePath(folderPath, fileName string) (string, error) {
	paths, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return "", err
	}

	for _, path := range paths {
		if path.IsDir() {
			appPath, err := s.getBinFilePath(filepath.Join(folderPath, path.Name()), fileName)
			if err != nil {
				continue
			}
			return appPath, nil
		} else {
			if path.Name() == fileName {
				return filepath.Join(folderPath, path.Name()), nil
			}
		}
	}

	return "", fmt.Errorf("服务主程序(%s)不存在", fileName)
}
