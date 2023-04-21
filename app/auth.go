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

	"github.com/metatron-code/metatron-agent/internal/intapi"
	"github.com/metatron-code/metatron-agent/tools"
)

type AuthConfig struct {
	Endpoint  string `json:"endpoint"`
	ThingName string `json:"thing_name"`

	CertificateID                string `json:"certificate_id"`
	CertificateDevice            string `json:"certificate_device"`
	CertificateKeypairPublicKey  string `json:"certificate_keypair_public_key"`
	CertificateKeypairPrivateKey string `json:"certificate_keypair_private_key"`
}

func (app *App) loadAuthConfig() (*AuthConfig, error) {
	confFile := path.Join(app.rootFilePath, "auth.dat")

start:
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

			dataEncrypted, err := tools.EncryptBytes(data, app.config.AgentUUID.String())
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

	data, err := tools.DecryptBytes(dataEncrypted, app.config.AgentUUID.String())
	if err != nil {
		if err := os.Remove(confFile); err != nil {
			return nil, err
		}

		goto start
	}

	conf := &AuthConfig{}
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	if verified, err := app.verifyAuthConfig(conf); err != nil {
		return nil, fmt.Errorf("error verify config: %s", err.Error())

	} else if !verified {
		if err := os.Remove(confFile); err != nil {
			return nil, err
		}

		goto start
	}

	return conf, nil
}

func (app *App) requestAuthConfig() (*AuthConfig, error) {
	client := intapi.NewHttpClient(app.metaVersion, app.metaCommit, app.metaSignKey)

	endpoint := &url.URL{
		Scheme: "https",
		Host:   app.cvmAddress,
		Path:   fmt.Sprintf("/Prod/registration/%s", app.config.AgentUUID.String()),
	}

	resp, err := client.Get(endpoint.String())
	if err != nil {
		return nil, err
	}

	conf := &AuthConfig{}

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

func (app *App) verifyAuthConfig(conf *AuthConfig) (bool, error) {
	client := intapi.NewHttpClient(app.metaVersion, app.metaCommit, app.metaSignKey)

	endpoint := &url.URL{
		Scheme: "https",
		Host:   app.cvmAddress,
		Path:   fmt.Sprintf("/Prod/info/%s", app.config.AgentUUID.String()),
	}

	resp, err := client.Get(endpoint.String())
	if err != nil {
		return false, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var data map[string]string

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}

		if err := json.Unmarshal(body, &data); err != nil {
			return false, err
		}

		if certID, ok := data["certificate_id"]; ok {
			if conf.CertificateID == certID {
				return true, nil
			}
		}

	case http.StatusForbidden:
		return false, fmt.Errorf("agent was blocked")

	case http.StatusNotFound:
		return false, nil
	}

	return false, nil
}
