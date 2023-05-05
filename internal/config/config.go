package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
)

// ServerConfig server configuration structure
type ServerConfig struct {
	Addr    string `env:"ADDRESS" json:"address"`
	DB      string `env:"DB" json:"db_address"`
	CfgFile string `env:"CFG_FILE"`
}

func NewServerConfig() (*ServerConfig, error) {

	var cfg ServerConfig

	flag.StringVar(&cfg.Addr,
		"a",
		"127.0.0.1",
		"the address where the server is running",
	)
	flag.StringVar(&cfg.Addr,
		"c",
		"",
		"configuration file",
	)
	flag.StringVar(&cfg.DB,
		"d",
		"",
		"the database where the server is running",
	)

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.CfgFile != "" {
		f, err := os.ReadFile(cfg.CfgFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(f, &cfg)
		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}
