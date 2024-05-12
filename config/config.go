package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Name       string `env:"APP_NAME"`
	Production bool   `env:"PRODUCTION"`
	Token      string `env:"TOKEN"`
}

func MustLoadConfig(env ...string) *AppConfig {
	var err error
	conf := &AppConfig{}
	err = godotenv.Load(env...)
	if err != nil {
		panic(err)
	}

	err = cleanenv.ReadEnv(conf)
	if err != nil {
		panic(err)
	}

	return conf
}
