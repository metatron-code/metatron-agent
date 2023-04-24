package intapi

import (
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

func TestGetAuthRequestSign(t *testing.T) {
	agentID := uuid.New()

	client := NewHTTPClient("0.1.2", "87f173b54157ab59626dd7692f4f317612a98a7f", agentID)

	sign, err := client.GetAuthRequestSign(http.MethodGet, "/")
	if err != nil {
		t.Error(err)
	}

	keys := strings.Split(sign, ";")
	if len(keys) != 5 {
		t.Errorf("error count of sign keys - got: %d, wants: 5", len(keys))
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

	requiredKeys := []string{"app", "agent_id", "timestamp", "version", "signature"}

	for _, row := range requiredKeys {
		if !slices.Contains(currentKeys, row) {
			t.Errorf("required key not found: %s", row)
		}
	}

	log.Printf("Auth sign: %s", sign)
}
