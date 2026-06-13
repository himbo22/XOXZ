package profile

import "github.com/labstack/echo/v5"

type ProfileController interface {
	// Me pattern (api/profile/me): require auth
	GetProfile(ctx *echo.Context) error
	UpdateProfile(ctx *echo.Context) error
	UpdateAvatar(ctx *echo.Context) error

	// public profile
	GetPublicProfile(c *echo.Context) error
}
