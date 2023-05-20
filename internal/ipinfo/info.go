package ipinfo

type Info struct {
	IP string `json:"ip"`

	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`

	CheckedON int64 `json:"checked_on"`
}
