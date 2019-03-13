package client

import "encoding/xml"

type Soap struct {
	XMLName struct{} `json:"-" xml:"soap12:Envelope"`
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
