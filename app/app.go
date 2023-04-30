package app

import (
	"os"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/metatron-code/metatron-agent/internal/vars"
)

type App struct {
	config       config
	rootFilePath string

	cvmAddress    string
	intAPIAddress string

	mqtt         *autopaho.ConnectionManager
	mqttAuthConf *AuthConfig
	mqttErrors   int

	startTime time.Time

	defaultEncryptPassword string

	shadowUpdated bool
}

func New() (*App, error) {
	app := &App{
		cvmAddress:    "4w8lflsa93.execute-api.eu-west-1.amazonaws.com",
		intAPIAddress: "z2xzt7xf18.execute-api.eu-west-1.amazonaws.com",

		startTime: time.Now(),

		defaultEncryptPassword: vars.DefaultEncryptPassword,
	}

	var okRootPath bool
	app.rootFilePath, okRootPath = os.LookupEnv("SNAP_COMMON")
	if !okRootPath {
		return nil, errorRunSnapd
	}

	if _, ok := os.LookupEnv("SNAP_COOKIE"); !ok {
		return nil, errorRunSnapd
	}

	if err := app.loadBaseConfig(); err != nil {
		return nil, err
	}

	if err := app.removeOldConfigFiles(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) SetDefaultEncryptPassword(pass string) {
	app.defaultEncryptPassword = pass
}
