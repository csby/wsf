package client

import (
	"encoding/xml"
	"testing"
)

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
	t.Log("\n", string(data))
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

func TestSoapResult_FromXml(t *testing.T) {
	input := `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope>
  <soap:Body>
    <MethodResponse>
      <MethodResult>
        <ID>1112</ID>
        <Name>Ktp</Name>
        <Remark>remark</Remark>
      </MethodResult>
    </MethodResponse>
    <soap:Fault>
      <soap:Code>
        <soap:Value>soap:Sender</soap:Value>
      </soap:Code>
      <soap:Reason>
        <soap:Text xml:lang="en">Unable to handle request without a valid action parameter.</soap:Text>
      </soap:Reason>
      <soap:Detail />
    </soap:Fault>
  </soap:Body>
 </soap:Envelope>`

	soap := struct {
		SoapResult
		Body struct {
			//XMLName struct{} `xml:"Body"`
			Fault SoapFault

			Response struct {
				XMLName struct{}     `xml:"MethodResponse"`
				Result  methodResult `xml:"MethodResult"`
			}
		}
	}{}

	err := xml.Unmarshal([]byte(input), &soap)
	if err != nil {
		t.Fatal(err)
	}

	raw, err := xml.MarshalIndent(soap, "", "	")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(raw))

	if soap.Body.Response.Result.ID != "1112" {
		t.Error("invalid ID:", soap.Body.Response.Result.ID)
	}
	if soap.Body.Response.Result.Name != "Ktp" {
		t.Error("invalid Name:", soap.Body.Response.Result.Name)
	}

	if soap.Body.Fault.Code.Value != "soap:Sender" {
		t.Error("invalid fault code:", soap.Body.Fault.Code.Value)
	}
	if soap.Body.Fault.Reason.Text != "Unable to handle request without a valid action parameter." {
		t.Error("invalid fault reason:", soap.Body.Fault.Reason.Text)
	}
}

type methodResult struct {
	ID   string `xml:"ID"`
	Name string `xml:"Name"`
}
