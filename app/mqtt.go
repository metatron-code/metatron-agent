package app

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (app *App) newMQTTClient() mqtt.Client {
	cert, err := tls.X509KeyPair(
		[]byte(app.mqttAuthConf.CertificateDevice),
		[]byte(app.mqttAuthConf.CertificateKeypairPrivateKey),
	)
	if err != nil {
		log.Println("error parse tls key:", err)
	}

	copts := mqtt.NewClientOptions()
	copts.SetClientID(app.mqttAuthConf.ThingName)
	copts.SetAutoReconnect(true)
	copts.SetMaxReconnectInterval(30 * time.Second)
	copts.SetOnConnectHandler(app.mqttOnConnect)

	copts.SetTLSConfig(&tls.Config{
		ClientAuth:   tls.NoClientCert,
		ClientCAs:    nil,
		Certificates: []tls.Certificate{cert},
	})

	copts.AddBroker(fmt.Sprintf("tcps://%s:8883/mqtt", app.mqttAuthConf.Endpoint))

	return mqtt.NewClient(copts)
}

func (app *App) mqttOnConnect(client mqtt.Client) {
	taskTopic := fmt.Sprintf("metatron-agent/%s/tasks", app.mqttAuthConf.ThingName)
	if token := client.Subscribe(taskTopic, 0, app.mqttEventTask); token.Wait() && token.Error() != nil {
		log.Println("error subscribe:", token.Error())
		app.mqttErrors++
	}
}

type state struct {
	Connected int64 `json:"connected"`
}

func (app *App) mqttSendState() {
	stateTopic := fmt.Sprintf("metatron-agent/%s/state", app.mqttAuthConf.ThingName)

	for {
		if app.mqtt == nil || (app.mqtt != nil && !app.mqtt.IsConnected()) {
			time.Sleep(10 * time.Second)
			continue
		}

		data := state{
			Connected: time.Now().Unix(),
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			log.Println("error marshal state:", err)
			continue
		}

		if token := app.mqtt.Publish(stateTopic, 1, false, dataBytes); token.Wait() && token.Error() != nil {
			log.Println("error publish:", err)
			app.mqttErrors++
		}

		time.Sleep(3 * time.Minute)
	}
}
