package main

import (
	"log"
	"os"

	"github.com/metatron-code/metatron-agent/app"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	defaultEncryptPassword = "qwerty"
)

//go:generate go run ./internal/gen.go

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	agent, err := app.New(version, commit, date)
	if err != nil {
		log.Println("error app initialization:", err)
		return
	}

	agent.SetDefaultEncryptPassword(defaultEncryptPassword)

	if err := agent.Execute(); err != nil {
		log.Println("error execute app:", err)
	}
}
