package media

import "github.com/labstack/echo/v5"

type MediaController interface {
	GeneratePresignedURL(ctx *echo.Context) error

	// grpc
	ConfirmUpload(ctx *echo.Context) error
	DeleteMedia(ctx *echo.Context) error
}
