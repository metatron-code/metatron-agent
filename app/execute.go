package app

import "log"

func (app *App) Execute() error {
	log.Println("Your HW-UID:", app.agentUUID.String())

	return nil
}
