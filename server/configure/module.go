package configure

type Module struct {
	Type    string `json:"type" note:"类型"`
	Name    string `json:"name" note:"名称"`
	Path    string `json:"path" note:"路径"`
	Version string `json:"version" note:"版本号"`
	Remark  string `json:"remark" note:"备注说明"`
}
