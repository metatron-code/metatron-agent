package app

import (
	"os"
	"runtime"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/getsentry/sentry-go"
)

type App struct {
	config       config
	rootFilePath string
	cvmAddress   string
	mqtt         mqtt.Client

	metaVersion string
	metaCommit  string
	metaDate    string
	metaSignKey string
}

func New(version, commit, date, signKey string) (*App, error) {
	app := &App{
		cvmAddress:  "cvm-prod.metatron.get-server.net",
		metaVersion: version,
		metaCommit:  commit,
		metaDate:    date,
		metaSignKey: signKey,
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

	sentryClientOptions := sentry.ClientOptions{
		SampleRate:    0.5,
		EnableTracing: false,
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
