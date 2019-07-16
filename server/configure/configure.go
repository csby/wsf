package configure

type Configure struct {
	Module    Module    `json:"-" note:"模块信息"`
	Log       Log       `json:"log" note:"日志"`
	Service   Service   `json:"svc" note:"系统服务"`
	Http      Http      `json:"http" note:"HTTP服务"`
	Https     Https     `json:"https" note:"HTTPS服务"`
	Tcp       Tcp       `json:"tcp" note:"TCP(TLS)服务"`
	Proxy     string    `json:"proxy" note:"代理服务器IP地址（客户端不是来自代理服务器时，远程地址为当前连接地址）"`
	Root      string    `json:"root" note:"网站根目录"`
	Document  Document  `json:"document" note:"接口文档"`
	Operation Operation `json:"operation" note:"服务管理"`
	Webapp    Webapp    `json:"webapp" note:"应用网站"`
}
