package controller

import (
	"fmt"
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

func (s *Service) CanRestartDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "后台服务")
	function := catalog.AddFunction(method, path, "是否可在线重启")
	function.SetNote("判断当前服务是否可以在线重启")
	function.SetOutputDataExample(true)
	function.SetInputContentType("")
}

func (s *Service) RestartDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "后台服务")
	function := catalog.AddFunction(method, path, "重启服务")
	function.SetNote("重新启动当前服务")
	function.SetOutputDataExample(true)
	function.SetInputContentType("")
}

func (s *Service) CanUpdateDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "后台服务")
	function := catalog.AddFunction(method, path, "是否可在线更新")
	function.SetNote("判断当前服务是否可以在线更新")
	function.SetOutputDataExample(true)
	function.SetInputContentType("")
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
