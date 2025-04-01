package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ERROR_REQUEST_CONF int = (1 << 0)
	ERROR_USER_CONF    int = (1 << 1)
	ERROR_OAUTH_CONF   int = (1 << 2)
	MAX_ERRORS         int = 0x07
)

// todo for more xomplex configs
// ConfigFile map[string]*AppConfig
type MegaConfig struct {
	Requests RequestService
	Oauth    OAuthProviderConfig
	User     UserService
}

func (app *MegaConfig) LoadConfig(fname string) error, int {
	var fp *os.File
	var err error
	if fp, err = os.Open(fname); err != nil {
		return fmt.Errorf("falied to open : %s", fname)
	}
	defer fp.Close()
	var errcode = 0
	decoder := yaml.NewDecoder(fp)
	if err := decoder.Decode(&app.Requests); err != nil {
		errcode |= ERROR_REQUEST_CONF
	}
	if err := decoder.Decode(&app.Oauth); err != nil {
		errcode |= ERROR_OAUTH_CONF
	}
	if err := decoder.Decode(&app.User); err != nil {
		errcode |= ERROR_USER_CONF
	}
	if errcode == MAX_ERRORS {
		return fmt.Errorf("falied to decode : %s error %d", fname, errcode)
	}
	return nil, errcode
}
