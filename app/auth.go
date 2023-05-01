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
	"github.com/metatron-code/metatron-agent/internal/tools"
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
	confFile := path.Join(app.rootFilePath, "auth_v1.dat")

start:
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		conf, err := app.requestAuthConfig()
		if err != nil {
			return nil, fmt.Errorf("ERR-AUTH-01: %s", err.Error())
		}

		if conf != nil {
			data, err := json.Marshal(conf)
			if err != nil {
				return nil, fmt.Errorf("ERR-AUTH-02: %s", err.Error())
			}

			dataEncrypted, err := tools.EncryptBytes(data, app.config.AgentUUID.String())
			if err != nil {
				return nil, fmt.Errorf("ERR-AUTH-03: %s", err.Error())
			}

			encoded := base64.RawURLEncoding.EncodeToString(dataEncrypted)

			if err := os.WriteFile(confFile, []byte(encoded), 0600); err != nil {
				if err != nil {
					return nil, fmt.Errorf("ERR-AUTH-04: %s", err.Error())
				}
			}
		}

		return conf, nil
	}

	confData, err := os.ReadFile(confFile)
	if err != nil {
		return nil, fmt.Errorf("ERR-AUTH-05: %s", err.Error())
	}

	dataEncrypted, err := base64.RawURLEncoding.DecodeString(string(confData))
	if err != nil {
		return nil, fmt.Errorf("ERR-AUTH-06: %s", err.Error())
	}

	data, err := tools.DecryptBytes(dataEncrypted, app.config.AgentUUID.String())
	if err != nil {
		if err := os.Remove(confFile); err != nil {
			return nil, fmt.Errorf("ERR-AUTH-07: %s", err.Error())
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
			return nil, fmt.Errorf("ERR-AUTH-08: %s", err.Error())
		}

		goto start
	}

	return conf, nil
}

func (app *App) requestAuthConfig() (*AuthConfig, error) {
	client := intapi.NewHTTPClient(app.config.AgentUUID)

	endpoint := &url.URL{
		Scheme: "https",
		Host:   app.cvmAddress,
		Path:   fmt.Sprintf("/registration/%s", app.config.AgentUUID.String()),
	}

	resp, err := client.Get(endpoint.String())
	if err != nil {
		return nil, fmt.Errorf("ERR-AUTH-09: %s", err.Error())
	}

	conf := &AuthConfig{}

	switch resp.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("ERR-AUTH-10: %s", err.Error())
		}

		if err := json.Unmarshal(body, &conf); err != nil {
			return nil, fmt.Errorf("ERR-AUTH-11: %s", err.Error())
		}

		if conf != nil {
			return conf, nil
		}

	case http.StatusUnauthorized:
		return nil, fmt.Errorf("error connect to int-api: %d", http.StatusUnauthorized)

	case http.StatusInternalServerError:
		return nil, fmt.Errorf("error connect to int-api: %d", http.StatusUnauthorized)
	}

	return nil, nil
}

func (app *App) verifyAuthConfig(conf *AuthConfig) (bool, error) {
	client := intapi.NewHTTPClient(app.config.AgentUUID)

	endpoint := &url.URL{
		Scheme: "https",
		Host:   app.cvmAddress,
		Path:   fmt.Sprintf("/info/%s", app.config.AgentUUID.String()),
	}

	resp, err := client.Get(endpoint.String())
	if err != nil {
		return false, fmt.Errorf("ERR-AUTH-12: %s", err.Error())
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var data map[string]string

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false, fmt.Errorf("ERR-AUTH-13: %s", err.Error())
		}

		if err := json.Unmarshal(body, &data); err != nil {
			return false, fmt.Errorf("ERR-AUTH-14: %s", err.Error())
		}

		if certID, ok := data["certificate_id"]; ok {
			if conf.CertificateID == certID {
				return true, nil
			}
		}

	case http.StatusUnauthorized:
		return false, fmt.Errorf("error connect to int-api")

	case http.StatusForbidden:
		return false, fmt.Errorf("agent was blocked")

	case http.StatusNotFound:
		return false, nil
	}

	return false, nil
}
