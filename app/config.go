package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
			return fmt.Errorf("ERR-CONF-01: %s", err.Error())
		}

		dataEncrypted, err := tools.EncryptBytes(data, app.defaultEncryptPassword)
		if err != nil {
			return fmt.Errorf("ERR-CONF-02: %s", err.Error())
		}

		encoded := base64.RawURLEncoding.EncodeToString(dataEncrypted)

		if err := os.WriteFile(confFile, []byte(encoded), 0600); err != nil {
			if err != nil {
				return fmt.Errorf("ERR-CONF-03: %s", err.Error())
			}
		}

	} else {
		confData, err := os.ReadFile(confFile)
		if err != nil {
			return fmt.Errorf("ERR-CONF-04: %s", err.Error())
		}

		dataEncrypted, err := base64.RawURLEncoding.DecodeString(string(confData))
		if err != nil {
			return fmt.Errorf("ERR-CONF-05: %s", err.Error())
		}

		data, err := tools.DecryptBytes(dataEncrypted, app.defaultEncryptPassword)
		if err != nil {
			if err := os.Remove(confFile); err != nil {
				return fmt.Errorf("ERR-CONF-06: %s", err.Error())
			}

			goto start
		}

		if err := json.Unmarshal(data, &conf); err != nil {
			return fmt.Errorf("ERR-CONF-07: %s", err.Error())
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
				return fmt.Errorf("ERR-CONF-08: %s", err.Error())
			}
		}
	}

	return nil
}
