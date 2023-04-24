package intapi

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type HTTPClient struct {
	c http.Client

	agentID uuid.UUID

	appVersion string
	appCommit  string
}

func NewHTTPClient(version, commit string, agentID uuid.UUID) *HTTPClient {
	return &HTTPClient{
		agentID:    agentID,
		appVersion: version,
		appCommit:  commit,
	}
}

func (c *HTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	sign, err := c.GetAuthRequestSign(req.Method, req.URL.Path)
	if err != nil {
		return nil, err
	}

	return c.Do(req, sign)
}

func (c *HTTPClient) Do(req *http.Request, sign string) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("HMAC-SHA256 %s", sign))

	req.Header.Set("User-Agent", fmt.Sprintf("Mozilla/5.0 (compatible; MetaTronAgent/%s; +https://metatron.vitalvas.dev)", c.appVersion))

	return c.c.Do(req)
}
