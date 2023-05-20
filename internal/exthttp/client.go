package exthttp

import (
	"fmt"
	"net/http"

	"github.com/metatron-code/metatron-agent/internal/vars"
)

type HTTPClient struct {
	c http.Client

	appVersion string
	appCommit  string
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		appVersion: vars.Version,
		appCommit:  vars.Commit,
	}
}

func (c *HTTPClient) SetVersion(version string) {
	c.appVersion = version
}

func (c *HTTPClient) SetCommit(commit string) {
	c.appCommit = commit
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", fmt.Sprintf("Mozilla/5.0 (compatible; MetaTronAgent/%s; +https://metatron.vitalvas.dev)", c.appVersion))

	return c.c.Do(req)
}

func (c *HTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}
