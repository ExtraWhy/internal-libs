package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// todo for more xomplex configs
// ConfigFile map[string]*AppConfig
type RequestService struct {
	DatabaseUrl     string `yaml:"database_url"`
	DatabaseType    string `yaml:"database_type"`
	RestServiceHost string `yaml:"rest_service_host"`
	RestServicePort string `yaml:"rest_service_port"`
}

func (app *RequestService) LoadConfig(fname string) error {
	var fp *os.File
	var err error
	if fp, err = os.Open(fname); err != nil {
		return fmt.Errorf("falied to open : %s", fname)
	}
	defer fp.Close()

	decoder := yaml.NewDecoder(fp)
	if err := decoder.Decode(app); err != nil {
		return fmt.Errorf("falied to decode : %s", fname)
	}

	return nil

}

type OAuthProviderConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectUrl  string `yaml:"redirect_url"`
}

type UserService struct {
	AllowedHosts     string              `yaml:"allowed_hosts"` //todo check for yaml arrays
	GoogleProvider   OAuthProviderConfig `yaml:"google_provider"`
	FacebookProvider OAuthProviderConfig `yaml:"facebook_provider"`
}

func (app *UserService) LoadConfig(fname string) error {
	var fp *os.File
	var err error
	if fp, err = os.Open(fname); err != nil {
		return fmt.Errorf("falied to open : %s", fname)
	}
	defer fp.Close()

	decoder := yaml.NewDecoder(fp)
	if err := decoder.Decode(app); err != nil {
		return fmt.Errorf("falied to decode : %s", fname)
	}

	return nil

}
