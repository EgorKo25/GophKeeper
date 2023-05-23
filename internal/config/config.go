package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
)

// ServerConfig server configuration structure
type ServerConfig struct {
	Addr         string `env:"ADDRESS" json:"address"`
	DB           string `env:"DB" json:"db_address"`
	AccessToken  string `env:"ACCESS" json:"access_token"`
	RefreshToken string `env:"REFRESH" json:"refresh_token"`
	PublicPath   string `env:"PUBLIC_PATH" json:"public_path"`
	CfgFile      string `env:"CFG_FILE"`
}

// NewServerConfig server config constructor
func NewServerConfig() (*ServerConfig, error) {

	var cfg ServerConfig

	flag.StringVar(&cfg.Addr,
		"a",
		"127.0.0.1:8080",
		"the address where the server is running",
	)
	flag.StringVar(&cfg.RefreshToken,
		"tr",
		"your-access-secret-key",
		"secret key for jwt access token",
	)
	flag.StringVar(&cfg.AccessToken,
		"ta",
		"your-refresh-secret-key",
		"secret key for jwt refresh token",
	)
	flag.StringVar(&cfg.DB,
		"d",
		"postgresql://postgres:jxrfhbr25@localhost:5432/gophkeeper",
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

// AgentConfig agent configuration structure
type AgentConfig struct {
	AddrServ string `env:"ADDRESS" json:"address"`
	CfgFile  string `env:"CFG_FILE"`
	Secret   string `env:"SECRET" json:"secret"`
}

// NewAgentConfig is a constructor client configuration
func NewAgentConfig() (*AgentConfig, error) {

	ac := AgentConfig{}
	flag.StringVar(&ac.Secret,
		"secret",
		"some-sec",
		"secret key",
	)
	flag.StringVar(&ac.AddrServ,
		"address",
		"http://127.0.0.1:8080",
		"address server",
	)
	flag.StringVar(&ac.CfgFile,
		"config",
		"",
		"configuration file",
	)

	flag.Parse()

	err := env.Parse(&ac)
	if err != nil {
		return nil, err
	}

	if ac.CfgFile != "" {
		f, err := os.ReadFile(ac.CfgFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(f, &ac)
		if err != nil {
			return nil, err
		}
	}

	return &ac, nil
}
