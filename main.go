package main

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/metatron-code/metatron-agent/app"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	defer func() {
		err := recover()

		if err != nil {
			sentry.CurrentHub().Recover(err)
			sentry.Flush(2 * time.Second)
		}
	}()

	agent, err := app.New()
	if err != nil {
		log.Println("error app initialization:", err)
		return
	}

	if err := agent.Execute(); err != nil {
		log.Println("error execute app:", err)
	}
}
