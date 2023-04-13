package app

import (
	"crypto/sha256"
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

			encoded := base64.RawURLEncoding.EncodeToString(dataEncrypted)

			if err := os.WriteFile(confFile, []byte(encoded), 0600); err != nil {
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

	dataEncrypted, err := base64.RawURLEncoding.DecodeString(string(confData))
	if err != nil {
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
	sign, err := app.getAuthRequestSign()
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("version", app.metaVersion)
	values.Add("commit", app.metaCommit)
	values.Add("sign", sign)

	endpoint := &url.URL{
		Scheme:   "https",
		Host:     app.cvmAddress,
		Path:     fmt.Sprintf("/registration/%s", app.config.AgentUUID.String()),
		RawQuery: values.Encode(),
	}

	resp, err := http.Get(endpoint.String())
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

func (app *App) getAuthRequestSign() (string, error) {
	hash := sha256.New()

	if _, err := hash.Write([]byte(app.metaCommit)); err != nil {
		return "", err
	}

	path, err := os.Executable()
	if err != nil {
		return "", err
	}

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(hash.Sum(nil)), nil
}
