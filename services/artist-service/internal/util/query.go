package util

import (
	"strconv"

	"github.com/labstack/echo/v5"
)

func ParseIntQuery(ctx *echo.Context, key string, defaultValue int) int {
	value := ctx.QueryParam(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}
