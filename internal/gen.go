//go:build ignore

package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type sendReleaseInfo struct {
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
	Commit     string `json:"commit"`
	SignKey    string `json:"sign_key"`
}

func main() {
	signKeyBytes := make([]byte, 24)
	if _, err := rand.Read(signKeyBytes); err != nil {
		log.Fatal(err)
	}

	signKey := base64.RawURLEncoding.EncodeToString(signKeyBytes)

	var packageTemplate = template.Must(template.New("").Parse(strings.TrimSpace(`
package vars

var (
	SignKey = "{{ .SignKey }}"
)
	`)))

	f, err := os.Create("internal/vars/vars.go")
	if err != nil {
		log.Fatal(err)
	}

	packageTemplate.Execute(f, struct {
		SignKey string
	}{
		SignKey: signKey,
	})

	f.Close()

	if os.Getenv("GITHUB_REF_TYPE") != "tag" {
		return
	}

	send := sendReleaseInfo{
		AppName: "metatron-agent",
		Commit:  os.Getenv("GITHUB_SHA"),
		SignKey: signKey,
	}

	send.AppVersion = strings.TrimPrefix(os.Getenv("GITHUB_REF"), "refs/tags/v")

	body, err := json.Marshal(send)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, os.Getenv("RELEASE_NOTIFY_URL"), bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Key %s", os.Getenv("RELEASE_AUTH_KEY")))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		log.Fatal("err put release info")
	}
}
