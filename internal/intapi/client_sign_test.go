package intapi

import (
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

func TestGetAuthRequestSign(t *testing.T) {
	agentID, err := uuid.Parse("34538c96-808c-465d-be8b-bba49339392f")
	if err != nil {
		t.Error(err)
	}

	var timestamp int64 = 1682466181

	client := NewHTTPClient(agentID)

	client.SetSignKey("-7sSHnpPQl3gq27jyu8qdl_gtZphGFgc")
	client.SetVersion("0.1.2")
	client.SetCommit("87f173b54157ab59626dd7692f4f317612a98a7f")

	sign, err := client.GetAuthRequestSign(http.MethodGet, "/", timestamp)
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

	if sign != "agent_id=34538c96-808c-465d-be8b-bba49339392f;app=metatron-agent;signature=TqvihO0Cd4sTlG0K6bI_SLFUjMpzEElFhylskSczABo;timestamp=1682466181;version=0.1.2" {
		t.Errorf("error validate sign: %s", sign)
	}
}
