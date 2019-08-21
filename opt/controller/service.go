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

	svcMgr   types.SvcUpdMgr
	bootTime time.Time
}

func NewService(log types.Log, cfg *configure.Configure, svcMgr types.SvcUpdMgr) *Service {
	instance := &Service{}
	instance.SetLog(log)
	instance.cfg = cfg
	instance.svcMgr = svcMgr

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
	function.AddOutputError(types.ErrTokenInvalid)
}

func (s *Service) CanRestartDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "后台服务")
	function := catalog.AddFunction(method, path, "是否可在线重启")
	function.SetNote("判断当前服务是否可以在线重启")
	function.SetOutputDataExample(true)
	function.SetInputContentType("")
	function.AddOutputError(types.ErrTokenInvalid)
}

func (s *Service) RestartDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "后台服务")
	function := catalog.AddFunction(method, path, "重启服务")
	function.SetNote("重新启动当前服务")
	function.SetOutputDataExample(true)
	function.SetInputContentType("")
	function.AddOutputError(types.ErrTokenInvalid)
}

func (s *Service) CanUpdateDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "后台服务")
	function := catalog.AddFunction(method, path, "是否可在线更新")
	function.SetNote("判断当前服务是否可以在线更新")
	function.SetOutputDataExample(true)
	function.SetInputContentType("")
	function.AddOutputError(types.ErrTokenInvalid)
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
	function.AddOutputError(types.ErrInternal)
	function.AddOutputError(types.ErrNotSupport)
	function.AddOutputError(types.ErrTokenInvalid)
}

func (s *Service) extractUploadFile(w http.ResponseWriter, r *http.Request, a types.Assistant) (string, string, bool) {
	uploadFile, _, err := r.FormFile("file")
	if err != nil {
		a.Error(types.ErrInputInvalid, "invalid file: ", err)
		return "", "", false
	}
	defer uploadFile.Close()

	buf := &bytes.Buffer{}
	fileSize, err := buf.ReadFrom(uploadFile)
	if err != nil {
		a.Error(types.ErrInputInvalid, "read file error: ", err)
		return "", "", false
	}
	if fileSize < 1 {
		a.Error(types.ErrInputInvalid, "invalid file: size is zero")
		return "", "", false
	}

	binFileFolder, oldBinFileName := filepath.Split(s.cfg.Module.Path)
	tempFolder := filepath.Join(binFileFolder, a.NewGuid())
	err = os.MkdirAll(tempFolder, 0777)
	if err != nil {
		a.Error(types.ErrInternal, fmt.Sprintf("create temp folder '%s' error:", tempFolder), err)
		return "", "", false
	}

	fileData := buf.Bytes()
	zipFile := &file.Zip{}
	err = zipFile.DecompressMemory(fileData, tempFolder)
	if err != nil {
		tarFile := &file.Tar{}
		err = tarFile.DecompressMemory(fileData, tempFolder)
		if err != nil {
			a.Error(types.ErrInternal, "decompress file error: ", err)
			return "", tempFolder, false
		}
	}

	newBinFilePath, err := s.getBinFilePath(tempFolder, oldBinFileName)
	if err != nil {
		a.Error(types.ErrInputInvalid, err)
		return "", tempFolder, false
	}
	module := &types.Module{Path: newBinFilePath}
	moduleName := module.Name()
	if moduleName != s.cfg.Module.Name {
		a.Error(types.ErrInputInvalid, fmt.Sprintf("模块名称(%s)无效", moduleName))
		return "", tempFolder, false
	}
	moduleType := module.Type()
	if moduleType != s.cfg.Module.Type {
		a.Error(types.ErrInputInvalid, fmt.Sprintf("模块名称(%s)无效", moduleType))
		return "", tempFolder, false
	}

	return newBinFilePath, tempFolder, true
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
