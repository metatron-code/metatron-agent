package app

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path"

	"github.com/google/uuid"
)

type config struct {
	AgentUUID uuid.UUID `json:"agent_uuid"`
}

func (app *App) loadBaseConfig() error {
	confFile := path.Join(app.rootFilePath, "conf.dat")

	var conf config

	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		conf.AgentUUID = uuid.New()

		data, err := json.Marshal(conf)
		if err != nil {
			return err
		}

		dataEncrypted, err := encryptBytes(data, defaultEncryptPassword)
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

		data, err := decryptBytes(dataEncrypted, defaultEncryptPassword)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &conf); err != nil {
			return err
		}
	}

	app.config = conf

	return nil
}
