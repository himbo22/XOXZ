package _const

import "net/http"

type CustomCode struct {
	Code       int
	Message    string
	Detail     *string
	HTTPStatus int
}

var (
	CodeSuccess = CustomCode{Code: 0, Message: "success", HTTPStatus: http.StatusOK}

	CodeInvalidRequest = CustomCode{Code: 4000, Message: "invalid request", HTTPStatus: http.StatusBadRequest}
	CodeArtistNotFound = CustomCode{Code: 4004, Message: "artist not found", HTTPStatus: http.StatusNotFound}
	CodeInternalError  = CustomCode{Code: 9999, Message: "internal server error", HTTPStatus: http.StatusInternalServerError}
)
