package tasks

import (
	"encoding/json"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type IcmpPing struct {
	Target  string `json:"target"`
	Count   int    `json:"count"`
	Network string `json:"network"`

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
	Jitter                float64   `json:"jitter"`
}

func NewIcmpPing(params []byte) (*IcmpPing, error) {
	task := &IcmpPing{}

	if err := json.Unmarshal(params, &task); err != nil {
		return nil, err
	}

	task.ping = probing.New(task.Target)

	task.ping.Count = 4
	if task.Count > 0 {
		task.ping.Count = task.Count
	}

	// TODO: add custom resolver due golang bug: https://github.com/golang/go/issues/28666
	switch task.Network {
	case "ipv4":
		task.ping.SetNetwork("ip4")

	case "ipv6":
		task.ping.SetNetwork("ip6")
	}

	return task, nil
}

func (t *IcmpPing) Run(timeout time.Duration) ([]byte, error) {
	t.ping.Timeout = timeout

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

	if stats.Rtts != nil {
		var rttsSum float64
		for _, row := range stats.Rtts {
			resp.Rtts = append(resp.Rtts, row.Seconds())
			rttsSum += row.Seconds()
		}

		resp.Jitter = rttsSum / float64(len(stats.Rtts))
	}

	respByte, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return respByte, nil
}
