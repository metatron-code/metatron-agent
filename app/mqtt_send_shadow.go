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

	localShadow := shadowReported{
		Version: app.metaVersion,
		EnvOS:   runtime.GOOS,
		EnvArch: runtime.GOARCH,
	}

	if localShadow != remote.State.Reported {
		sendShadow := shadowReported{}

		if remote.State.Reported.Version != localShadow.Version {
			sendShadow.Version = localShadow.Version
		}

		if remote.State.Reported.EnvOS != localShadow.EnvOS {
			sendShadow.EnvOS = localShadow.EnvOS
		}

		if remote.State.Reported.EnvArch != localShadow.EnvArch {
			sendShadow.EnvArch = localShadow.EnvArch
		}

		sendData := shadowData{}
		sendData.State.Reported = sendShadow

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
	sendData := shadowData{}
	sendData.State.Reported = shadowReported{
		Version: app.metaVersion,
		EnvOS:   runtime.GOOS,
		EnvArch: runtime.GOARCH,
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

	app.shadowUpdated = true
}
