package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" json:"server_address"`
	DatabaseDSN   string `env:"DATABASE_DSN" json:"database_dsn"`
	ConfigFile    string `env:"CONFIG"`
	JWTSecret     string `env:"JWT_SECRET"`
}

func (cfg *Config) Init() {

	flag.StringVarP(&cfg.ServerAddress, "address", "a", "localhost:8080", "server address")
	flag.StringVarP(&cfg.DatabaseDSN, "database", "d", "host=localhost port=5555 user=postgres password=password dbname=postgres sslmode=disable", "database connection string")
	flag.StringVarP(&cfg.ConfigFile, "config", "c", "", "file config path")
	flag.StringVarP(&cfg.JWTSecret, "jwt", "j", "^is_j5Wxv;byl]#oI[LPH8+CFh&}eY3>", "jwt secret key")

	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.ConfigFile != "" {
		var c Config

		f, err := os.Open(cfg.ConfigFile)

		if err != nil {
			log.Fatal("Error while open file")
		}

		if err := json.NewDecoder(f).Decode(&c); err != nil {
			log.Fatal("Error while decoding config")
		}

		if cfg.ServerAddress == "" {
			cfg.ServerAddress = c.ServerAddress
		}

		if cfg.DatabaseDSN == "" {
			cfg.DatabaseDSN = c.DatabaseDSN
		}

		if cfg.JWTSecret == "" {
			cfg.JWTSecret = c.JWTSecret
		}
	}

}
