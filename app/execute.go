package app

import (
	"context"
	"log"
	"time"
)

func (app *App) Execute() error {
	log.Println("Your HW-UID:", app.config.AgentUUID.String())

	var errCount int
	var err error

	for {
		time.Sleep(10 * time.Second)

		app.mqttAuthConf, err = app.loadAuthConfig()
		if err != nil || app.mqttAuthConf == nil {
			if err != nil {
				log.Println("load auth config error:", err)
			}

			if errCount <= 10 {
				errCount++
			}

			time.Sleep(time.Minute * time.Duration(errCount))
			continue
		}

		if app.mqttAuthConf != nil {
			break
		}
	}

	go app.mqttSendState()
	go app.mqttSendShadow()

	for {
		if app.mqttErrors >= 10 {
			app.mqtt.Disconnect(context.Background())
			app.mqtt = nil
			app.mqttErrors = 0
		}

		if app.mqtt == nil {
			app.mqtt, err = app.newMQTTClient()
			if err != nil {
				log.Println("error mqtt client:", err)
				app.mqtt = nil
			}
		}

		if err := app.mqtt.AwaitConnection(context.Background()); err != nil {
			log.Println("error mqtt await connection", err)
		}

		time.Sleep(10 * time.Second)
	}
}
