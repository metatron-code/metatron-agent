package ipinfo

import (
	"encoding/json"
	"time"

	"github.com/metatron-code/metatron-agent/internal/exthttp"
)

type IPInfo struct {
	City     string `json:"city"`
	Country  string `json:"country"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Region   string `json:"region"`
	Timezone string `json:"timezone"`
}

func GetIPInfo() (*Info, error) {
	client := exthttp.NewHTTPClient()

	resp, err := client.Get("https://ipinfo.io/json")
	if err != nil {
		return nil, err
	}

	var info IPInfo

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &Info{
		IP: info.IP,

		Country: info.Country,
		Region:  info.Region,
		City:    info.City,

		CheckedON: time.Now().Unix(),
	}, nil
}
