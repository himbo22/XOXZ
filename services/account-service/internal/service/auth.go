package service

import (
	"context"

	"github.com/himbo22/xoxz/account-service/internal/logic"
	"github.com/himbo22/xoxz/account-service/internal/model"
	"github.com/himbo22/xoxz/common-service/monitoring/telemetry"
)

type AuthService interface {
	AuthenticateWithGoogle(ctx context.Context, req model.GoogleLoginRequest) (model.AuthTokenResponse, error)
	RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (model.AuthTokenResponse, error)
	Logout(ctx context.Context, req model.LogoutRequest) error
	RevokeAllSessions(ctx context.Context, req model.RevokeAllSessionsRequest) error
}

type authService struct {
	authLogic *logic.AuthLogic
}

func NewAuthService(authLogic *logic.AuthLogic) AuthService {
	return &authService{
		authLogic: authLogic,
	}
}

// RevokeAllSessions implements AuthService.
func (a *authService) RevokeAllSessions(ctx context.Context, req model.RevokeAllSessionsRequest) error {
	return a.authLogic.RevokeAllSessions(ctx, req)
}

func (a *authService) AuthenticateWithGoogle(ctx context.Context, req model.GoogleLoginRequest) (res model.AuthTokenResponse, err error) {
	ctx, end := telemetry.StartSpan(ctx, "test1", "service")
	defer func() { end(err) }()

	res, err = a.authLogic.AuthenticateWithGoogle(ctx, req)

	return res, err
}

func (a *authService) RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (model.AuthTokenResponse, error) {
	return a.authLogic.RefreshToken(ctx, req)
}

func (a *authService) Logout(ctx context.Context, req model.LogoutRequest) error {
	return a.authLogic.Logout(ctx, req)
}
