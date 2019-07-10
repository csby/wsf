package configure

type Configure struct {
	Http      Http      `json:"http" note:"HTTP服务"`
	Https     Https     `json:"https" note:"HTTPS服务"`
	Tcp       Tcp       `json:"tcp" note:"TCP(TLS)服务"`
	Proxy     string    `json:"proxy" note:"代理服务器IP地址（客户端不是来自代理服务器时，远程地址为当前连接地址）"`
	Root      string    `json:"root" note:"网站根目录"`
	Document  Document  `json:"document" note:"接口文档"`
	Operation Operation `json:"operation" note:"服务管理"`
	Webapp    Webapp    `json:"webapp" note:"应用网站"`
}
