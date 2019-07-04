package client

import (
	"encoding/xml"
	"testing"
)

func TestHttp_PostJson(t *testing.T) {
	url := "http://localhost:8080/doc.api/test/api"
	argument := &inputArgument{
		ID:   11,
		Name: "Json",
	}
	client := &Http{}
	input, output, _, _, err := client.PostJson(url, argument)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("input: ", string(input[:]))
	t.Log("output:", string(output[:]))
}

func TestHttp_PostXml(t *testing.T) {
	url := "http://localhost:8080/doc.api/test/api"
	argument := &inputArgument{
		ID:   11,
		Name: "Xml",
	}
	client := &Http{}
	input, output, _, _, err := client.PostXml(url, argument)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("input: ", string(input[:]))
	t.Log("output:", string(output[:]))

	argument.Name = "Xml-[]byte"
	argumentData, err := xml.Marshal(argument)
	if err != nil {
		t.Fatal(err)
	}
	input, output, _, _, err = client.PostXml(url, argumentData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("input: ", string(input[:]))
	t.Log("output:", string(output[:]))

	argument.Name = "Xml-string"
	argumentData, err = xml.Marshal(argument)
	if err != nil {
		t.Fatal(err)
	}
	input, output, _, _, err = client.PostXml(url, string(argumentData))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("input: ", string(input[:]))
	t.Log("output:", string(output[:]))
}

func TestHttp_PostSoap(t *testing.T) {
	argument := &identity{
		bodyData: bodyData{
			XmlNs: "http://tempuri.org/",
		},
		InputID: `{"username":"test1","password":"1"}`,
	}
	url := "http://localhost:8001/service.asmx"

	client := &Http{}
	input, output, _, _, err := client.PostSoap(url, argument)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("soap input: ", string(input[:]))
	t.Log("soap output:", string(output[:]))
}

type inputArgument struct {
	ID      uint64   `json:"id" xml:"id"`
	Name    string   `json:"name" xml:"name"`
	XMLName struct{} `json:"-" xml:"argument"`
}
