package util

import (
	"net/http"

	"github.com/himbo22/xoxz/common-service/xoxz/model"
	"github.com/himbo22/xoxz/common-service/xoxz/util"
	_const "github.com/himbo22/xoxz/livestream-service/internal/const"
	"github.com/labstack/echo/v5"
)

func SuccessResponse(c *echo.Context, httpCode int, customCode int, message string, data any) error {
	// If 204 No Content, no need to return a JSON body
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

// Helper function to quickly create an error
func NewError(httpCode, customCode int, message string) *util.AppError {
	return &util.AppError{
		HTTPCode:   httpCode,
		CustomCode: customCode,
		Message:    message,
	}
}

func NewErrorByCode(code _const.CustomCode, message ...string) *util.AppError {
	msg := code.Message
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	return &util.AppError{
		HTTPCode:   code.HTTPStatus,
		CustomCode: code.Code,
		Message:    msg,
		Detail:     code.Detail,
	}
}
