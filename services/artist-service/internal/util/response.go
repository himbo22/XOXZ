package util

import (
	"net/http"

	"github.com/himbo22/xoxz/artist-service/internal/model"
	commonUtil "github.com/himbo22/xoxz/common-service/xoxz/util"
	"github.com/labstack/echo/v5"
)

func SuccessResponse(c *echo.Context, httpCode int, customCode int, message string, data any) error {
	if httpCode == http.StatusNoContent {
		return c.NoContent(httpCode)
	}

	return c.JSON(httpCode, model.ResponseModel{
		StatusCode: customCode,
		Message:    message,
		Error:      nil,
		Data:       data,
	})
}

func NewError(httpCode, customCode int, message string) *commonUtil.AppError {
	return &commonUtil.AppError{
		HTTPCode:   httpCode,
		CustomCode: customCode,
		Message:    message,
	}
}
