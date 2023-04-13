package app

import (
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

	for {
		if app.mqtt == nil {
			app.mqtt = app.newMQTTClient()
		}

		if !app.mqtt.IsConnected() {
			if token := app.mqtt.Connect(); token.Wait() && token.Error() != nil {
				return token.Error()
			}
		}

		time.Sleep(10 * time.Second)
	}
}
