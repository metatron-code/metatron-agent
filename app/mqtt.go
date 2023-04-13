package app

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (app *App) newMQTTClient() mqtt.Client {
	cert, err := tls.X509KeyPair(
		[]byte(app.mqttAuthConf["certificate_device"]),
		[]byte(app.mqttAuthConf["certificate_keypair_private_key"]),
	)
	if err != nil {
		log.Println("error parse tls key:", err)

	}

	copts := mqtt.NewClientOptions()
	copts.SetClientID(app.mqttAuthConf["thing_name"])
	copts.SetAutoReconnect(true)
	copts.SetMaxReconnectInterval(30 * time.Second)
	copts.SetOnConnectHandler(app.mqttOnConnect)
	copts.SetTLSConfig(&tls.Config{
		Certificates: []tls.Certificate{cert},
	})

	copts.AddBroker(fmt.Sprintf("tcps://%s:8883/mqtt", app.mqttAuthConf["endpoint"]))

	return mqtt.NewClient(copts)
}

func (app *App) mqttOnConnect(client mqtt.Client) {
	taskTopic := fmt.Sprintf("metatron-agent/%s/tasks", app.mqttAuthConf["thing_name"])
	if token := client.Subscribe(taskTopic, 0, app.mqttEventTask); token.Wait() && token.Error() != nil {
		log.Println("error subscribe:", token.Error())
	}
}
