package app

import (
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"strings"
	"tracing/cmd/opentelemetry/config"
	"tracing/internal/opentelemetry/router"
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
		AllowHeaders: "Origin,Authorization,Content-Type",
		MaxAge:       0,
	}))

	r.Use(otelfiber.Middleware())
	r.Use(middleware.GetContextLoggerMiddleWare())

	router.Init(r, a.Cfg)

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
