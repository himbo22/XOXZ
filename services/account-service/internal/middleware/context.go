package middleware

import (
	"time"

	_const "github.com/himbo22/xoxz/account-service/internal/const"
	"github.com/labstack/echo/v5"
)

// ContextMiddleware extracts client information from headers and stores in context
func ContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Extract IP address
			ipAddress := c.RealIP()

			// Timezone
			timezone := c.Request().Header.Get("X-Timezone")
			if timezone == "" {
				timezone = "Asia/Ho_Chi_Minh"
			} else if _, err := time.LoadLocation(timezone); err != nil {
				timezone = "Asia/Ho_Chi_Minh"
			}

			// User-Agent
			userAgent := c.Request().Header.Get("User-Agent")

			// Store in context
			c.Set(string(_const.IPAddressKey), ipAddress)
			c.Set(string(_const.TimezoneKey), timezone)
			c.Set(string(_const.UserAgentKey), userAgent)
			// c.Set(string(AccessTokenKey), accessToken)

			return next(c)
		}
	}
}

// SetContextData sets IP address, timezone and user agent data into the fiber context
func SetContextData(c *echo.Context, ipAddress, timezone, userAgent string) {
	c.Set(string(_const.IPAddressKey), ipAddress)
	c.Set(string(_const.TimezoneKey), timezone)
	c.Set(string(_const.UserAgentKey), userAgent)
}

// GetIPAddress retrieves IP address from context
func GetIPAddress(c *echo.Context) string {
	if ip, ok := c.Get(string(_const.IPAddressKey)).(string); ok {
		return ip
	}
	return ""
}

// GetTimezone retrieves timezone from context
func GetTimezone(c *echo.Context) string {
	if tz, ok := c.Get(string(_const.TimezoneKey)).(string); ok {
		return tz
	}
	return "Asia/Ho_Chi_Minh"
}

// GetUserAgent retrieves user agent from context
func GetUserAgent(c *echo.Context) string {
	if ua, ok := c.Get(string(_const.UserAgentKey)).(string); ok {
		return ua
	}
	return ""
}
