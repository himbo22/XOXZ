package oauth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/himbo22/xoxz/account-service/internal/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuthAdapter struct {
	config     *oauth2.Config
	httpClient *http.Client
}

type GoogleOAuthConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURL  string `yaml:"redirect_url"`
}

func NewGoogleOAuthAdapter(cfg GoogleOAuthConfig) *GoogleOAuthAdapter {
	return &GoogleOAuthAdapter{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (g *GoogleOAuthAdapter) VerifyGoogleToken(ctx context.Context, token, clientID string) (*model.TokenPayload, error) {
	payload, err := idtoken.Validate(ctx, token, clientID)
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
