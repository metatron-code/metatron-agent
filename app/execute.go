package app

import (
	"log"
	"time"
)

func (app *App) Execute() error {
	log.Println("Your HW-UID:", app.agentUUID.String())

	for {
		time.Sleep(1 * time.Second)
	}

}
