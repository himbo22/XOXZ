package router

import (
	profileController "github.com/himbo22/xoxz/account-service/internal/controller/http/profile"
	"github.com/himbo22/xoxz/account-service/internal/middleware"
	"github.com/labstack/echo/v5"
)

func SetupProfileRoutes(parentRouter *echo.Group, profile profileController.ProfileController, authMiddleware *middleware.AuthMiddleware) {
	profileGroup := parentRouter.Group("/profile")

	profileGroup.GET("/:username", profile.GetPublicProfile)

	meGroup := profileGroup.Group("/me", authMiddleware.RequireAuth())
	{
		meGroup.GET("", profile.GetProfile)
		meGroup.PUT("", profile.UpdateProfile, authMiddleware.SecureSession())
		meGroup.PUT("/avatar", profile.UpdateAvatar, authMiddleware.SecureSession())
	}
}
