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

type LoginService struct {
	ClientID     string `yaml:"client_id"`
	ClientSec    string `yaml:"client_sec"`
	RedirectUrl  string `yaml:"redir_url"`
	AllowedHosts string `yaml:"allowed_hosts"` //todo check for yaml arrays
}

func (app *LoginService) LoadConfig(fname string) error {
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
