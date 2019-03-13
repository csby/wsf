package example

import "github.com/csby/wsf/server/configure"

var cfg = &Config{
	Server: configure.Configure{
		Http: configure.Http{
			Enabled: true,
			Port:    8086,
		},
		Https: configure.Https{
			Enabled: true,
			Port:    8446,
			Cert: configure.Certificate{
				Ca: configure.CertificateCa{
					File: caCrtFilePath,
				},
				Server: configure.CertificatePfx{
					File:     serverCrtFilePath,
					Password: serverCrtPassword,
				},
			},
			RequestClientCert: true,
		},
		Document: configure.Document{
			Enabled: true,
		},
	},
}

type Config struct {
	Server configure.Configure
}
