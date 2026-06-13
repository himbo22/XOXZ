package middleware

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/himbo22/xoxz/account-service/internal/config"
	_const "github.com/himbo22/xoxz/account-service/internal/const"
	"github.com/himbo22/xoxz/account-service/internal/util"
	"github.com/labstack/echo/v5"
)

// for authenticated user
func (m *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// get deviceID from header
			deviceID, ok := GetDeviceIDFromHeader(c)
			if !ok {
				return util.NewErrorByCode(_const.CodeInvalidRequest, "Missing device ID")
			}

			// Device Type
			deviceType, ok := GetDeviceTypeFromHeader(c)
			if !ok {
				return util.NewErrorByCode(_const.CodeInvalidRequest, "Invalid device type")
			}

			// Extract & Validate Token
			authHeader := c.Request().Header.Get("Authorization")
			var accessToken string
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				accessToken = strings.TrimSpace(parts[1])
			}

			claims, err := util.ParseTokenWithClaims(accessToken, config.AppPublicKey)
			if err != nil {
				return err
			}
			if claims == nil {
				return util.NewErrorByCode(_const.CodeUnauthorized)
			}

			// Red alert
			if claims.DeviceID != deviceID {
				_ = m.redisRepo.RevokeAllSessions(c.Request().Context(), claims.UserID)
				return util.NewErrorByCode(_const.CodeForbidden, "Unusual access detected")
			}

			// TODO: check role

			// 2. Store values in Echo context for the transport/HTTP layer.
			c.Set(string(_const.DeviceTypeKey), deviceType)
			c.Set(string(_const.RoleIDKey), claims.RoleID)
			c.Set(string(_const.DeviceIDKey), claims.DeviceID)
			c.Set(string(_const.UserIDKey), claims.UserID)
			c.Set(string(_const.IssuerAtKey), claims.RegisteredClaims.IssuedAt.Unix())
			c.Set(string(_const.AccessTokenKey), accessToken)

			// 3. Pass control to the next handler.
			return next(c)
		}
	}
}

// layer 2
func (m *AuthMiddleware) SecureSession() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userID, ok := GetUserID(c)
			if !ok {
				return util.NewErrorByCode(_const.CodeUnauthorized, "invalid token")
			}
			key := _const.BlackAccessTokenKey(userID)
			logoutTime, err := m.redisRepo.Get(c.Request().Context(), key)
			if err != nil {
				return err
			}
			if logoutTime == "" {
				return next(c)
			}
			iat := GetIssuerAt(c)
			lgTime, err := strconv.ParseInt(logoutTime, 10, 64)
			if err != nil {
				return err
			}
			if iat <= lgTime {
				return util.NewErrorByCode(_const.CodeInvalidRequest, "invalid token")
			}
			return next(c)
		}
	}
}

func GetIssuerAt(c *echo.Context) int64 {
	v := c.Get(string(_const.IssuerAtKey))
	if v == nil {
		return 0
	}

	iat, ok := v.(int64)
	if !ok {
		return 0
	}

	return iat
}

func GetUserID(c *echo.Context) (uuid.UUID, bool) {
	v := c.Get(string(_const.UserIDKey))
	if v == nil {
		return uuid.Nil, false
	}

	userID, ok := v.(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return uuid.Nil, false
	}

	return userID, true
}

// get AT from echo context
func GetAccessToken(c *echo.Context) string {
	v := c.Get(string(_const.AccessTokenKey))
	if v == nil {
		return ""
	}

	token, ok := v.(string)
	if !ok || token == "" {
		return ""
	}

	return token
}

func GetDeviceIDFromHeader(c *echo.Context) (string, bool) {
	deviceID := c.Request().Header.Get(string(_const.HeaderDeviceID))
	if deviceID == "" {
		return "", false
	}
	return deviceID, true
}

func GetDeviceTypeFromHeader(c *echo.Context) (_const.DeviceType, bool) {
	deviceType := c.Request().Header.Get(string(_const.HeaderDeviceType))
	if deviceType == "" {
		return "", false
	}
	dt := util.ParseDeviceType(deviceType)
	if dt == _const.NilType {
		return "", false
	}
	return dt, true
}

// get device type from echo context
func GetDeviceType(c *echo.Context) _const.DeviceType {
	v := c.Get(string(_const.DeviceTypeKey))
	if v == nil {
		return _const.NilType
	}
	deviceType, ok := v.(_const.DeviceType)
	if !ok {
		return _const.NilType
	}

	return deviceType
}

// get device id from echo context
func GetDeviceID(c *echo.Context) string {
	v := c.Get(string(_const.DeviceIDKey))
	if v == nil {
		return ""
	}

	deviceID, ok := v.(string)
	if !ok {
		return ""
	}

	return deviceID
}

// Get refresh token from cookies
func GetRefreshToken(c *echo.Context) (string, bool) {
	cookie, err := c.Cookie(string(_const.RefreshTokenCookieKey))
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}
