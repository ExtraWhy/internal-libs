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
	GoogleClientID     string `yaml:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `yaml:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectUrl  string `yaml:"GOOGLE_REDIRECT_URL"`
	AllowedHosts       string `yaml:"ALLOWED_HOSTS"` //todo check for yaml arrays
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

/*
GOOGLE_CLIENT_ID=173962797108-p2fkcc16vpereds09mcflf79k7j1qtef.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=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
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
ALLOWED_HOSTS=http://localhost:3000, https://cryptowin-ten.vercel.app
*/
