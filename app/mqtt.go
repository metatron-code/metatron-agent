package app

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (app *App) newMQTTClient(conf map[string]string) mqtt.Client {
	cert, err := tls.X509KeyPair(
		[]byte(conf["certificate_device"]),
		[]byte(conf["certificate_keypair_private_key"]),
	)
	if err != nil {
		log.Println("error parse tls key:", err)

	}

	copts := mqtt.NewClientOptions()
	copts.SetClientID(conf["thing_name"])
	copts.SetAutoReconnect(true)
	copts.SetMaxReconnectInterval(30 * time.Second)
	copts.SetOnConnectHandler(app.mqttOnConnect)
	copts.SetTLSConfig(&tls.Config{
		Certificates: []tls.Certificate{cert},
	})

	copts.AddBroker(fmt.Sprintf("tcps://%s:8883/mqtt", conf["endpoint"]))

	return mqtt.NewClient(copts)
}

func (app *App) mqttOnConnect(client mqtt.Client) {
	log.Println("mqtt connected")
}
