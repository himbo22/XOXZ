package router

import (
	"github.com/himbo22/xoxz/account-service/internal/controller/http/admin"
	"github.com/himbo22/xoxz/account-service/internal/middleware"
	"github.com/labstack/echo/v5"
)

func SetupAdminRoutes(
	parent *echo.Group,
	adminController admin.AdminController,
	authMiddleware *middleware.AuthMiddleware,
) {
	adminGroup := parent.Group("/artist", authMiddleware.AdminMiddleware())
	{
		adminGroup.POST("/account", adminController.CreateArtistInvite)
		adminGroup.DELETE("/account", adminController.RevokeArtistAccount)
	}
}
