package app

import (
	"fmt"
	"log"

	"github.com/eclipse/paho.golang/paho"
)

func (app *App) mqttRouter(msg *paho.Publish) {
	switch msg.Topic {
	case fmt.Sprintf("metatron-agent/%s/tasks", app.mqttAuthConf.ThingName):
		app.mqttEventTask(msg)

	default:
		log.Println("unknown message from unknown topic:", msg.Topic)
	}
}
