package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
	"log"
	"sync"
)

type Config struct {
	Env          string `hcl:"env" env:"ENV" default:"local"`
	DatabaseDSN  string `hcl:"database_dsn" env:"DB_DSN" default:"postgres://postgres:postgres@localhost:5437/postgres?sslmode=disable"`
	AgentsNumber int    `hcl:"agents_number" env:"AGENTS_NUMBER" default:"2"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			EnvPrefix: "NFB",
			Files:     []string{"./config.hcl", "./config.local.hcl"},
			FileDecoders: map[string]aconfig.FileDecoder{
				".hcl": aconfighcl.New(),
			},
		})

		if err := loader.Load(); err != nil {
			log.Printf("failed to load config: %v", err)
		}
	})
	return cfg
}
