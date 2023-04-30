package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/paho"
	"github.com/metatron-code/metatron-agent/internal/intapi"
	"github.com/metatron-code/metatron-agent/internal/vars"
)

type state struct {
	AgentID   string `json:"agent_id"`
	Connected int64  `json:"connected"`
	Uptime    int64  `json:"uptime"`
	Version   string `json:"app_version"`

	IPInfo map[string]string `json:"ip_info"`
}

func (app *App) mqttSendState() {
	delay := 3
	willMessageExpiry := uint32((delay * 3) * 60)

	for {
		if app.mqtt == nil {
			time.Sleep(10 * time.Second)
			continue
		}

		func() {
			data := state{
				AgentID:   app.config.AgentUUID.String(),
				Connected: time.Now().Unix(),
				Uptime:    int64(time.Since(app.startTime).Seconds()),
				Version:   vars.Version,
			}

			ipInfo, err := app.getIPInfo()
			if err != nil {
				log.Println("error get ip info:", err.Error())
			} else {
				data.IPInfo = ipInfo
			}

			dataBytes, err := json.Marshal(data)
			if err != nil {
				log.Println("error marshal state:", err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			msg := &paho.Publish{
				Topic:   fmt.Sprintf("metatron-agent/%s/state", app.mqttAuthConf.ThingName),
				QoS:     1,
				Payload: dataBytes,
				Properties: &paho.PublishProperties{
					ContentType:   "application/json",
					MessageExpiry: &willMessageExpiry,
				},
			}

			if _, err := app.mqtt.Publish(ctx, msg); err != nil {
				log.Println("error publish:", err)
				app.mqttErrors++
			}
		}()

		time.Sleep(time.Duration(delay) * time.Minute)
	}
}

func (app *App) getIPInfo() (map[string]string, error) {
	client := intapi.NewHTTPClient(app.config.AgentUUID)

	endpoint := &url.URL{Scheme: "https", Host: app.intAPIAddress, Path: "/ipinfo"}

	resp, err := client.Get(endpoint.String())
	if err != nil {
		return nil, err
	}

	var data map[string]string

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code for get IP info")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}
