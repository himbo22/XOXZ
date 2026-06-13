package router

import (
	"github.com/himbo22/xoxz/media-service/internal/controller/http/media"
	"github.com/labstack/echo/v5"
)

func SetupMediaController(parentRouter *echo.Group, mediaCo media.MediaController) {
	parentRouter.POST("/create-presigned-url", mediaCo.GeneratePresignedURL)
}
