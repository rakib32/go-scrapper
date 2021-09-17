package http

import (
	"net/http"
	"time"
)

func Get(url string, client *http.Client) (*http.Response, error) {

	return client.Get(url)
}

func CustomClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 20
	transport.MaxConnsPerHost = 20
	transport.MaxIdleConnsPerHost = 20

	client := http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &client
}
