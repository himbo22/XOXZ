package util

import (
	"fmt"

	_const "github.com/himbo22/xoxz/media-service/internal/const"
)

type AppError struct {
	HTTPCode   int
	CustomCode int
	Message    string
	Detail     *string
}

// Must implement the Golang error interface
func (e *AppError) Error() string {
	if e.Detail != nil && *e.Detail != "" {
		return fmt.Sprintf("Code: %d, Message: %s, Detail: %s", e.CustomCode, e.Message, *e.Detail)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.CustomCode, e.Message)
}

// Helper function to quickly create an error
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
