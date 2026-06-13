package artist

import "github.com/labstack/echo/v5"

type ArtistController interface {
	Create(ctx *echo.Context) error
	GetByID(ctx *echo.Context) error
	List(ctx *echo.Context) error
	Update(ctx *echo.Context) error
	Delete(ctx *echo.Context) error
}
