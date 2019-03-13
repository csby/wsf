package host

import (
	"crypto/tls"
	"fmt"
	"github.com/csby/security/certificate"
	"github.com/csby/wsf/server/configure"
	"github.com/csby/wsf/server/handler"
	"github.com/csby/wsf/types"
	"net/http"
	"sync"
)

func NewHost(log types.Log, cfg *configure.Configure, httpHandler types.HttpHandler, tcpHandler types.TcpHandler) types.Host {
	instance := &host{cfg: cfg, httpHandler: httpHandler, tcpHandler: tcpHandler}
	instance.SetLog(log)

	return instance
}

type host struct {
	types.Base

	cfg         *configure.Configure
	httpHandler types.HttpHandler
	tcpHandler  types.TcpHandler

	httpServer  *http.Server
	httpsServer *http.Server
}

func (s *host) Run() error {
	if s.cfg == nil {
		return fmt.Errorf(s.LogError("invalid configure: nil"))
	}

	wg := &sync.WaitGroup{}

	if s.httpHandler != nil {
		router, err := handler.NewHttpHandler(s.GetLog(), s.httpHandler)
		if err != nil {
			s.LogError("NewHttpHandler error: ", err)
			return err
		}

		// http
		if s.cfg.Http.Enabled {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer s.LogInfo("http server stopped")

				err := s.runHttp(router)
				if err != nil {
					s.LogError("http server error: ", err)
				}

			}()
		}

		// https
		if s.cfg.Https.Enabled {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer s.LogInfo("https server stopped")

				err := s.runHttps(router)
				if err != nil {
					s.LogError("https server error: ", err)
				}
			}()
		}
	}

	if s.tcpHandler != nil {
		// tcp
		if s.cfg.Tcp.Enabled {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer s.LogInfo("tcp server stopped")

				err := s.runTcp()
				if err != nil {
					s.LogError("tcp server error: ", err)
				}
			}()
		}
	}

	wg.Wait()
	return nil
}

func (s *host) Close() (err error) {
	if s.httpServer != nil {
		e := s.httpServer.Close()
		if e != nil {
			err = e
		}
	}

	if s.httpsServer != nil {
		e := s.httpsServer.Close()
		if e != nil {
			err = e
		}
	}

	return
}

func (s *host) runHttp(handler http.Handler) error {
	defer func() {
		if err := recover(); err != nil {
			s.LogError("http server exception: ", err)
		}
	}()

	addr := fmt.Sprintf("%s:%d", s.cfg.Http.Address, s.cfg.Http.Port)
	s.LogInfo("http server running on \"", addr, "\"")

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	err := s.httpServer.ListenAndServe()
	s.httpServer = nil

	return err
}

func (s *host) runHttps(handler http.Handler) error {
	defer func() {
		if err := recover(); err != nil {
			s.LogError("https server exception: ", err)
		}
	}()

	caFilePath := s.cfg.Https.Cert.Ca.File
	s.LogInfo("https server ca file: ", caFilePath)
	pfxFilePath := s.cfg.Https.Cert.Server.File
	s.LogInfo("https server pfx file: ", pfxFilePath)
	pfx := &certificate.CrtPfx{}
	err := pfx.FromFile(pfxFilePath, s.cfg.Https.Cert.Server.Password)
	if err != nil {
		return fmt.Errorf("load pfx file fail: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", s.cfg.Https.Address, s.cfg.Https.Port)
	s.httpsServer = &http.Server{
		Addr:    addr,
		Handler: handler,
		TLSConfig: &tls.Config{
			Certificates: pfx.TlsCertificates(),
			ClientAuth:   tls.NoClientCert,
		},
	}
	if s.cfg.Https.RequestClientCert {
		s.httpsServer.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	if len(caFilePath) > 0 {
		crt := &certificate.Crt{}
		err = crt.FromFile(caFilePath)
		if err != nil {
			return fmt.Errorf("load ca file fail: %v", err)
		}
		s.httpsServer.TLSConfig.ClientCAs = crt.Pool()
	}

	s.LogInfo("https server running on \"", addr, "\"")
	err = s.httpsServer.ListenAndServeTLS("", "")
	s.httpsServer = nil

	return err
}

func (s *host) runTcp() error {
	defer func() {
		if err := recover(); err != nil {
			s.LogError("tcp server exception: ", err)
		}
	}()

	addr := fmt.Sprintf("%s:%d", s.cfg.Tcp.Address, s.cfg.Tcp.Port)
	s.LogInfo("tcp server running on \"", addr, "\"")

	return nil
}
