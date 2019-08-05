package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/csby/wsf/doc/web"
	"github.com/csby/wsf/file"
	"github.com/csby/wsf/server/configure"
	"github.com/csby/wsf/types"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Site struct {
	controller

	optPath    string
	webappPath string
}

func NewSite(log types.Log, cfg *configure.Configure, db types.TokenDatabase, chs types.SocketChannelCollection, optPath, webappPath string) *Site {
	instance := &Site{}
	instance.SetLog(log)
	instance.cfg = cfg
	instance.dbToken = db
	instance.wsChannels = chs
	instance.optPath = optPath
	instance.webappPath = webappPath

	return instance
}

func (s *Site) RootInfo(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	data := make([]*types.SiteFile, 0)

	files, err := ioutil.ReadDir(s.cfg.Root)
	if err != nil {
		a.Error(types.ErrInternal, err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		item := &types.SiteFile{
			Name:       file.Name(),
			UploadTime: types.DateTime(file.ModTime()),
			Url:        fmt.Sprintf("%s://%s/%s", a.Schema(), r.Host, file.Name()),
		}

		data = append(data, item)
	}

	a.Success(data)
}

func (s *Site) RootInfoDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "网站管理", "根站点")
	function := catalog.AddFunction(method, path, "获取文件列表")
	function.SetNote("获取根站点所有文件列表，但不包括文件夹")
	function.SetOutputDataExample([]types.SiteFile{
		{
			Url:        "http://192.168.1.1:8080/test.txt",
			Name:       "test.txt",
			UploadTime: types.DateTime(time.Now()),
		},
	})
	function.SetInputContentType("")
}

func (s *Site) RootUploadFile(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	uploadFile, head, err := r.FormFile("file")
	if err != nil {
		a.Error(types.ErrInputInvalid, "file invalid: ", err)
		return
	}
	defer uploadFile.Close()
	buf := &bytes.Buffer{}
	fileSize, err := buf.ReadFrom(uploadFile)
	if err != nil {
		a.Error(types.ErrInputInvalid, "read file error: ", err)
		return
	}
	if fileSize < 1 {
		a.Error(types.ErrInputInvalid, "file size is zero")
		return
	}

	folderPath := s.cfg.Root
	err = os.MkdirAll(folderPath, 0777)
	if err != nil {
		a.Error(types.ErrInternal, "create folder for root error: ", err)
		return
	}

	filePath := filepath.Join(folderPath, head.Filename)
	fileWriter, err := os.Create(filePath)
	if err != nil {
		a.Error(types.ErrInternal, "create file in root folder error: ", err)
		return
	}
	defer fileWriter.Close()

	_, err = io.Copy(fileWriter, uploadFile)
	if err != nil {
		a.Error(types.ErrInternal, "save file in root folder error: ", err)
		return
	}

	data := &types.SiteFile{
		Name:       head.Filename,
		UploadTime: types.DateTime(time.Now()),
		Url:        fmt.Sprintf("%s://%s/%s", a.Schema(), r.Host, head.Filename),
	}

	a.Success(data)

	s.writeWebSocketMessage(a.Token(), types.WSRootSiteUploadFile, data)
}

func (s *Site) RootUploadFileDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "网站管理", "根站点")
	function := catalog.AddFunction(method, path, "上传文件")
	function.SetNote("上传文件到根站点所在目录，成功返回已上传文件的访问地址")
	function.SetOutputDataExample(string("http://192.168.1.1:8080/test.txt"))
	function.SetInputContentType("multipart/form-data")
	function.AddInputForm(true, "file", "文件", 1, nil)
}

func (s *Site) RootDeleteFile(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	argument := &types.SiteFileFilter{}
	err := a.GetJson(argument)
	if err != nil {
		a.Error(types.ErrInputInvalid, err)
		return
	}
	if len(argument.Name) < 1 {
		a.Error(types.ErrInputInvalid, "文件名称(name)文空")
		return
	}

	filePath := filepath.Join(s.cfg.Root, argument.Name)
	err = os.Remove(filePath)
	if err != nil {
		a.Error(types.ErrInputInvalid, err)
		return
	}

	a.Success(argument)

	s.writeWebSocketMessage(a.Token(), types.WSRootSiteDeleteFile, argument)
}

func (s *Site) RootDeleteFileDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "网站管理", "根站点")
	function := catalog.AddFunction(method, path, "删除文件")
	function.SetNote("删除根站点所在目录的文件")
	function.SetInputExample(&types.SiteFileFilter{
		Name: "test.txt",
	})
	function.SetOutputDataExample(nil)
}

func (s *Site) OptInfo(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	data, err := s.getInfo(s.cfg.Operation.Root, s.optPath, r.Host, a.Schema())
	if err != nil {
		a.Error(types.ErrInternal, err)
		return
	}

	a.Success(data)
}

func (s *Site) OptInfoDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "网站管理", "后台服务")
	function := catalog.AddFunction(method, path, "获取管理网站信息")
	function.SetNote("获取管理网站信息，包括访问地址及版本等")
	function.SetOutputDataExample(&types.SiteInfo{
		Url:        "http://192.168.1.1:8080/opt/",
		Version:    "1.0.1.0",
		DeployTime: types.DateTime(time.Now()),
	})
	function.SetInputContentType("")
}

func (s *Site) OptUpload(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	root := s.cfg.Operation.Root
	if !s.upload(root, w, r, p, a) {
		return
	}

	path := s.optPath
	data, _ := s.getInfo(root, path, r.Host, a.Schema())
	a.Success(data)

	if data != nil {
		s.writeWebSocketMessage(a.Token(), types.WSOptSiteUpload, data)
	}
}

func (s *Site) OptUploadDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "网站管理", "后台服务")
	function := catalog.AddFunction(method, path, "上传管理网站")
	function.SetNote("上传网站打包文件(.zip或.tar.gz)，并替换之前已发布的网站")
	function.SetOutputDataExample(&types.SiteInfo{
		Url:        "http://192.168.1.1:8080/api/info",
		Version:    "1.0.1.0",
		DeployTime: types.DateTime(time.Now()),
	})
	function.SetInputContentType("multipart/form-data")
	function.AddInputForm(true, "file", "网站打包文件(.zip或.tar.gz)", 1, nil)
}

func (s *Site) DocInfo(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	data, err := s.getInfo(s.cfg.Document.Root, web.SitePath, r.Host, a.Schema())
	if err != nil {
		a.Error(types.ErrInternal, err)
		return
	}

	a.Success(data)
}

func (s *Site) DocInfoDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "网站管理", "接口文档")
	function := catalog.AddFunction(method, path, "获取接口文档网站信息")
	function.SetNote("获取接口文档网站信息，包括访问地址及版本等")
	function.SetOutputDataExample(&types.SiteInfo{
		Url:        "http://192.168.1.1:8080/doc/",
		Version:    "1.0.1.0",
		DeployTime: types.DateTime(time.Now()),
	})
	function.SetInputContentType("")
}

func (s *Site) DocUpload(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	root := s.cfg.Document.Root
	if !s.upload(root, w, r, p, a) {
		return
	}

	path := web.SitePath
	data, _ := s.getInfo(root, path, r.Host, a.Schema())
	a.Success(data)

	if data != nil {
		s.writeWebSocketMessage(a.Token(), types.WSDocSiteUpload, data)
	}
}

func (s *Site) DocUploadDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "网站管理", "接口文档")
	function := catalog.AddFunction(method, path, "上传接口文档网站")
	function.SetNote("上传网站打包文件(.zip或.tar.gz)，并替换之前已发布的网站")
	function.SetOutputDataExample(&types.SiteInfo{
		Url:        "http://192.168.1.1:8080/doc/",
		Version:    "1.0.1.0",
		DeployTime: types.DateTime(time.Now()),
	})
	function.SetInputContentType("multipart/form-data")
	function.AddInputForm(true, "file", "网站打包文件(.zip或.tar.gz)", 1, nil)
}

func (s *Site) WebappInfo(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	baseUrl := fmt.Sprintf("%s://%s%s", a.Schema(), r.Host, s.webappPath)
	data := &types.SiteAppTree{}
	data.ParseChildren(s.cfg.Webapp.Root, baseUrl)

	a.Success(data.Children)
}

func (s *Site) WebappInfoDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "网站管理", "网站应用")
	function := catalog.AddFunction(method, path, "获取网站应用信息")
	function.SetNote("获取所有网站应用的信息")
	function.SetOutputDataExample([]types.SiteAppTree{
		{
			Type:       1,
			Path:       "test",
			UploadTime: types.DateTime(time.Now()),
			UploadUser: "admin",
			Version:    "1.0.1.0",
			Url:        "https://www.example.com/webapp/test",
		},
	})
	function.SetInputContentType("")
}

func (s *Site) WebappUpload(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	path := strings.TrimSpace(strings.ToLower(r.FormValue("path")))
	if len(path) < 1 {
		a.Error(types.ErrInputInvalid, "路径(path)为空")
		return
	}

	uploadFile, _, err := r.FormFile("file")
	if err != nil {
		a.Error(types.ErrInputInvalid, "file invalid: ", err)
		return
	}
	defer uploadFile.Close()
	buf := &bytes.Buffer{}
	fileSize, err := buf.ReadFrom(uploadFile)
	if err != nil {
		a.Error(types.ErrInputInvalid, "read file error: ", err)
		return
	}
	if fileSize < 1 {
		a.Error(types.ErrInputInvalid, "file size is zero")
		return
	}

	tempFolder := filepath.Join(s.cfg.Webapp.Root, a.NewGuid())
	err = os.MkdirAll(tempFolder, 0777)
	if err != nil {
		a.Error(types.ErrInternal, "create temp folder error: ", err)
		return
	}
	fileData := buf.Bytes()
	zipFile := &file.Zip{}
	err = zipFile.DecompressMemory(fileData, tempFolder)
	if err != nil {
		tarFile := &file.Tar{}
		err = tarFile.DecompressMemory(fileData, tempFolder)
		if err != nil {
			a.Error(types.ErrInternal, "decompress file error: ", err)
			return
		}
	}

	appFolder := filepath.Join(s.cfg.Webapp.Root, path)
	err = os.RemoveAll(appFolder)
	if err != nil {
		a.Error(types.ErrInternal, fmt.Sprintf("remove original webapp '%s' error: ", path), err)
		return
	}
	os.MkdirAll(filepath.Dir(appFolder), 0777)
	err = os.Rename(tempFolder, appFolder)
	if err != nil {
		a.Error(types.ErrInternal, fmt.Sprintf("rename webapp folder '%s' error: ", path), err)
		return
	}

	appInfo := &types.SiteApp{
		UploadTime: types.DateTime(time.Now()),
	}
	appInfo.Remark = r.FormValue("remark")
	appInfo.Version = r.FormValue("version")
	token := s.getToken(a.Token())
	if token != nil {
		appInfo.UploadUser = token.UserAccount
	}
	appInfo.SaveToFolder(appFolder)

	a.Success(appInfo)

	s.writeWebSocketMessage(a.Token(), types.WSWebappSiteUpload, nil)
}

func (s *Site) WebappUploadDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "网站管理", "网站应用")
	function := catalog.AddFunction(method, path, "上传网站应用")
	function.SetNote("上传网站打包文件(.zip或.tar.gz)，并替换之前已发布的网站")
	function.SetOutputDataExample(&types.SiteApp{
		Version:    "1.0.1.0",
		Remark:     "说明",
		UploadUser: "Admin",
		UploadTime: types.DateTime(time.Now()),
	})
	function.SetInputContentType("multipart/form-data")
	function.AddInputForm(true, "file", "网站打包文件(.zip或.tar.gz)", 1, nil)
	function.AddInputForm(true, "path", "路径，如test或group/item", 0, "")
	function.AddInputForm(false, "version", "版本，如1.0.1.0", 0, "")
	function.AddInputForm(false, "remark", "说明，如XXX网站原型", 0, "")
}

func (s *Site) WebappDelete(w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) {
	argument := &types.SiteAppPath{}
	err := a.GetJson(argument)
	if err != nil {
		a.Error(types.ErrInputInvalid, err)
		return
	}
	if len(argument.Path) < 1 {
		a.Error(types.ErrInputInvalid, "路径(path)文空")
		return
	}

	filePath := filepath.Join(s.cfg.Webapp.Root, argument.Path)
	err = os.RemoveAll(filePath)
	if err != nil {
		a.Error(types.ErrInputInvalid, err)
		return
	}

	a.Success(argument)

	s.writeWebSocketMessage(a.Token(), types.WSWebappSiteDelete, nil)
}

func (s *Site) WebappDeleteDoc(doc types.Doc, method string, path types.HttpPath) {
	catalog := s.createCatalog(doc, "网站管理", "网站应用")
	function := catalog.AddFunction(method, path, "删除网站应用")
	function.SetNote("删除指定路径的应用网站")
	function.SetInputExample(&types.SiteAppPath{
		Path: "group/item",
	})
	function.SetOutputDataExample(nil)
}

func (s *Site) upload(root string, w http.ResponseWriter, r *http.Request, p types.Params, a types.Assistant) bool {
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

	tempFolder := filepath.Join(filepath.Dir(root), a.NewGuid())
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

	uploadFolder := root
	err = os.RemoveAll(uploadFolder)
	if err != nil {
		a.Error(types.ErrInternal, "remove original site error: ", err)
		return false
	}
	os.MkdirAll(filepath.Dir(uploadFolder), 0777)
	err = os.Rename(tempFolder, uploadFolder)
	if err != nil {
		a.Error(types.ErrInternal, fmt.Sprintf("rename folder '%s' error:", uploadFolder), err)
		return false
	}

	return true
}

func (s *Site) getInfo(root, path, host, schema string) (*types.SiteInfo, error) {
	fi, err := os.Stat(root)
	if os.IsNotExist(err) {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("the root '%s' is not folder", root)
	}

	info := &types.SiteInfo{DeployTime: types.DateTime(fi.ModTime())}
	info.Version, _ = s.getVersion(root)
	info.Url = fmt.Sprintf("%s://%s%s/", schema, host, path)
	info.Root = root

	return info, nil
}

func (s *Site) getVersion(folderPath string) (string, error) {
	/*
		const version = "1.0.1.1"

		export default {
			version
		}
	*/
	filePath := filepath.Join(folderPath, "version.js")
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	bufReader := bufio.NewReader(file)
	for {
		line, err := bufReader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if len(line) <= 0 {
			continue
		}
		if line[0] == '/' {
			continue
		}

		keyValue := strings.Split(line, "=")
		if len(keyValue) < 2 {
			continue
		}
		if strings.TrimSpace(keyValue[0]) != "const version" {
			continue
		}
		value := strings.TrimSpace(keyValue[1])
		value = strings.TrimLeft(value, "'")
		value = strings.TrimLeft(value, "\"")
		value = strings.TrimRight(value, "'")
		value = strings.TrimRight(value, "\"")

		return value, nil
	}

	return "", nil
}
