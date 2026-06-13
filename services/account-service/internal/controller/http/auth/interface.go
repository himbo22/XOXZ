package auth

import "github.com/labstack/echo/v5"

type AuthController interface {
	Google(ctx *echo.Context) error
	Refresh(ctx *echo.Context) error
	Logout(ctx *echo.Context) error
	RevokeAllSessions(ctx *echo.Context) error
}
