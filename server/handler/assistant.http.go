package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/csby/wsf/types"
	"io/ioutil"
	"net/http"
	"time"
)

type httpAssistant struct {
	assistant

	response   http.ResponseWriter
	request    *http.Request
	outputCode *int

	method string
	schema string
	path   string
	token  string
	param  []byte
}

func (s *httpAssistant) GetBody() ([]byte, error) {
	return ioutil.ReadAll(s.request.Body)
}

func (s *httpAssistant) GetJson(v interface{}) error {
	err := json.NewDecoder(s.request.Body).Decode(v)
	if err == nil {
		s.input, err = json.Marshal(v)
	}

	return err
}

func (s *httpAssistant) GetXml(v interface{}) error {
	bodyData, err := ioutil.ReadAll(s.request.Body)
	if err != nil {
		return err
	}
	defer s.request.Body.Close()
	s.input = bodyData

	err = xml.Unmarshal(bodyData, v)

	return err
}

func (s *httpAssistant) OutputJson(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Fprint(s.response, err)
	} else {
		s.response.Header().Add("Access-Control-Allow-Origin", "*")
		s.response.Header().Set("Content-Type", "application/json;charset=utf-8")
		s.response.Write(data)
		s.output = data
	}

}

func (s *httpAssistant) OutputXml(v interface{}) {
	if s.response == nil {
		return
	}

	if v != nil {
		switch v.(type) {
		case []byte:
			s.output = v.([]byte)
		case string:
			s.output = []byte(v.(string))
		default:
			bodyData, err := xml.MarshalIndent(v, "", "	")
			if err != nil {
				fmt.Fprint(s.response, err)
				s.output = []byte(err.Error())
				return
			} else {
				s.output = bodyData
			}
		}
	}

	s.response.Header().Add("Access-Control-Allow-Origin", "*")
	s.response.Header().Set("Content-Type", "application/xml;charset=utf-8")
	if len(s.output) > 0 {
		s.response.Write(s.output)
	}
}

func (s *httpAssistant) Success(data interface{}) {
	result := &types.Result{
		Code:   0,
		Data:   data,
		Elapse: time.Now().Sub(s.enterTime).String(),
		Serial: s.rid,
	}
	s.outputCode = &result.Code

	s.OutputJson(result)
}

func (s *httpAssistant) Error(errCode types.ErrorCode, errDetails ...interface{}) {
	result := &types.Result{
		Code:   errCode.Code(),
		Elapse: time.Now().Sub(s.enterTime).String(),
		Serial: s.rid,
		Error: types.ResultError{
			Summary: errCode.Summary(),
			Detail:  fmt.Sprint(errDetails...),
		},
	}
	s.outputCode = &result.Code

	s.OutputJson(result)
}

func (s *httpAssistant) IsError() bool {
	if s.outputCode == nil {
		return false
	}

	if *s.outputCode == 0 {
		return false
	}

	return true
}

func (s *httpAssistant) Method() string {
	return s.method
}

func (s *httpAssistant) Schema() string {
	return s.schema
}

func (s *httpAssistant) Path() string {
	return s.path
}

func (s *httpAssistant) Token() string {
	return s.token
}

func (s *httpAssistant) GetParam() []byte {
	return s.param
}
