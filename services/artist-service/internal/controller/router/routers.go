package router

import (
	"net/http"

	"github.com/himbo22/xoxz/artist-service/internal/controller/http/artist"
	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

type Controllers struct {
	ArtistController artist.ArtistController
}

func SetupRouters(app *echo.Echo, controllers Controllers) {
	internal := app.Group("/api/v1/internal")
	{
		internal.GET("/swagger/*", echoSwagger.WrapHandler)
		internal.GET("/ping", func(c *echo.Context) error {
			return c.String(http.StatusOK, "pong")
		})
		internal.GET("/health", func(c *echo.Context) error {
			return c.NoContent(http.StatusOK)
		})
	}

	public := app.Group("/api/v1/public")
	{
		SetupArtistRoutes(public, controllers.ArtistController)
	}
}
