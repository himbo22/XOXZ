package webhook

import (
	"github.com/labstack/echo/v5"
)

type WebhookController interface {
	ServeHTTP(ctx *echo.Context) error
}
