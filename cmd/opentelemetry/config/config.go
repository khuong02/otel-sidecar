package config

import (
	"fmt"
	"log"
	"os"
	configs "tracing/pkg/config"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	configs.App    `yaml:"APP"`
	configs.Http   `yaml:"HTTP"`
	configs.Tracer `yaml:"TRACER"`
	ServiceProxy   *ServiceProxy `yaml:"SERVICE_PROXY"`
}

type ServiceProxy struct {
	Hosts []string `env-required:"true" yaml:"HOSTS"    env:"SERVICE_PROXY_HOSTS"`
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
	// overwrite default env
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
