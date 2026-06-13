package model

type ResponseModel struct {
	StatusCode int     `json:"status_code"`
	Message    string  `json:"message"`
	Error      *string `json:"error,omitempty"`
	Data       any     `json:"data,omitempty"`
}

type AuthTokenResponseDoc struct {
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Error      *string           `json:"error,omitempty"`
	Data       AuthTokenResponse `json:"data"`
}

type MessageResponseDoc struct {
	StatusCode int     `json:"status_code"`
	Message    string  `json:"message"`
	Error      *string `json:"error,omitempty"`
	Data       any     `json:"data,omitempty"`
}
