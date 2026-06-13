package router

import (
	"net/http"

	"github.com/himbo22/xoxz/account-service/internal/controller/http/admin"
	"github.com/himbo22/xoxz/account-service/internal/controller/http/auth"
	"github.com/himbo22/xoxz/account-service/internal/controller/http/profile"
	"github.com/himbo22/xoxz/account-service/internal/middleware"
	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

type Controllers struct {
	AuthController    auth.AuthController
	ProfileController profile.ProfileController
	AdminController   admin.AdminController
}

func SetupRouters(app *echo.Echo, controllers Controllers, authMiddleware *middleware.AuthMiddleware) {
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
		SetupAuthRoutes(public, controllers.AuthController, authMiddleware)
		SetupProfileRoutes(public, controllers.ProfileController, authMiddleware)
		SetupAdminRoutes(public, controllers.AdminController, authMiddleware)
	}
}
