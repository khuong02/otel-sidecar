package app

import (
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"log"
	"strings"
	"tracing/cmd/opentelemetry/config"
	"tracing/pkg/middleware"
)

type App struct {
	Cfg config.Config
}

func New(
	cfg config.Config,
) *App {
	return &App{
		Cfg: cfg,
	}
}

func (a *App) Init() *fiber.App {
	r := fiber.New()

	r.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOriginsFunc: nil,
		AllowOrigins:     "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		MaxAge: 0,
	}))

	r.Use(otelfiber.Middleware())
	r.Use(middleware.GetContextLoggerMiddleWare())

	r.Use(proxy.Balancer(proxy.Config{
		Servers: a.Cfg.ServiceProxy.Hosts,
	}))

	return r
}

func (a *App) Run() {
	var err error
	errs := make(chan error)

	// http
	go func() {
		errs <- a.Init().Listen(a.Cfg.Http.Host + ":" + a.Cfg.Http.Port)
	}()

	err = <-errs
	if err != nil {
		log.Fatal("Start service fail", "err: ", err)
	}
}
