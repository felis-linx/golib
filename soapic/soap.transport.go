package soapic

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

var debug = false

// Call ( action string, payload interface{}, result *interface{}) (err error)
// Call send SOAP request
func (client Client) Call(action string, payload interface{}, requestID string) (result []byte, err error) {
	data, err := xml.MarshalIndent(payload, "", "  ")
	if err != nil {
		return
	}

	fmt.Println("Request:\n", string(data))
	// requestLog, _ := getMaskedLog(client.url, action, requestID, payload)
	// noty.Info(string(requestLog), fmt.Sprintf("%s Request: %s", client.trace, action), &requestID)

	if debug == true {
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

	// fmt.Println(string(bodyBytes))
	// err = xml.Unmarshal(responsePayload, &result)
	// if err != nil {
	// 	return
	// }

	fmt.Println("Response:\n", string(result))
	return
}
