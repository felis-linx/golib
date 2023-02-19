package soapic

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// Call ( action string, payload interface{}, result *interface{}) (err error)
// Call send SOAP request
func (client Client) Call(action string, payload interface{}, requestID string) (result []byte, err error) {
	data, err := xml.MarshalIndent(payload, "", "  ")
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", client.url, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", action)
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	response, err := client.httpClient.Do(req)
	if err != nil {
		return
	}

	result, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	defer response.Body.Close()

	return
}
