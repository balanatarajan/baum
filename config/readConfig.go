package config

import "github.com/balanatarajan/baum/config/cdb"
import "fmt"

type Config struct {
	is_cdb bool     // Boolean that indicates whether CDB was used as ac onfiguration mechanism
	cdb    *cdb.Cdb //CDB file that contains the config details
}

func ReadConfig(file string) (*Config, error) {

	c, err := cdb.Open(file)

	if err != nil {
		return nil, err
	}

	var cfg Config
	cfg.cdb = c
	cfg.is_cdb = true

	return &cfg, nil
}

func (c *Config) Find(s string) (string, error) {

	v, err := c.cdb.Data([]byte(s))

	if err != nil {
		fmt.Println("Unable to find key", s)
		return "", err
	}

	return string(v), nil
}

func (c *Config) Close() error {

	err := c.cdb.Close()

	c.cdb = nil
	c.is_cdb = false

	return err
}
