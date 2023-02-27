package soapic

import (
	"net/http"
)

// Client struct
// Client is a SOAP client
type Client struct {
	httpClient *http.Client
	url        string
	trace      string
}

// New (httpClient http.Client, url string) Client
func New(httpClient *http.Client, url string, trace string) Client {
	return Client{
		httpClient: httpClient,
		url:        url,
		trace:      trace,
	}
}
