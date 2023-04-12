package app

import (
	"log"
	"time"
)

func (app *App) Execute() error {
	log.Println("Your HW-UID:", app.config.AgentUUID.String())

	var errCount int
	for {
		time.Sleep(10 * time.Second)

		authConf, err := app.loadAuthConfig()
		if err != nil || authConf == nil {
			if err != nil {
				log.Println("load auth config error:", err)
			}

			if errCount <= 10 {
				errCount++
			}

			time.Sleep(time.Minute * time.Duration(errCount))
			continue
		}

		errCount = 0

		if app.mqtt == nil {
			app.mqtt = app.newMQTTClient(authConf)
		}

	}

}
