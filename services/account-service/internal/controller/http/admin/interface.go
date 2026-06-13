package admin

import "github.com/labstack/echo/v5"

type AdminController interface {
	CreateArtistInvite(ctx *echo.Context) error
	RevokeArtistAccount(ctx *echo.Context) error
}
