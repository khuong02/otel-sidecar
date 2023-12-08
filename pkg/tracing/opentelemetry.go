package tracing

import (
	"context"
	"crypto/x509"
	"github.com/gofiber/fiber/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	"os"
)

type TracerConfig struct {
	CollectorURL       string
	Insecure           string
	ServiceName        string
	Lang               string
	ServerNameOverride string
	CP                 *x509.CertPool
	Stdouttrace        bool
	Env                string
}

type Tracer struct {
	cfg TracerConfig
}

func NewTracer(cfg TracerConfig) *Tracer {
	return &Tracer{
		cfg: cfg,
	}
}

func (t *Tracer) InitTracer() func(context.Context) error {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	options := make([]sdktrace.TracerProviderOption, 0)
	options = append(options, sdktrace.WithSampler(sdktrace.AlwaysSample()))

	if t.cfg.Stdouttrace {
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			log.Error("Could not set resources: ", "err:", err)

			os.Exit(1)
		}

		options = append(options, sdktrace.WithBatcher(exporter))
	}

	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(t.cfg.CP, t.cfg.ServerNameOverride))
	if len(t.cfg.Insecure) > 0 {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(t.cfg.CollectorURL),
		),
	)
	if err != nil {
		log.Error("Exporter fail", "err:", err)

		os.Exit(1)
	}

	options = append(options, sdktrace.WithBatcher(exporter))

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", t.cfg.ServiceName),
			attribute.String("library.language", t.cfg.Lang),
		),
	)
	if err != nil {
		log.Warn("Could not set resources: ", "err:", err)
	}

	options = append(options, sdktrace.WithResource(resources))

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			options...,
		),
	)

	return exporter.Shutdown
}
