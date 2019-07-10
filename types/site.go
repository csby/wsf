package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
)

type SiteInfo struct {
	Url        string   `json:"url" note:"访问地址"`
	Version    string   `json:"version" note:"版本号"`
	DeployTime DateTime `json:"deployTime" note:"发布时间"`
}

type SiteFile struct {
	Url        string   `json:"url" note:"访问地址"`
	Name       string   `json:"name" note:"文件名称"`
	UploadTime DateTime `json:"uploadTime" note:"上传时间"`
}

type SiteApp struct {
	mutex sync.RWMutex

	Version    string   `json:"version" note:"版本号"`
	UploadTime DateTime `json:"uploadTime" note:"上传时间"`
	UploadUser string   `json:"uploadUser" note:"上传者账号"`
	Remark     string   `json:"remark" note:"说明"`
}

func (s *SiteApp) LoadFromFile(filePath string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

type SiteAppTree struct {
	parent *SiteAppTree

	Type       int      `json:"type" note:"0-folder; 1-app"`
	Name       string   `json:"name" note:"应用程序名称"`
	Path       string   `json:"path" note:"应用程序路径"`
	Url        string   `json:"url" note:"访问地址"`
	Version    string   `json:"version" note:"版本号"`
	UploadTime DateTime `json:"uploadTime" note:"上传时间"`
	UploadUser string   `json:"uploadUser" note:"上传者账号"`
	Remark     string   `json:"remark" note:"说明"`

	Children []*SiteAppTree `json:"children"`
}

func (s *SiteAppTree) ParseChildren(folderPath, baseUrl string) {
	s.Children = make([]*SiteAppTree, 0)

	paths, err := ioutil.ReadDir(folderPath)
	if err == nil {
		for _, path := range paths {
			if !path.IsDir() {
				continue
			}

			child := &SiteAppTree{
				Name:       path.Name(),
				parent:     s,
				Type:       0,
				UploadTime: DateTime(path.ModTime()),
			}
			child.Path = child.getPath()
			s.Children = append(s.Children, child)

			appInfo := filepath.Join(folderPath, path.Name(), "app.info")
			app := &SiteApp{
				Version:    "1.0.1.0",
				UploadTime: DateTime(path.ModTime()),
			}
			err = app.LoadFromFile(appInfo)
			if err == nil {
				child.Type = 1
				child.Version = app.Version
				child.UploadUser = app.UploadUser
				child.Remark = app.Remark
				child.Url = fmt.Sprintf("%s/%s", baseUrl, child.Path)
				child.Children = make([]*SiteAppTree, 0)
			} else {
				child.ParseChildren(filepath.Join(folderPath, path.Name()), baseUrl)
			}
		}
	}
}

func (s *SiteAppTree) getPath() string {
	path := s.Name

	parent := s.parent
	for parent != nil {
		if parent.Name == "" {
			break
		}

		path = parent.Name + "/" + path
		parent = parent.parent
	}

	return path
}
