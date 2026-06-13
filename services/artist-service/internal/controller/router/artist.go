package router

import (
	"github.com/himbo22/xoxz/artist-service/internal/controller/http/artist"
	"github.com/labstack/echo/v5"
)

func SetupArtistRoutes(parentRouter *echo.Group, artistController artist.ArtistController) {
	artistRouteGroup := parentRouter.Group("/artists")
	{
		artistRouteGroup.POST("", artistController.Create)
		artistRouteGroup.GET("", artistController.List)
		artistRouteGroup.GET("/:id", artistController.GetByID)
		artistRouteGroup.PUT("/:id", artistController.Update)
		artistRouteGroup.DELETE("/:id", artistController.Delete)
	}
}
