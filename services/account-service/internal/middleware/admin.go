package middleware

import (
	"github.com/google/uuid"
	_const "github.com/himbo22/xoxz/account-service/internal/const"
	"github.com/himbo22/xoxz/account-service/internal/util"
	"github.com/labstack/echo/v5"
)

func (m *AuthMiddleware) AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			roleID, ok := GetRoleId(ctx)
			if !ok || roleID == uuid.Nil {
				return util.NewErrorByCode(_const.CodeForbidden, "1: Invalid token")
			}

			if !m.permissions.HasPermission(roleID, string(_const.ACCOUNT_PROVIDE)) {
				return util.NewErrorByCode(_const.CodeForbidden, "2: Invalid permission")
			}

			return next(ctx)
		}
	}
}

// get role id from echo context
func GetRoleId(c *echo.Context) (uuid.UUID, bool) {
	v := c.Get(string(_const.RoleIDKey))
	if v == nil {
		return uuid.Nil, false
	}

	roleID, ok := v.(uuid.UUID)
	if !ok || roleID == uuid.Nil {
		return uuid.Nil, false
	}

	return roleID, true
}
