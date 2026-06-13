package util

import "fmt"

type AppError struct {
	HTTPCode   int
	CustomCode int
	Message    string
	Detail     *string
}

func (e *AppError) Error() string {
	if e.Detail != nil && *e.Detail != "" {
		return fmt.Sprintf("Code: %d, Message: %s, Detail: %s", e.CustomCode, e.Message, *e.Detail)
	}
	return fmt.Sprintf("Code: %d, Message: %s", e.CustomCode, e.Message)
}
