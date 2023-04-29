package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/eclipse/paho.golang/paho"
	"github.com/metatron-code/metatron-agent/internal/vars"
)

type state struct {
	AgentID   string `json:"agent_id"`
	Connected int64  `json:"connected"`
	Uptime    int64  `json:"uptime"`
	Version   string `json:"app_version"`
}

func (app *App) mqttSendState() {
	willMessageExpiry := uint32(5 * 60)

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

		time.Sleep(3 * time.Minute)
	}
}
