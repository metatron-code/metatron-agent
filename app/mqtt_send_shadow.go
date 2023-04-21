package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/eclipse/paho.golang/paho"
)

type shadowData struct {
	State struct {
		Reported shadowReported `json:"reported"`
	} `json:"state"`
}
type shadowReported struct {
	Version string `json:"version"`
	EnvOS   string `json:"env_os"`
	EnvArch string `json:"env_arch"`
}

func (app *App) mqttSendShadow() {
	for {
		if app.mqtt == nil {
			time.Sleep(10 * time.Second)
			continue
		}

		func() {
			data := shadowData{}
			data.State.Reported = shadowReported{
				Version: app.metaVersion,
				EnvOS:   runtime.GOOS,
				EnvArch: runtime.GOARCH,
			}

			dataBytes, err := json.Marshal(data)
			if err != nil {
				log.Println("error marshal shadow:", err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			msg := &paho.Publish{
				Topic:   fmt.Sprintf("$aws/things/%s/shadow/update", app.mqttAuthConf.ThingName),
				Payload: dataBytes,
				Properties: &paho.PublishProperties{
					ContentType: "application/json",
				},
			}

			if _, err := app.mqtt.Publish(ctx, msg); err != nil {
				log.Println("error publish:", err)
				app.mqttErrors++
			}
		}()

		time.Sleep(time.Hour)
	}
}
