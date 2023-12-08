package main

import (
	"context"
	"crypto/x509"
	"log"
	"strings"
	"tracing/cmd/opentelemetry/config"
	"tracing/internal/opentelemetry/app"
	"tracing/pkg/tracing"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// init trace
	var cp *x509.CertPool
	if strings.TrimSpace(cfg.Tracer.CPPath) == "" {
		cp = nil
	}

	tracerCfg := tracing.TracerConfig{
		CollectorURL:       cfg.Tracer.CollectorURL,
		Insecure:           cfg.Tracer.Insecure,
		ServiceName:        cfg.App.Name,
		Lang:               cfg.Tracer.Lang,
		ServerNameOverride: cfg.Tracer.ServerNameOverride,
		Stdouttrace:        cfg.Tracer.StdOutTrace,
		CP:                 cp,
	}

	tracer := tracing.NewTracer(tracerCfg)
	cleanup := tracer.InitTracer()
	defer cleanup(context.Background())

	service := app.New(*cfg)
	service.Run()
}
