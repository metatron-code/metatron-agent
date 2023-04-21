//go:build ignore

package main

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"os"
	"strings"
)

func main() {
	signKeyBytes := make([]byte, 24)
	if _, err := rand.Read(signKeyBytes); err != nil {
		log.Fatal(err)
	}

	signKey := base64.RawURLEncoding.EncodeToString(signKeyBytes)

	var packageTemplate = template.Must(template.New("").Parse(strings.TrimSpace(`
package vars

var (
	SignKey = "{{ .SignKey}}"
)
	`)))

	f, err := os.Create("internal/vars/vars.go")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	packageTemplate.Execute(f, struct {
		SignKey string
	}{
		SignKey: signKey,
	})
}
