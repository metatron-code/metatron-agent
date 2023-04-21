package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/url"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

func (app *App) newMQTTClient() (*autopaho.ConnectionManager, error) {
	cert, err := tls.X509KeyPair(
		[]byte(app.mqttAuthConf.CertificateDevice),
		[]byte(app.mqttAuthConf.CertificateKeypairPrivateKey),
	)
	if err != nil {
		log.Println("error parse tls key:", err)
	}

	clientConf := autopaho.ClientConfig{
		BrokerUrls: []*url.URL{
			{
				Scheme: "tcps",
				Host:   fmt.Sprintf("%s:8883", app.mqttAuthConf.Endpoint),
				Path:   "/mqtt",
			},
		},
		ClientConfig: paho.ClientConfig{
			ClientID: app.mqttAuthConf.ThingName,
			Router:   paho.NewSingleHandlerRouter(app.mqttRouter),
		},
		KeepAlive: 30,
		TlsCfg: &tls.Config{
			ClientAuth:   tls.NoClientCert,
			ClientCAs:    nil,
			Certificates: []tls.Certificate{cert},
		},
		OnConnectionUp: app.mqttOnConnect,
	}

	return autopaho.NewConnection(context.Background(), clientConf)
}

func (app *App) mqttOnConnect(cm *autopaho.ConnectionManager, _ *paho.Connack) {
	subscriptions := map[string]paho.SubscribeOptions{
		fmt.Sprintf("metatron-agent/%s/tasks", app.mqttAuthConf.ThingName): {QoS: 1},
	}

	if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{Subscriptions: subscriptions}); err != nil {
		log.Println("error subscribe:", err)
		app.mqttErrors++
	}
}
