package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// todo for more xomplex configs
// ConfigFile map[string]*AppConfig
type MegaConfig struct {
}

func (app *MegaConfig) LoadConfig(fname string, cf interface{}) error {
	var fp *os.File
	var err error

	if fp, err = os.Open(fname); err != nil {
		return fmt.Errorf("falied to open : %s", fname)
	}
	defer fp.Close()
	decoder := yaml.NewDecoder(fp)
	if err := decoder.Decode(cf); err != nil {
		return fmt.Errorf("falied to decode : %s ", fname)

	}
	return nil
}
