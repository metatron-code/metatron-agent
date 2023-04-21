package intapi

import (
	"fmt"
	"net/http"
)

type HttpClient struct {
	c http.Client

	appVersion string
	appCommit  string
	signKey    string
}

func NewHttpClient(version, commit, signKey string) *HttpClient {
	return &HttpClient{
		appVersion: version,
		appCommit:  commit,
		signKey:    signKey,
	}
}

func (c *HttpClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	sign, err := c.getAuthRequestSign(req.Method, req.URL.Path, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, sign)
}

func (c *HttpClient) Do(req *http.Request, sign string) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("HMAC-SHA256 %s", sign))

	return c.c.Do(req)
}
