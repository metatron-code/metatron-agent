package intapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
)

func (c *HTTPClient) GetAuthRequestSign(method, path string, timestamp int64) (string, error) {
	timestampStr := fmt.Sprintf("%d", timestamp)
	agentID := c.agentID.String()

	nonceStr := fmt.Sprintf("%s/%s/%s", timestampStr, c.appVersion, c.appCommit)
	nonce := sha256.Sum256([]byte(nonceStr))

	var b bytes.Buffer
	b.Write(nonce[:])

	b.WriteString("\n")
	b.WriteString(strings.ToLower(method))

	b.WriteString("\n")
	b.WriteString(strings.ToLower(path))

	b.WriteString("\n")
	b.WriteString(agentID)

	signKey, err := base64.RawURLEncoding.DecodeString(c.signKey)
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
		"timestamp": timestampStr,
		"signature": base64.RawURLEncoding.EncodeToString(hash.Sum(nil)),
	}

	dataSlice := make([]string, 0)

	for key, val := range data {
		dataSlice = append(dataSlice, fmt.Sprintf("%s=%s", key, val))
	}

	sort.Strings(dataSlice)

	return strings.Join(dataSlice, ";"), nil
}
