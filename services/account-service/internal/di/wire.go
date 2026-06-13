//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/himbo22/xoxz/account-service/internal/config"
)

func InitializeApp(cfg *config.Config) (*App, func(), error) {
	wire.Build(
		InfrastructureSet,
		AuthSet,
		ProfileSet,
		AdminSet,
		RepositorySet,
		ControllerSet,
		CacheSet,
		MiddlewareSet,
		NewApp,
	)
	return nil, nil, nil
}
