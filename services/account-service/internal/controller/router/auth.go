package router

import (
	"github.com/himbo22/xoxz/account-service/internal/controller/http/auth"
	"github.com/himbo22/xoxz/account-service/internal/middleware"
	"github.com/labstack/echo/v5"
)

func SetupAuthRoutes(parentRouter *echo.Group, authController auth.AuthController, authMiddleware *middleware.AuthMiddleware) {
	authGroup := parentRouter.Group("/auth")

	// ==========================================
	// 1. PUBLIC ROUTES
	// ==========================================
	authGroup.POST("/google", authController.Google)

	// ==========================================
	// 2. Session Group
	// ==========================================
	// we set /session for logout and refresh because logout need refresh token to remove it from redis
	sessionGroup := authGroup.Group("/session")
	{
		sessionGroup.POST("/refresh", authController.Refresh)
		sessionGroup.POST("/logout", authController.Logout, authMiddleware.RequireAuth())
		sessionGroup.POST("/revoke-all", authController.RevokeAllSessions, authMiddleware.RequireAuth())
	}
}
