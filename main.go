package main

import (
	"log"
	"os"

	"github.com/metatron-code/metatron-agent/app"
)

var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	signKey   = "none"
	sentryDsn = ""

	defaultEncryptPassword = "qwerty"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	// defer func() {
	// 	err := recover()

	// 	if err != nil {
	// 		log.Println("global error:", err)

	// 		sentry.CurrentHub().Recover(err)
	// 		sentry.Flush(2 * time.Second)
	// 	}
	// }()

	agent, err := app.New(version, commit, date, signKey, sentryDsn)
	if err != nil {
		log.Println("error app initialization:", err)
		return
	}

	agent.SetDefaultEncryptPassword(defaultEncryptPassword)

	if err := agent.Execute(); err != nil {
		log.Println("error execute app:", err)
	}
}
