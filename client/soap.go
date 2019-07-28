package client

import "encoding/xml"

type Soap struct {
	XMLName struct{} `xml:"soap12:Envelope"`
	Xsi     string   `xml:"xmlns:xsi,attr"`
	Xsd     string   `xml:"xmlns:xsd,attr"`
	Soap12  string   `xml:"xmlns:soap12,attr"`

	Body SoapBody `xml:"soap12:Body"`
}

func (s *Soap) ToXml() ([]byte, error) {
	data, err := xml.MarshalIndent(s, "", "	")
	if err != nil {
		return nil, err
	}

	return data, nil
}

type SoapBody struct {
	Data interface{}
}

// <?xml version="1.0" encoding="utf-8"?>
// <soap:Envelope>
//   <soap:Body>
//     <MethodResponse>
//       <MethodResult>
//         <ID>1112</ID>
//         <Name>Ktp</Name>
//       </MethodResult>
//     </MethodResponse>
//     <soap:Fault>
//       <soap:Code>
//         <soap:Value>soap:Sender</soap:Value>
//       </soap:Code>
//       <soap:Reason>
//         <soap:Text xml:lang="en">Unable to handle request without a valid action parameter.</soap:Text>
//       </soap:Reason>
//       <soap:Detail />
//     </soap:Fault>
//   </soap:Body>
// </soap:Envelope>
type SoapResult struct {
	XMLName struct{} `xml:"Envelope"`
}

type SoapFault struct {
	XMLName struct{} `xml:"Fault"`

	Code   SoapFaultCode   `xml:"Code"`
	Reason SoapFaultReason `xml:"Reason"`
}

type SoapFaultCode struct {
	Value string `xml:"Value"`
}

type SoapFaultReason struct {
	Text string `xml:"Text"`
}
