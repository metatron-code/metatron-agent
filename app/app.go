package app

import (
	"os"
	"runtime"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/getsentry/sentry-go"
)

type App struct {
	config       config
	rootFilePath string
	cvmAddress   string
	mqtt         mqtt.Client
	mqttAuthConf *AuthConfig
	mqttErrors   int

	startTime time.Time

	metaVersion string
	metaCommit  string
	metaDate    string
	metaSignKey string

	defaultEncryptPassword string
}

func New(version, commit, date, signKey, sentryDsn string) (*App, error) {
	app := &App{
		cvmAddress:  "oyw2ltnb9c.execute-api.eu-west-1.amazonaws.com",
		metaVersion: version,
		metaCommit:  commit,
		metaDate:    date,
		metaSignKey: signKey,

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

	sentryClientOptions := sentry.ClientOptions{
		SampleRate:    0.5,
		EnableTracing: false,
	}

	if len(sentryDsn) > 10 {
		sentryClientOptions.Dsn = sentryDsn
	}

	if version, ok := os.LookupEnv("SNAP_VERSION"); ok {
		sentryClientOptions.Release = version
	}

	if err := sentry.Init(sentryClientOptions); err != nil {
		return nil, err
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTags(map[string]string{
			"OS":      runtime.GOOS,
			"ARCH":    runtime.GOARCH,
			"Version": version,
			"Commit":  commit,
		})

		scope.SetUser(sentry.User{
			ID: app.config.AgentUUID.String(),
		})
	})

	return app, nil
}

func (app *App) SetDefaultEncryptPassword(pass string) {
	app.defaultEncryptPassword = pass
}
