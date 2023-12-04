package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"tracing/cmd/opentelemetry/config"
)

type Router struct {
	cfg *config.Config
}

func Init(group *fiber.App, cfg config.Config) {
	group.All("*", func(c *fiber.Ctx) error {
		url := cfg.ServiceProxy.Host + c.Path()
		err := proxy.Do(c, url)
		if err != nil {
			panic(err)
		}

		return nil
	})
}
