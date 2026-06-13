package repository

import (
	"context"

	"github.com/himbo22/xoxz/account-service/internal/domain/entity"
)

type IdentityRepository interface {
	FindBySubjectAndProvider(ctx context.Context, subject string, provider string) (*entity.Identity, error)
	Upsert(ctx context.Context, identity *entity.Identity) error
}
