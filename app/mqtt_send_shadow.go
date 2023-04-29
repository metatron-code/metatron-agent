package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/eclipse/paho.golang/paho"
	"github.com/metatron-code/metatron-agent/internal/vars"
)

type shadowData struct {
	State shadowState `json:"state"`
}

type shadowState struct {
	Reported shadowReported `json:"reported"`
}

type shadowReported struct {
	Version string `json:"version,omitempty"`
	EnvOS   string `json:"env_os,omitempty"`
	EnvArch string `json:"env_arch,omitempty"`
}

func (app *App) mqttSendShadow() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := &paho.Publish{
		Topic:   fmt.Sprintf("$aws/things/%s/shadow/get", app.mqttAuthConf.ThingName),
		Payload: []byte("{}"),
	}

	if _, err := app.mqtt.Publish(ctx, msg); err != nil {
		log.Println("error publish:", err)
		app.mqttErrors++
	}

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)

		if app.shadowUpdated {
			return
		}
	}

	if !app.shadowUpdated {
		app.mqttForceSendShadow()
	}
}

func (app *App) mqttEventShadow(msg *paho.Publish) {
	if msg.Payload == nil {
		return
	}

	var remote shadowData

	if err := json.Unmarshal(msg.Payload, &remote); err != nil {
		log.Println("error unmarshal remove shadow state:", err)
		return
	}

	localShadow := getShadowData()

	if localShadow.State.Reported != remote.State.Reported {
		sendShadow := shadowReported{}

		if remote.State.Reported.Version != localShadow.State.Reported.Version {
			sendShadow.Version = localShadow.State.Reported.Version
		}

		if remote.State.Reported.EnvOS != localShadow.State.Reported.EnvOS {
			sendShadow.EnvOS = localShadow.State.Reported.EnvOS
		}

		if remote.State.Reported.EnvArch != localShadow.State.Reported.EnvArch {
			sendShadow.EnvArch = localShadow.State.Reported.EnvArch
		}

		sendData := shadowData{
			State: shadowState{
				Reported: sendShadow,
			},
		}

		sendDataBytes, err := json.Marshal(sendData)
		if err != nil {
			log.Println("error marshal shadow data:", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		msg := &paho.Publish{
			Topic:   fmt.Sprintf("$aws/things/%s/shadow/update", app.mqttAuthConf.ThingName),
			Payload: sendDataBytes,
		}

		if _, err := app.mqtt.Publish(ctx, msg); err != nil {
			log.Println("error publish:", err)
			app.mqttErrors++
		}
	}

	app.shadowUpdated = true
}

func (app *App) mqttForceSendShadow() {
	sendData := getShadowData()

	sendDataBytes, err := json.Marshal(sendData)
	if err != nil {
		log.Println("error marshal shadow data:", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := &paho.Publish{
		Topic:   fmt.Sprintf("$aws/things/%s/shadow/update", app.mqttAuthConf.ThingName),
		Payload: sendDataBytes,
	}

	if _, err := app.mqtt.Publish(ctx, msg); err != nil {
		log.Println("error publish:", err)
		app.mqttErrors++
	}

	app.shadowUpdated = true
}

func getShadowData() shadowData {
	return shadowData{
		State: shadowState{
			Reported: shadowReported{
				Version: vars.Version,
				EnvOS:   runtime.GOOS,
				EnvArch: runtime.GOARCH,
			},
		},
	}
}
