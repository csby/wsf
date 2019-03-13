package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Http struct {
	Transport *http.Transport // usually for https request
	Timeout   int64           // timeout in seconds unit, zero meas not timeout
}

func (s *Http) Get(url string, argument interface{}) (output []byte, connState *tls.ConnectionState, statusCode int, err error) {
	client := s.newClient()
	resp, e := client.Get(url)
	if e != nil {
		err = e
		return
	}
	defer resp.Body.Close()

	connState = resp.TLS
	statusCode = resp.StatusCode

	output, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}

func (s *Http) PostJson(url string, argument interface{}) (input, output []byte, connState *tls.ConnectionState, statusCode int, err error) {
	input = nil
	var body io.Reader = nil
	if argument != nil {
		switch argument.(type) {
		case []byte:
			body = bytes.NewBuffer(argument.([]byte))
			input = argument.([]byte)
		default:
			bodyData, e := json.Marshal(argument)
			if e != nil {
				err = e
				return
			}
			body = bytes.NewBuffer([]byte(bodyData))
			input = bodyData
		}
	}

	client := s.newClient()
	resp, e := client.Post(url, "application/json;charset=utf-8", body)
	if e != nil {
		err = e
		return
	}
	defer resp.Body.Close()

	connState = resp.TLS
	statusCode = resp.StatusCode

	output, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}

func (s *Http) PostXml(url string, argument interface{}) (input, output []byte, connState *tls.ConnectionState, statusCode int, err error) {
	input = nil
	var body io.Reader = nil
	if argument != nil {
		switch argument.(type) {
		case []byte:
			body = bytes.NewBuffer(argument.([]byte))
			input = argument.([]byte)
		case string:
			body = bytes.NewBufferString(argument.(string))
			input = []byte(argument.(string))
		default:
			bodyData, e := xml.MarshalIndent(argument, "", "	")
			if e != nil {
				err = e
				return
			}
			body = bytes.NewBuffer([]byte(bodyData))
			input = bodyData
		}
	}

	client := s.newClient()
	resp, e := client.Post(url, "application/xml;charset=utf-8", body)
	if e != nil {
		err = e
		return
	}
	defer resp.Body.Close()

	connState = resp.TLS
	statusCode = resp.StatusCode

	output, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}

func (s *Http) PostSoap(url string, argument interface{}) (input, output []byte, connState *tls.ConnectionState, statusCode int, err error) {
	soap := &Soap{
		Xsi:    "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:    "http://www.w3.org/2001/XMLSchema",
		Soap12: "http://www.w3.org/2003/05/soap-envelope",
		Body: SoapBody{
			Data: argument,
		},
	}

	input, err = xml.MarshalIndent(soap, "", "	")
	if err != nil {
		return
	}
	body := bytes.NewBuffer([]byte(input))

	client := s.newClient()
	resp, e := client.Post(url, "application/soap+xml;charset=utf-8", body)
	if e != nil {
		err = e
		return
	}
	defer resp.Body.Close()

	connState = resp.TLS
	statusCode = resp.StatusCode

	output, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}

func (s *Http) Download(url string, argument interface{}) ([]byte, *tls.ConnectionState, error) {
	client := s.newClient()
	resp, err := client.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return bodyData, resp.TLS, nil
}

func (s *Http) Certificate() *tls.Certificate {
	if s.Transport == nil {
		return nil
	}
	cfg := s.Transport.TLSClientConfig
	if cfg == nil {
		return nil
	}
	if len(cfg.Certificates) == 0 {
		return nil
	}

	return &cfg.Certificates[0]
}

func (s *Http) newClient() *http.Client {
	client := &http.Client{}
	if s.Transport != nil {
		client.Transport = s.Transport
	}
	if s.Timeout > 0 {
		timeout := s.Timeout * time.Second.Nanoseconds()
		client.Timeout = time.Duration(timeout)
	}

	return client
}
