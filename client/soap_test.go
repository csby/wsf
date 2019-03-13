package client

import "testing"

func TestSoap_ToXml(t *testing.T) {
	soap := &Soap{
		Xsi:    "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:    "http://www.w3.org/2001/XMLSchema",
		Soap12: "http://www.w3.org/2003/05/soap-envelope",

		Body: SoapBody{
			Data: &identity{
				bodyData: bodyData{
					XmlNs: "http://tempuri.org/",
				},
				InputID: `{"username":"test1","password":"1"}`,
			},
		},
	}

	data, err := soap.ToXml()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

type bodyData struct {
	XmlNs string `xml:"xmlns,attr"`
}

type identity struct {
	XMLName struct{} `json:"-" xml:"Identity"`
	bodyData

	InputID   string `xml:"inputID"`
	InputData string `xml:"inputData"`
}
