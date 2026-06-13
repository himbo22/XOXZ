package middleware

import (
	"github.com/himbo22/xoxz/account-service/internal/bootstrap"
	"github.com/himbo22/xoxz/account-service/internal/domain/repository"
)

type AuthMiddleware struct {
	redisRepo   repository.RedisRepository
	permissions *bootstrap.AccessControl
}

func NewAuthMiddleware(
	redisRepo repository.RedisRepository,
	permissions *bootstrap.AccessControl,
) *AuthMiddleware {
	return &AuthMiddleware{
		redisRepo:   redisRepo,
		permissions: permissions,
	}
}
