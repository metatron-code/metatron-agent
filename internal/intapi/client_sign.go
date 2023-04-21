package intapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (c *HTTPClient) GetAuthRequestSign(method, path string, body []byte) (string, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	nonce := sha256.Sum256([]byte(fmt.Sprintf("%s/%s/%s", timestamp, c.appVersion, c.appCommit)))

	var b bytes.Buffer
	b.Write(nonce[:])
	b.WriteString("\n")
	b.WriteString(method)
	b.WriteString("\n")
	b.WriteString(path)

	if method != http.MethodGet && method != http.MethodHead {
		bodyHash := sha256.Sum256(body)
		b.WriteString("\n")
		b.Write(bodyHash[:])
	}

	hash := hmac.New(sha256.New, []byte(c.signKey))

	if _, err := hash.Write(b.Bytes()); err != nil {
		return "", err
	}

	data := map[string]string{
		"app":       "metatron-agent",
		"version":   c.appVersion,
		"timestamp": timestamp,
		"signature": base64.RawURLEncoding.EncodeToString(hash.Sum(nil)),
	}

	dataSlice := make([]string, 0)

	for key, val := range data {
		dataSlice = append(dataSlice, fmt.Sprintf("%s=%s", key, val))
	}

	return strings.Join(dataSlice, ";"), nil
}
