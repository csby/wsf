package configure

type Document struct {
	Enabled bool   `json:"enabled" note:"是否启用"`
	Root    string `json:"root" note:"网站物理路径, 如: /home/doc"`
}
