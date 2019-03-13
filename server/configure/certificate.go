package configure

type Certificate struct {
	Ca     CertificateCa  `json:"ca" note:"CA证书"`
	Server CertificatePfx `json:"server" note:"服务器证书"`
}
