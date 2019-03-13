package configure

type Https struct {
	Address           string      `json:"address" note:"监听地址，空表示监听所有地址"`
	Port              int         `json:"port" note:"监听端口号"`
	Enabled           bool        `json:"enabled" note:"是否启用"`
	BehindProxy       bool        `json:"behindProxy" note:"是否位于代理服务器之后"`
	RequestClientCert bool        `json:"requestClientCert" note:"是否要求客户端证书"`
	Cert              Certificate `json:"cert" note:"证书"`
}
