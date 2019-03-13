package example

import (
	"crypto/tls"
	"fmt"
	"github.com/csby/security/certificate"
	"github.com/csby/wsf/client"
	"github.com/csby/wsf/types"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	wg        sync.WaitGroup
	transport *http.Transport
}

func (s *Client) Run() <-chan error {
	// load certificate for client
	clientCrt := &certificate.CrtPfx{}
	err := clientCrt.FromFile(clientCrtFilePath, clientCrtPassword)
	if err != nil {
		log.Error("load client certificate fail: ", err)
	} else {
		s.transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: clientCrt.TlsCertificates(),
			},
		}

		caCrt := &certificate.Crt{}
		err = caCrt.FromFile(caCrtFilePath)
		if err != nil {
			log.Error("load ca certificate fail: ", err)
		} else {
			s.transport.TLSClientConfig.RootCAs = caCrt.Pool()
		}
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runTestApi()
	}()

	ch := make(chan error)
	go func(err chan<- error) {
		time.Sleep(time.Second)
		s.wg.Wait()
		err <- nil
	}(ch)

	return ch
}

func (s *Client) runTestApi() {
	httpClient := client.Http{}
	url := s.getTestApiUrl("http", uriHello)
	_, output, connState, _, err := httpClient.PostJson(url, nil)
	if err != nil {
		log.Error("test api hello fail: ", err)
	} else {
		if connState != nil {
			log.Error("connState: ", connState)
		}
		result := &types.Result{}
		err = result.Unmarshal(output)
		if err != nil {
			log.Error("test api hello fail: ", err)
		} else {
			log.Info("test api hello success")
			fmt.Println(result.FormatString())
			if result.Code != 0 {
				log.Error("test api hello fail: error code = ", result.Code)
			}
		}
	}

	url = s.getTestApiUrl("https", uriHello)
	_, output, connState, _, err = httpClient.PostJson(url, nil)
	if err == nil {
		log.Error("test api hello fail: ", "no transport should be err")
	} else {
		log.Debug("no transport error info: ", err)
	}

	httpClient.Transport = s.transport
	_, output, connState, _, err = httpClient.PostJson(url, nil)
	if err != nil {
		log.Error("test api hello fail: ", err)
	} else {
		if connState == nil {
			log.Error("connState: nil")
		} else {
			serverCrt := &certificate.Crt{}
			serverCrt.FromConnectionState(connState)
			sou := serverCrt.OrganizationalUnit()
			if sou != serverCrtOU {
				log.Error("test api hello server organization unit invalid: ", sou)
			}
		}
		result := &types.Result{}
		err = result.Unmarshal(output)
		if err != nil {
			log.Error("test api hello fail: ", err)
		} else {
			log.Info("test api hello success")
			fmt.Println(result.FormatString())
			if result.Code != 0 {
				log.Error("test api hello fail: error code = ", result.Code)
			}
		}
	}
}

func (s *Client) getTestApiUrl(schema, uri string) string {
	port := cfg.Server.Http.Port
	if schema == "https" {
		port = cfg.Server.Https.Port
	}
	return fmt.Sprintf("%s://localhost:%d%s%s", schema, port, testPath.Prefix, uri)
}
