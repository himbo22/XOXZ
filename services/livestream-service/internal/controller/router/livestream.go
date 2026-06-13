package router

import (
	"github.com/himbo22/xoxz/livestream-service/internal/controller/http/livestream"
	"github.com/labstack/echo/v5"
)

func SetupLiveRoutes(parentRouter *echo.Group, liveController livestream.LiveStreamController) {
	liveRouteGroup := parentRouter.Group("/streams")
	{
		liveRouteGroup.POST("", liveController.Create)
		liveRouteGroup.PUT("/:id", liveController.Stop)
	}
}
