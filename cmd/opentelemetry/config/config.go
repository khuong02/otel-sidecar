package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	configs "tracing/pkg/config"
)

type Config struct {
	configs.App    `yaml:"APP"`
	configs.Http   `yaml:"HTTP"`
	configs.Tracer `yaml:"TRACER"`
	ServiceProxy   *ServiceProxy `yaml:"SERVICE_PROXY"`
}

type ServiceProxy struct {
	Host string `env-required:"true" yaml:"HOST"    env:"SERVICE_PROXY_HOST"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()
	//flags.GetFlag()
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// debug
	fmt.Println(dir + "/config.yaml")

	err = cleanenv.ReadConfig(dir+"/config.yaml", cfg)
	if err != nil {
		err = cleanenv.ReadEnv(cfg)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	return cfg, nil
}
