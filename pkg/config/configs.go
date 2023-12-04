package config

import (
	"tracing/utils/stage"
)

type (
	App struct {
		Name    string          `env-required:"true" yaml:"NAME"    env:"APP_NAME"`
		Version string          `env-required:"true" yaml:"VERSION" env:"APP_VERSION"`
		Stage   stage.StageType `env-required:"true" yaml:"STAGE"    env:"APP_STAGE"`
	}
	Http struct {
		Port string `env-required:"true" yaml:"PORT"    env:"HTTP_PORT"`
		Host string `env-required:"true" yaml:"HOST"    env:"HTTP_HOST"`
	}

	Tracer struct {
		CollectorURL       string `env-required:"true" yaml:"COLLECTOR_URL"    env:"TRACER_COLLECTOR_URL"`
		Insecure           string `env-required:"true" yaml:"INSECURE"    env:"TRACER_INSECURE"`
		Lang               string `env-required:"true" yaml:"LANG"    env:"TRACER_LANG"`
		ServerNameOverride string `yaml:"SERVER_NAME_OVERRIDE"    env:"TRACER_SERVER_NAME_OVERRIDE"`
		StdOutTrace        bool   `yaml:"STD_OUT_TRACE"    env:"TRACER_STD_OUT_TRACE"`
		CPPath             string `yaml:"CP_PATH"    env:"TRACER_CP_PATH"`
	}
)
