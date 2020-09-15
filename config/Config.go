package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	ConfigFile    string `env:"CONFIG_FILE,required"`
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDatabase int    `env:"REDIS_DATABASE" envDefault:"0"`
}

var EnvConfig Config

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		PadLevelText:    true,
	})
	EnvConfig = Config{}
	if err := env.Parse(&EnvConfig); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
