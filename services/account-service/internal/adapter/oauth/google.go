package oauth

import (
	"context"
	"errors"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/himbo22/xoxz/account-service/internal/model"
)

func VerifyGoogleToken(ctx context.Context, token string, clientId string) (*model.TokenPayload, error) {
	payload, err := idtoken.Validate(ctx, token, clientId)
	if err != nil {
		return nil, err
	}
	if payload == nil {
		return nil, errors.New("google token verification success but payload is unexpectedly nil")
	}
	return &model.TokenPayload{
		Issuer:   payload.Issuer,
		Audience: payload.Audience,
		Expires:  payload.Expires,
		IssuedAt: payload.IssuedAt,
		Subject:  payload.Subject,
		Claims:   payload.Claims,
	}, nil
}
