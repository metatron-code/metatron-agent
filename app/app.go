package app

import (
	"os"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
)

type App struct {
	config       config
	rootFilePath string
	cvmAddress   string
	mqtt         *autopaho.ConnectionManager
	mqttAuthConf *AuthConfig
	mqttErrors   int

	startTime time.Time

	metaVersion string
	metaCommit  string
	metaDate    string

	defaultEncryptPassword string

	shadowUpdated bool
}

func New(version, commit, date string) (*App, error) {
	app := &App{
		cvmAddress:  "4w8lflsa93.execute-api.eu-west-1.amazonaws.com",
		metaVersion: version,
		metaCommit:  commit,
		metaDate:    date,

		startTime: time.Now(),
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
