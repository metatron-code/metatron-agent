package app

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/metatron-code/metatron-agent/internal/tools"
)

type config struct {
	AgentUUID uuid.UUID `json:"agent_uuid"`
}

func (app *App) loadBaseConfig() error {
	confFile := path.Join(app.rootFilePath, "conf_v1.dat")

	var conf config

start:
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		conf.AgentUUID = uuid.New()

		data, err := json.Marshal(conf)
		if err != nil {
			return err
		}

		dataEncrypted, err := tools.EncryptBytes(data, app.defaultEncryptPassword)
		if err != nil {
			return err
		}

		encoded := base64.RawURLEncoding.EncodeToString(dataEncrypted)

		if err := os.WriteFile(confFile, []byte(encoded), 0600); err != nil {
			if err != nil {
				return err
			}
		}

	} else {
		confData, err := os.ReadFile(confFile)
		if err != nil {
			return err
		}

		dataEncrypted, err := base64.RawURLEncoding.DecodeString(string(confData))
		if err != nil {
			return err
		}

		data, err := tools.DecryptBytes(dataEncrypted, app.defaultEncryptPassword)
		if err != nil {
			if err := os.Remove(confFile); err != nil {
				return err
			}

			goto start
		}

		if err := json.Unmarshal(data, &conf); err != nil {
			return err
		}
	}

	app.config = conf

	return nil
}

func (app *App) removeOldConfigFiles() error {
	for _, name := range []string{"hw.uid", "conf.dat", "auth.dat"} {
		confFile := path.Join(app.rootFilePath, name)

		if _, err := os.Stat(confFile); err == nil {
			if err := os.Remove(confFile); err != nil {
				return err
			}
		}
	}

	return nil
}
