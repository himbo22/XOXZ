package router

import (
	"net/http"

	"github.com/himbo22/xoxz/media-service/internal/controller/http/media"
	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

type Controllers struct {
	MediaController media.MediaController
}

func SetupRouters(app *echo.Echo, controllers Controllers) {
	internal := app.Group("/api/v1/internal")
	{
		// Swagger
		internal.GET("/swagger/*", echoSwagger.WrapHandler)

		// Ping
		internal.GET("/ping", func(c *echo.Context) error {
			return c.String(200, "pong")
		})

		// healcheck
		internal.GET("/health", func(c *echo.Context) error {
			return c.NoContent(http.StatusOK)
		})
	}

	public := app.Group("/api/v1/public")
	{
		SetupMediaController(public, controllers.MediaController)
	}
}
