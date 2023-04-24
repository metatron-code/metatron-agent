package intapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/metatron-code/metatron-agent/internal/vars"
)

func (c *HTTPClient) GetAuthRequestSign(method, path string) (string, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	agentID := c.agentID.String()
	nonce := sha256.Sum256([]byte(fmt.Sprintf("%s/%s/%s", timestamp, c.appVersion, c.appCommit)))

	var b bytes.Buffer
	b.Write(nonce[:])

	b.WriteString("\n")
	b.WriteString(strings.ToLower(method))

	b.WriteString("\n")
	b.WriteString(strings.ToLower(path))

	b.WriteString("\n")
	b.WriteString(agentID)

	signKey, err := base64.RawURLEncoding.DecodeString(vars.SignKey)
	if err != nil {
		return "", err
	}

	hash := hmac.New(sha256.New, signKey)

	if _, err := hash.Write(b.Bytes()); err != nil {
		return "", err
	}

	data := map[string]string{
		"app":       "metatron-agent",
		"agent_id":  agentID,
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
