package tasks

import (
	"encoding/json"

	probing "github.com/prometheus-community/pro-bing"
)

type IcmpPing struct {
	Target string `json:"target"`
	Count  int    `json:"count"`

	ping *probing.Pinger
}

type IcmpPingResponse struct {
	PacketsSent           int       `json:"packets_sent"`
	PacketsRecv           int       `json:"packets_recv"`
	PacketsRecvDuplicates int       `json:"packets_recv_duplicates"`
	PacketLoss            float64   `json:"packet_loss"`
	IPAddr                string    `json:"ip_address"`
	Addr                  string    `json:"address"`
	MinRtt                float64   `json:"min_rtt"`
	MaxRtt                float64   `json:"maxrtt"`
	AvgRtt                float64   `json:"avg_rtt"`
	StdDevRtt             float64   `json:"std_dev_rtt"`
	Rtts                  []float64 `json:"rtts"`
}

func NewIcmpPing(params []byte) (*IcmpPing, error) {
	task := &IcmpPing{}

	if err := json.Unmarshal(params, &task); err != nil {
		return nil, err
	}

	var err error
	task.ping, err = probing.NewPinger(task.Target)
	if err != nil {
		return nil, err
	}

	task.ping.Count = 4
	if task.Count > 0 {
		task.ping.Count = task.Count
	}

	return task, nil
}

func (t *IcmpPing) Run() ([]byte, error) {
	if err := t.ping.Run(); err != nil {
		return nil, err
	}

	stats := t.ping.Statistics()

	resp := IcmpPingResponse{
		PacketsSent:           stats.PacketsSent,
		PacketsRecv:           stats.PacketsRecv,
		PacketsRecvDuplicates: stats.PacketsRecvDuplicates,
		PacketLoss:            stats.PacketLoss,
		IPAddr:                stats.IPAddr.String(),
		Addr:                  stats.Addr,

		MinRtt: stats.MinRtt.Seconds(),
		MaxRtt: stats.MaxRtt.Seconds(),
		AvgRtt: stats.AvgRtt.Seconds(),
	}

	for _, row := range stats.Rtts {
		resp.Rtts = append(resp.Rtts, row.Seconds())
	}

	respByte, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return respByte, nil
}
