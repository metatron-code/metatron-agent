package app

import (
	"os"
	"path"
	"runtime"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
)

type App struct {
	agentUUID uuid.UUID
}

func New() (*App, error) {
	app := &App{
		agentUUID: uuid.New(),
	}

	rootPath, okRootPath := os.LookupEnv("SNAP_COMMON")
	if !okRootPath {
		return nil, errorRunSnapd
	}

	if _, ok := os.LookupEnv("SNAP_COOKIE"); !ok {
		return nil, errorRunSnapd
	}

	if err := app.loadAgentID(rootPath); err != nil {
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
			"OS":   runtime.GOOS,
			"ARCH": runtime.GOARCH,
		})

		scope.SetUser(sentry.User{
			ID: app.agentUUID.String(),
		})
	})

	return app, nil
}

func (app *App) loadAgentID(rootPath string) error {
	uidPath := path.Join(rootPath, "hw.uid")

	if _, err := os.Stat(uidPath); os.IsNotExist(err) {
		uidDataNew, err := uuid.NewUUID()
		if err != nil {
			return err
		}

		uidDataText, err := uidDataNew.MarshalText()
		if err != nil {
			return err
		}

		if err := os.WriteFile(uidPath, uidDataText, 0600); err != nil {
			if err != nil {
				return err
			}
		}

	} else {
		uidData, err := os.ReadFile(uidPath)
		if err != nil {
			return err
		}

		if err := app.agentUUID.UnmarshalText(uidData); err != nil {
			return err
		}

	}

	return nil
}
