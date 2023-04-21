package intapi

import (
	"net/http"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestGetAuthRequestSign(t *testing.T) {
	signKey := "PTCSfWLuuCoF8AHLCTzgY8EOd8Xy79yU"
	client := NewHTTPClient("0.1.2", "87f173b54157ab59626dd7692f4f317612a98a7f", signKey)

	sign, err := client.GetAuthRequestSign(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
	}

	keys := strings.Split(sign, ";")
	if len(keys) != 4 {
		t.Errorf("error count of sign keys - got: %d, wants: 4", len(keys))
	}

	currentKeys := make([]string, len(keys))

	for idx, row := range keys {
		if !strings.Contains(row, "=") {
			t.Errorf("error get key-valu - got: %s", row)
		}

		keyValue := strings.Split(row, "=")
		if len(keyValue) != 2 {
			t.Errorf("error parse key-value: %#v", keyValue)
		}

		currentKeys[idx] = keyValue[0]
	}

	requiredKeys := []string{"app", "timestamp", "version", "signature"}

	for _, row := range requiredKeys {
		if !slices.Contains(currentKeys, row) {
			t.Errorf("required key not found: %s", row)
		}
	}
}
