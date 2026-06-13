package livestream

import "github.com/labstack/echo/v5"

type LiveStreamController interface {
	Create(ctx *echo.Context) error
	Stop(ctx *echo.Context) error
}
