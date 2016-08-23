package api

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func sendHttpRequest(method, rawurl string, body io.Reader, header map[string]string) (*http.Response, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return &http.Response{}, err
	}

	var client *http.Client
	switch url.Scheme {
	case "":
		fallthrough
	case "https":
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	case "http":
		client = &http.Client{}
	default:
		return &http.Response{}, fmt.Errorf("bad url schema: %v", url.Scheme)
	}

	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return &http.Response{}, err
	}
	req.URL.RawQuery = req.URL.Query().Encode()
	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}
