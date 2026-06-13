package di

import (
	"github.com/himbo22/xoxz/artist-service/internal/config"
	"github.com/labstack/echo/v5"
)

type App struct {
	Config  *config.Config
	EchoApp *echo.Echo
}
