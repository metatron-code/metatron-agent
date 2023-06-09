package intapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/metatron-code/metatron-agent/internal/vars"
)

type HTTPClient struct {
	c http.Client

	agentID uuid.UUID

	appVersion string
	appCommit  string

	signKey string
}

func NewHTTPClient(agentID uuid.UUID) *HTTPClient {
	return &HTTPClient{
		agentID:    agentID,
		appVersion: vars.Version,
		appCommit:  vars.Commit,

		signKey: vars.SignKey,
	}
}

func (c *HTTPClient) SetSignKey(key string) {
	c.signKey = key
}

func (c *HTTPClient) SetVersion(version string) {
	c.appVersion = version
}

func (c *HTTPClient) SetCommit(commit string) {
	c.appCommit = commit
}

func (c *HTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	sign, err := c.GetAuthRequestSign(req.Method, req.URL.Path, time.Now().Unix())
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
