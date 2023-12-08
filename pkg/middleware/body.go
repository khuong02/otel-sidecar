package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"
	"tracing/utils/common"
)

const (
	Request  = "Request"
	Response = "Response"
)

type BodyLog struct {
	startTime time.Time
	span      trace.Span
}

func GetContextLoggerMiddleWare() fiber.Handler {
	bodyLog := BodyLog{}

	return func(c *fiber.Ctx) error {
		reqCtx := c.UserContext()

		reqHeader := make(http.Header)
		c.Request().Header.VisitAll(func(k, v []byte) {
			reqHeader.Add(string(k), string(v))
		})
		span := trace.SpanFromContext(otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(reqHeader)))
		defer span.SpanContext()

		bodyLog.span = span

		bodyLog.PrintLogging(c, Request)

		c.SetUserContext(reqCtx)

		bodyLog.startTime = time.Now()
		err := c.Next()
		if err != nil {
			return err
		}

		bodyLog.PrintLogging(c, Response)

		return err
	}
}

func (b BodyLog) SetLoggingRequestMethod(c *fiber.Ctx) map[string]interface{} {
	info := map[string]interface{}{
		"headers": c.GetReqHeaders(),
		"path":    c.Path(),
		"method":  c.Method(),
		"spanID":  b.span.SpanContext().SpanID().String(),
		"traceID": b.span.SpanContext().TraceID().String(),
	}

	if c.Request().Body() != nil {
		info["body"] = string(c.Request().Body())
	}

	if c.Queries() != nil {
		info["queries"] = c.Queries()
	}

	return info
}

func (b BodyLog) SetLoggingResponseMethod(c *fiber.Ctx) map[string]interface{} {
	status := c.Response().StatusCode()
	timeElapsed := time.Since(b.startTime).Milliseconds()

	info := map[string]interface{}{
		"status":      status,
		"body":        string(c.Response().Body()),
		"headers":     c.GetReqHeaders(),
		"path":        c.Path(),
		"method":      c.Method(),
		"timeElapsed": fmt.Sprintf("%v ms", timeElapsed),
		"spanID":      b.span.SpanContext().SpanID().String(),
		"traceID":     b.span.SpanContext().TraceID().String(),
	}

	return info
}

func (b BodyLog) PrintLogging(c *fiber.Ctx, kindLog string) {
	strLogging := ""

	switch kindLog {
	case Request:
		strLogging = common.MapToString(b.SetLoggingRequestMethod(c))

	default:
		strLogging = common.MapToString(b.SetLoggingResponseMethod(c))
	}

	log.Infof(
		"Logging %s info, %v",
		kindLog,
		strLogging,
	)
}
