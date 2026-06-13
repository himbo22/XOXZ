package _const

import (
	"errors"
	"net/http"
)

type CustomCode struct {
	Code       int
	Message    string
	Detail     *string
	HTTPStatus int
}

func (c CustomCode) WithMessage(message string) CustomCode {
	c.Message = message
	return c
}

func (c CustomCode) WithDetail(detail string) CustomCode {
	c.Detail = &detail
	return c
}

var (
	ErrRecordNotFound = errors.New("domain: record not found")
	ErrInternalServer = errors.New("domain: internal server error")
	ErrInvalidData    = errors.New("domain: invalid data format")
)

var (
	CodeSuccess              = CustomCode{Code: 2000, Message: "success", HTTPStatus: http.StatusOK}
	CodeSuccessNoContent     = CustomCode{Code: 2001, Message: "success", HTTPStatus: http.StatusNoContent}
	CodeInvalidRequest       = CustomCode{Code: 4000, Message: "invalid request", HTTPStatus: http.StatusBadRequest}
	CodeInvalidGoogleToken   = CustomCode{Code: 4001, Message: "invalid google token", HTTPStatus: http.StatusUnauthorized}
	CodeUserEmailNotVerified = CustomCode{Code: 4002, Message: "email is not verified", HTTPStatus: http.StatusForbidden}
	CodeUserNotFound         = CustomCode{Code: 4003, Message: "user not found", HTTPStatus: http.StatusNotFound}
	CodeGetRoleFailed        = CustomCode{Code: 4004, Message: "Get role fail", HTTPStatus: http.StatusOK}
	CodeUpdateUserFailed     = CustomCode{Code: 4005, Message: "Update user fail", HTTPStatus: http.StatusOK}
	CodeCreateUserFailed     = CustomCode{Code: 4006, Message: "Create user fail", HTTPStatus: http.StatusOK}
	CodeUpdateIdentityFailed = CustomCode{Code: 4007, Message: "Update identity fail", HTTPStatus: http.StatusOK}
	CodeMarshalPayloadFailed = CustomCode{Code: 4008, Message: "Marshal payload fail", HTTPStatus: http.StatusOK}
	CodeUnauthorized         = CustomCode{Code: 4009, Message: "Unauthorized user", HTTPStatus: http.StatusUnauthorized}
	CodeForbidden            = CustomCode{Code: 4009, Message: "Forbidden user", HTTPStatus: http.StatusForbidden}
	CodeRoleNotFound         = CustomCode{Code: 4010, Message: "Role not found", HTTPStatus: http.StatusForbidden}
	CodeExpiredAccessToken   = CustomCode{Code: 4011, Message: "expired access token", HTTPStatus: http.StatusUnauthorized}
	CodeExpiredRefreshToken  = CustomCode{Code: 4012, Message: "expired refresh token", HTTPStatus: http.StatusUnauthorized}
	CodeStorageError         = CustomCode{Code: 4013, Message: "invalid file", HTTPStatus: http.StatusUnprocessableEntity}
	CodeInternalError        = CustomCode{Code: 5000, Message: "Internal error", HTTPStatus: http.StatusInternalServerError}
)
