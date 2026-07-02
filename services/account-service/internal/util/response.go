package util

import (
	"net/http"

	_const "github.com/himbo22/xoxz/account-service/internal/const"
	"github.com/himbo22/xoxz/account-service/internal/model"
	"github.com/labstack/echo/v5"
)

func SuccessResponse(c *echo.Context, httpCode int, customCode int, message string, data any) error {
	// For 204 No Content, do not return a JSON body.
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

func SuccessResponseByCode(c *echo.Context, code _const.CustomCode, data any) error {
	// For 204 No Content, do not return a JSON body.
	if code.HTTPStatus == http.StatusNoContent {
		return c.NoContent(code.HTTPStatus)
	}

	return c.JSON(code.HTTPStatus, model.ResponseModel{
		StatusCode: code.Code,
		Message:    code.Message,
		Error:      nil,
		Data:       data,
	})
}

// NewError creates an AppError quickly.
func NewError(httpCode, customCode int, message string) *AppError {
	return &AppError{
		HTTPCode:   httpCode,
		CustomCode: customCode,
		Message:    message,
	}
}

func NewErrorByCode(code _const.CustomCode, message ...string) *AppError {
	msg := code.Message
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	return &AppError{
		HTTPCode:   code.HTTPStatus,
		CustomCode: code.Code,
		Message:    msg,
		Detail:     code.Detail,
	}
}
