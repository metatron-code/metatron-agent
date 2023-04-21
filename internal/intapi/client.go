package intapi

import (
	"fmt"
	"net/http"
)

type HTTPClient struct {
	c http.Client

	appVersion string
	appCommit  string
	signKey    string
}

func NewHTTPClient(version, commit, signKey string) *HTTPClient {
	return &HTTPClient{
		appVersion: version,
		appCommit:  commit,
		signKey:    signKey,
	}
}

func (c *HTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	sign, err := c.GetAuthRequestSign(req.Method, req.URL.Path, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, sign)
}

func (c *HTTPClient) Do(req *http.Request, sign string) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("HMAC-SHA256 %s", sign))

	return c.c.Do(req)
}
