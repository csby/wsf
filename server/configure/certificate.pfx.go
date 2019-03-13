package configure

type CertificatePfx struct {
	File     string `json:"file" note:"证书文件路径"`
	Password string `json:"password" note:"证书秘密"`
}
