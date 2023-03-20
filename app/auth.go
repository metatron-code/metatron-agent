package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func (app *App) loadAuthConfig() (map[string]string, error) {
	confFile := path.Join(app.rootFilePath, "auth.dat")

	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		conf, err := app.requestAuthConfig()
		if err != nil {
			return nil, err
		}

		if conf != nil {
			data, err := json.Marshal(conf)
			if err != nil {
				return nil, err
			}

			dataEncrypted, err := encryptBytes(data, app.config.AgentUUID.String())
			if err != nil {
				return nil, err
			}

			var encoded []byte
			base64.RawURLEncoding.Encode(encoded, dataEncrypted)

			if err := os.WriteFile(confFile, encoded, 0600); err != nil {
				if err != nil {
					return nil, err
				}
			}
		}

		return conf, nil
	}

	confData, err := os.ReadFile(confFile)
	if err != nil {
		return nil, err
	}

	var dataEncrypted []byte

	if _, err := base64.RawURLEncoding.Decode(dataEncrypted, confData); err != nil {
		return nil, err
	}

	data, err := decryptBytes(dataEncrypted, app.config.AgentUUID.String())
	if err != nil {
		return nil, err
	}

	var conf map[string]string
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func (app *App) requestAuthConfig() (map[string]string, error) {
	endpoint := fmt.Sprintf("https://%s/registration/%s", app.cvmAddress, app.config.AgentUUID.String())

	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(endpointURL.String())
	if err != nil {
		return nil, err
	}

	var conf map[string]string

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(body, &conf); err != nil {
			return nil, err
		}

		if conf != nil {
			return conf, nil
		}
	}

	return nil, nil
}
