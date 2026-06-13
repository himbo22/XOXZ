package model

type GoogleRequest struct {
	Token string `json:"token"`
}

type TokenPayload struct {
	Issuer   string         `json:"iss"`
	Audience string         `json:"aud"`
	Expires  int64          `json:"exp"`
	IssuedAt int64          `json:"iat"`
	Subject  string         `json:"sub,omitempty"`
	Claims   map[string]any `json:"claims"`
}
